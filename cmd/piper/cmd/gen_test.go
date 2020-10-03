package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenCommand_Success(t *testing.T) {
	const testCount = 9

	for ti := 0; ti < testCount; ti++ {
		_, filename, _, _ := runtime.Caller(0)
		inputPath := fmt.Sprintf("samples/inputs/s%03d", ti)
		inputPath = path.Join(path.Dir(filename), "../../..", inputPath)

		// assert.NoError(t, command.Run())
		t.Run(fmt.Sprintf("case: %s", inputPath), func(t *testing.T) {
			_ = os.Remove(path.Join(inputPath, "piper_gen.go"))

			genCommand := newGenCommand()
			genCommand.SetArgs([]string{inputPath})
			err := genCommand.Execute()
			assert.NoError(t, err)

			command := exec.Command("go", "test")
			command.Dir = inputPath

			output, err := command.CombinedOutput()
			if err != nil {
				t.Log(string(output))
			}
			assert.NoError(t, err)
		})
	}
}

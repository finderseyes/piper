package cmd

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenCommand_Success(t *testing.T) {
	const testCount = 8

	for ti := 0; ti < testCount; ti++ {
		_, filename, _, _ := runtime.Caller(0)
		inputPath := fmt.Sprintf("samples/inputs/s%03d", ti)
		inputPath = path.Join(path.Dir(filename), "../../..", inputPath)

		// assert.NoError(t, command.Run())
		t.Run(inputPath, func(t *testing.T) {
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

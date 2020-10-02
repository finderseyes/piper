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

func TestNewRootCommand(t *testing.T) {
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
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err = command.Run()
			assert.NoError(t, err)
			// fmt.Printf("%s\n", output)
		})
	}
}

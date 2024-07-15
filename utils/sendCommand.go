package utils

import (
	"errors"
	"fmt"
	"os"
)

func SendCommandPipe(cmd string, pipeName string) error {
	pipe, err := os.OpenFile(pipeName, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return errors.New("No running server, make one with the `serve` command.")
	}
	defer pipe.Close()

	if _, err := pipe.Write([]byte(cmd)); err != nil {
		return errors.New(fmt.Sprintf("failed to write to pipe: %v", err))
	}
	return nil
}

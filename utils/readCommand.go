package utils

import (
	"io"
	"log"
	"os"
	"strings"
)

func ReadCommandPipe(pipeName string) string {

	for {
		pipe, err := os.OpenFile(pipeName, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatalf("failed to open pipe: %v", err)
		}
		defer pipe.Close()

		// Read the command from the named pipe
		buf := make([]byte, 1024)
		n, err := pipe.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("failed to read from pipe: %v", err)
		}

		return strings.TrimSpace(string(buf[:n]))
	}
}

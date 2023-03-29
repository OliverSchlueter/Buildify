package util

import (
	"io"
	"os"
)

var (
	BuildScriptPath *string
	ResultPath      *string
)

func FastCopyFile(source, destination *os.File) error {
	buf := make([]byte, 4096)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

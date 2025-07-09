package storage

import (
	"io"
	"os"
)

func writePieceToFile(filePath string, offset int64, reader io.Reader, length int64) (int64, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR, defaultFileMode)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	if _, err = file.Seek(offset, io.SeekStart); err != nil {
		return 0, err
	}
	return tryWriteWithBuffer(file, reader, length)
}

func readPieceFromFile(filePath string, offset int64, length int64) (io.Reader, io.Closer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	if _, err = file.Seek(offset, io.SeekStart); err != nil {
		file.Close()
		return nil, nil, err
	}
	return io.LimitReader(file, length), file, nil
}

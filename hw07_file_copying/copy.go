package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrFileFromNotProvided   = errors.New("file from not provided")
	ErrFileToNotProvided     = errors.New("file to not provided")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type Copier struct {
	pathFrom string
	pathTo   string
	offset   int64
	limit    int64
}

func buildCopier(fromPath, toPath string, offset, limit int64) Copier {
	if fromPath == "" {
		panic(ErrFileFromNotProvided)
	}

	if toPath == "" {
		panic(ErrFileToNotProvided)
	}

	return Copier{
		pathFrom: fromPath,
		pathTo:   toPath,
		offset:   offset,
		limit:    limit,
	}
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	copier := buildCopier(fromPath, toPath, offset, limit)
	return copier.copy()
}

func (c *Copier) copy() error {
	fileFrom, err := os.OpenFile(c.pathFrom, os.O_RDWR, 0o755)
	if err != nil {
		return ErrUnsupportedFile
	}

	defer func(fileFrom *os.File) {
		_ = fileFrom.Close()
	}(fileFrom)

	fileStat, err := fileFrom.Stat()
	if err != nil {
		return err
	}

	fileSize := fileStat.Size()

	if c.offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = fileFrom.Seek(c.offset, io.SeekStart)

	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("", "tmp-")
	if err != nil {
		return err
	}

	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(tmpFile)

	readLimit := calculateReadLimit(fileSize, c.offset, c.limit)

	bar := pb.New64(readLimit)
	barReader := bar.NewProxyReader(fileFrom)
	_, err = io.CopyN(tmpFile, barReader, readLimit)

	if err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), c.pathTo)
	bar.Finish()

	if err != nil {
		return err
	}

	return nil
}

func calculateReadLimit(fileSize, offset, limit int64) int64 {
	fileSizeLeft := fileSize - offset
	readLimit := fileSizeLeft

	if limit > 0 && limit <= fileSizeLeft {
		readLimit = limit
	}

	return readLimit
}

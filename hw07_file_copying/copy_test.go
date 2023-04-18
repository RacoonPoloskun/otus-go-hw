package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const fromTest = "testdata/input.txt"

func TestCopy(t *testing.T) {
	t.Run("copy", func(t *testing.T) {
		cases := []struct {
			name   string
			offset int64
			limit  int64
			file   string
		}{
			{
				name:   "full file",
				offset: 0,
				limit:  0,
			},
			{
				name:   "low limit",
				offset: 0,
				limit:  10,
			},
			{
				name:   "big limit",
				offset: 0,
				limit:  10000,
			},
			{
				name:   "low limit low offset",
				offset: 100,
				limit:  1000,
			},
			{
				name:   "big limit big offset",
				offset: 6000,
				limit:  1000,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				tmpFile, err := os.CreateTemp("", "out.*.txt")
				if err != nil {
					fmt.Println(err)
					return
				}
				defer os.Remove(tmpFile.Name())

				copyErr := Copy(fromTest, tmpFile.Name(), c.offset, c.limit)
				require.NoError(t, copyErr, "there shouldn't be an error")

				copiedContent, err := os.ReadFile(tmpFile.Name())
				require.NotErrorIs(t, err, fs.ErrNotExist, "target file must exist")

				expectedContent, err := os.ReadFile(tmpFile.Name())
				if err != nil {
					t.Fatalf("Can't read expected file '%v'", tmpFile.Name())
				}
				require.Equal(t, expectedContent, copiedContent, "contents should be equal")
			})
		}
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		toTest, err := os.CreateTemp("", "out.*.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTest.Name())

		offsetTest := int64(7000)
		limitTest := int64(0)
		err = Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
	})

	t.Run("offset = 100, 1000 byte copied successfully", func(t *testing.T) {
		toTest, _ := os.CreateTemp("", "out.*.txt")

		offsetTest := int64(100)
		limitTest := int64(1000)
		err := Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		require.Nil(t, err)
		f1, err := os.ReadFile(toTest.Name())
		if err != nil {
			return
		}

		f2, err := os.ReadFile("testdata/out_offset100_limit1000.txt")
		if err != nil {
			return
		}

		require.True(t, bytes.Equal(f1, f2))
	})

	t.Run("end of file reached", func(t *testing.T) {
		toTest, _ := os.CreateTemp("", "out.*.txt")

		offsetTest := int64(6000)
		limitTest := int64(1000)
		err := Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		require.Nil(t, err)

		to, _ := os.Open(toTest.Name())
		defer to.Close()

		toInfo, _ := to.Stat()

		require.Nil(t, err)
		require.True(t, toInfo.Size() == 617)
	})
}

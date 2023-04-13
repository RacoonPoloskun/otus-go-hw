package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const fromTest = "testdata/input.txt"

func TestCopy(t *testing.T) {
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

	t.Run("unsupported file", func(t *testing.T) {
		toUnsupportedTest := ""
		err := Copy(fromTest, toUnsupportedTest, 0, 0)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile))

		fromTest := ""
		toTest, err := os.CreateTemp("", "out.*.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTest.Name())

		offsetTest := int64(0)
		limitTest := int64(0)
		err = Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile))

		fromTest = "/dev/urandom"
		err = Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		require.NotNil(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile))
	})

	t.Run("full copied successfully", func(t *testing.T) {
		toTest, err := os.CreateTemp("", "out.*.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(toTest.Name())

		offsetTest := int64(0)
		limitTest := int64(0)
		err = Copy(fromTest, toTest.Name(), offsetTest, limitTest)

		info, _ := toTest.Stat()

		require.Nil(t, err)
		require.True(t, info.Size() == 6617)
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

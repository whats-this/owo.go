package owo

import (
	"fmt"
	"io"
	"math"
	"os"

	"bytes"
)

// FilesToNamedReaders Converts a list of file names to named readers.
func FilesToNamedReaders(names []string) (files []NamedReader, err error) {
	if len(names) > FileCountLimit {
		err = ErrTooManyFiles{len(names)}
		return
	}
	files = make([]NamedReader, len(names))
	var file *os.File
	var stat os.FileInfo

	for idx, name := range names {
		file, err = os.Open(name)
		if err != nil {
			return
		}
		stat, err = file.Stat()
		if err != nil {
			return
		}
		if stat.Size() > FileUploadLimit {
			err = ErrFileTooBig{file.Name(), uint64(stat.Size())}
			return
		}
		var bb bytes.Buffer
		_, err = io.Copy(&bb, file)
		if err != nil {
			return
		}
		files[idx] = NamedReader{bytes.NewReader(bb.Bytes()), stat.Name()}
		err = file.Close()
		if err != nil {
			return
		}
	}
	return
}

// ErrFileTooBig thrown when file exceeds hardcoded filesize limit
type ErrFileTooBig struct {
	Filename string
	Filesize uint64
}

var sizes = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
var fileUploadLimitSize = humanateBytes(FileUploadLimit, 1024, sizes)

func (e ErrFileTooBig) Error() string {
	return fmt.Sprintf("[pre-flight] File '%s' exceeds upload limit (%s > %s)", e.Filename, humanateBytes(e.Filesize, 1024, sizes), fileUploadLimitSize)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

// ErrTooManyFiles thrown when file count in upload exceeds filecount limit
type ErrTooManyFiles struct {
	Count int
}

func (e ErrTooManyFiles) Error() string {
	return fmt.Sprintf("[pre-flight] Too many files (%d > %d)", e.Count, FileCountLimit)
}

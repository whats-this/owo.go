package owo

import (
	"fmt"
	"os"

	humanize "github.com/dustin/go-humanize"
)

// FilesToNamedReaders Converts a list of file names to named readers.
func FilesToNamedReaders(names []string) (files []NamedReader, err error) {
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
		files[idx] = NamedReader{file, name}
		defer file.Close()
	}
	return
}

// ErrFileTooBig thrown when file exceeds hardcoded filesize limit
type ErrFileTooBig struct {
	Filename string
	Filesize uint64
}

func (e ErrFileTooBig) Error() string {
	return fmt.Sprintf("[pre-flight] File '%s' exceeds upload limit (%s > %s)", e.Filename, humanize.Bytes(e.Filesize), humanize.Bytes(FileUploadLimit))
}

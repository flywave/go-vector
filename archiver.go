package govector

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mholt/archiver/v3"
)

type WalkFunc func(filename string, f io.ReadCloser, size int64) error

type Archiver struct {
	ext      string
	fileName string
	file     io.Reader
	tempDir  string
	archive  string
}

func NewArchiver(fileName string, file io.Reader) *Archiver {
	return &Archiver{fileName: fileName, file: file}
}

func (a *Archiver) Close() error {
	if a.archive != "" {
		os.Remove(a.archive)
	}
	return os.RemoveAll(a.tempDir)
}

func (a *Archiver) Valid() error {
	_, err := archiver.ByExtension(a.fileName)

	if err != nil {
		return err
	}
	return nil
}

func (a *Archiver) Walk(f WalkFunc) error {
	filename, err := a.writeArchive()

	if err != nil {
		return err
	}

	err = archiver.Walk(filename, func(file archiver.File) error {
		return f(file.FileInfo.Name(), file.ReadCloser, file.FileInfo.Size())
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *Archiver) writeArchive() (string, error) {
	if a.archive != "" && FileExists(a.archive) {
		return a.archive, nil
	}

	if FileExists(a.fileName) {
		return a.fileName, nil
	}

	name := strings.TrimSuffix(path.Base(a.fileName), a.ext)
	p, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", name, a.ext))

	if err != nil {
		return "", err
	}

	defer p.Close()

	a.archive = p.Name()

	return p.Name(), nil
}

func (a *Archiver) Unarchive() (string, error) {
	filename, err := a.writeArchive()

	if err != nil {
		return "", err
	}

	name := strings.TrimSuffix(path.Base(a.fileName), a.ext)

	p, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("%s-", name))

	if err != nil {
		return "", err
	}

	a.tempDir = p

	if !FileExists(filename) {
		return "", errors.New("not found archiver")
	}

	err = archiver.Unarchive(filename, a.tempDir)
	if err != nil {
		return "", err
	}

	return a.tempDir, nil
}

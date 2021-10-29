package govector

import (
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/flywave/go-geobuf"
	"github.com/flywave/go-geom"
)

type GeoBufProvider struct {
	reader *geobuf.Reader
}

func (p *GeoBufProvider) Open(filename string, file io.Reader) error {
	if FileExists(filename) {
		p.reader = geobuf.ReaderFile(filename)
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		p.reader = geobuf.ReaderBuf(data)
	}
	return nil
}

func (p *GeoBufProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".geobuf" {
		return false
	}
	var reader *geobuf.Reader
	if FileExists(filename) {
		reader = geobuf.ReaderFile(filename)
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return false
		}
		reader = geobuf.ReaderBuf(data)
	}
	if reader.MetaDataBool {
		return true
	}
	return false
}

func (p *GeoBufProvider) Close() error {
	p.reader.Close()
	return nil
}

func (p *GeoBufProvider) Reset() error {
	if p.reader == nil {
		return errors.New("reader is null")
	}
	p.reader.Reset()
	return nil
}

func (p *GeoBufProvider) Next() bool {
	if p.reader == nil {
		return false
	}
	return p.reader.Next()
}

func (p *GeoBufProvider) Read() *geom.Feature {
	if p.reader == nil {
		return nil
	}
	return p.reader.Feature()
}

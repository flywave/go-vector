package govector

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/flywave/go-geobuf"
	"github.com/flywave/go-geom"
)

type GeoBufProvider struct {
	reader   *geobuf.Reader
	filename string
}

func (p *GeoBufProvider) tempFileFromReader(file io.Reader) (string, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "*-geobuf")
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	if _, err := io.Copy(tmp, file); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}

func (p *GeoBufProvider) Open(filename string, file io.Reader) error {
	if FileExists(filename) {
		p.filename = filename
		p.reader = geobuf.ReaderFile(filename)
		return nil
	}
	name, err := p.tempFileFromReader(file)
	if err != nil {
		return err
	}
	p.filename = name
	p.reader = geobuf.ReaderFile(name)
	return nil
}

func (p *GeoBufProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".geobuf" {
		return false
	}
	if FileExists(filename) {
		r := geobuf.ReaderFile(filename)
		if r.Next() {
			r.Close()
			return true
		}
		r.Close()
		return false
	}
	tmp, err := p.tempFileFromReader(file)
	if err != nil {
		return false
	}
	defer os.Remove(tmp)
	r := geobuf.ReaderFile(tmp)
	if r.Next() {
		r.Close()
		return true
	}
	r.Close()
	return false
}

func (p *GeoBufProvider) Close() error {
	if p.reader != nil {
		p.reader.Close()
	}
	if p.filename != "" && !FileExists(p.filename) {
		// temp file was created from reader
		os.Remove(p.filename)
	}
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

type GeoBufExporter struct {
	out      io.WriteCloser
	filename string
	writer   *geobuf.Writer
}

func newGeoBufExporter(writer io.WriteCloser) Exporter {
	tmp, err := ioutil.TempFile(os.TempDir(), "*-export.geobuf")
	if err != nil {
		return nil
	}
	name := tmp.Name()
	tmp.Close()

	w := geobuf.WriterFile(name)
	return &GeoBufExporter{
		out:      writer,
		filename: name,
		writer:   w,
	}
}

func (e *GeoBufExporter) WriteFeature(feature *geom.Feature) error {
	e.writer.WriteFeature(feature)
	return nil
}

func (e *GeoBufExporter) WriteFeatureCollection(feature *geom.FeatureCollection) error {
	for i := range feature.Features {
		e.writer.WriteFeature(feature.Features[i])
	}
	return nil
}

func (e *GeoBufExporter) Flush() error {
	return fmt.Errorf("Flush not supported for file-backed GeoBuf exporter")
}

func (e *GeoBufExporter) Close() error {
	defer e.out.Close()
	defer os.Remove(e.filename)

	e.writer.Close()

	f, err := os.Open(e.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(e.out, f)
	return err
}

func (e *GeoBufExporter) Bytes() ([]byte, error) {
	e.writer.Close()
	defer os.Remove(e.filename)
	return ioutil.ReadFile(e.filename)
}

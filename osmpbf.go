package govector

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/flywave/go-geobuf"
	"github.com/flywave/go-geom"
	"github.com/flywave/go-osm"
)

type featureWriter struct {
	filename string
	writer   *geobuf.Writer
}

func newFeatureWriter(filename string) *featureWriter {
	return &featureWriter{writer: geobuf.WriterFile(filename), filename: filename}
}

func (w *featureWriter) WriteFeature(feature *geom.Feature) error {
	if w.writer != nil {
		w.writer.WriteFeature(feature)
	}
	return nil
}

func (w *featureWriter) Close() {
	os.Remove(w.filename)
	w.writer.Close()
}

func (w *featureWriter) Reader() *geobuf.Reader {
	return w.writer.Reader()
}

type OSMPbfProvider struct {
	decoder *osm.Decoder
	writer  *featureWriter
	reader  *geobuf.Reader
	wg      sync.WaitGroup
}

func (p *OSMPbfProvider) Open(filename string, file io.Reader) error {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(path.Base(filename), ext)

	tempPbf, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", path.Base(name), ext))

	if err != nil {
		return err
	}

	_, err = io.Copy(tempPbf, file)

	if err != nil {
		return err
	}

	tempPbf.Sync()
	tempPbf.Seek(0, io.SeekStart)

	tempBuf, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", path.Base(name), ".geobuf"))

	if err != nil {
		return err
	}

	bufFileName := tempBuf.Name()

	tempBuf.Sync()
	tempBuf.Close()

	p.writer = newFeatureWriter(bufFileName)
	p.decoder = osm.ReadDecoder(tempPbf, int(math.MaxInt64), p.writer)

	p.wg.Add(1)

	go func() {
		defer tempPbf.Close()
		p.decoder.ProcessFile()
		p.reader = p.writer.Reader()
		p.wg.Done()
	}()

	return nil
}

func (p *OSMPbfProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".pbf" {
		return false
	}
	return true
}

func (p *OSMPbfProvider) Reset() error {
	p.wg.Wait()

	if p.reader == nil {
		return errors.New("reader is null")
	}
	p.reader.Reset()
	return nil
}

func (p *OSMPbfProvider) Close() error {
	p.wg.Wait()

	if p.decoder != nil {
		p.decoder.Close()
	}
	if p.writer != nil {
		p.writer.Close()
	}
	if p.reader != nil {
		p.reader.Close()
	}
	return nil
}

func (p *OSMPbfProvider) Next() bool {
	p.wg.Wait()

	if p.reader == nil {
		return false
	}
	return p.reader.Next()
}

func (p *OSMPbfProvider) Read() *geom.Feature {
	if p.reader == nil {
		return nil
	}
	return p.reader.Feature()
}

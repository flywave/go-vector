package govector

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flywave/go-geom"
	"github.com/mholt/archiver/v3"
)

var GZ_MAGIC = []byte("\x1f\x8b")

type GeoJSONGZProvider struct {
	GeoJSONGSeqProvider
}

func (p *GeoJSONGZProvider) Open(filename string, file io.Reader) error {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(path.Base(filename), ext)

	reader, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", name, ext))

	if err != nil {
		return err
	}

	gz := archiver.NewGz()

	gz.Decompress(file, reader)

	reader.Sync()

	reader.Seek(0, io.SeekStart)

	defer reader.Close()

	jsonname := "in.json"

	return p.GeoJSONGSeqProvider.Open(jsonname, reader)
}

func (p *GeoJSONGZProvider) Match(filename string, file io.Reader) bool {
	if !strings.HasSuffix(filename, ".geojson.gz") && !strings.HasSuffix(filename, ".json.gz") {
		return false
	}
	data := make([]byte, 3)
	n, _ := file.Read(data)
	if n < 3 {
		return false
	}
	return bytes.HasPrefix(data, GZ_MAGIC)
}

type GeoJSONGZExporter struct {
	jsonExporter *GeoJSONGSeqExporter
	tempFile     *os.File
	writer       io.WriteCloser
}

func newGeoJSONGZExporter(writer io.WriteCloser) Exporter {
	tempFile, err := ioutil.TempFile(os.TempDir(), "*-export.json")
	if err != nil {
		return nil
	}
	return &GeoJSONGZExporter{
		jsonExporter: newGeoJSONGSeqExporter(tempFile).(*GeoJSONGSeqExporter),
		tempFile:     tempFile,
		writer:       writer,
	}
}

func (e *GeoJSONGZExporter) WriteFeature(feature *geom.Feature) error {
	return e.jsonExporter.WriteFeature(feature)
}

func (e *GeoJSONGZExporter) WriteFeatureCollection(feature *geom.FeatureCollection) error {
	return e.jsonExporter.WriteFeatureCollection(feature)
}

func (e *GeoJSONGZExporter) Flush() error {
	return e.jsonExporter.Flush()
}

func (e *GeoJSONGZExporter) Close() error {
	defer e.writer.Close()
	defer func() {
		e.tempFile.Close()
		os.Remove(e.tempFile.Name())
	}()

	err := e.jsonExporter.Close()
	if err != nil {
		return err
	}

	e.tempFile.Seek(0, io.SeekStart)

	gzWriter := gzip.NewWriter(e.writer)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, e.tempFile)
	return err
}

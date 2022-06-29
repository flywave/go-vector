package govector

import (
	"bytes"
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
	GeoJSONProvider
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

	return p.GeoJSONProvider.Open(jsonname, reader)
}

func (p *GeoJSONGZProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".gz" && (!strings.HasSuffix(filename, ".geojson.gz") || !strings.HasSuffix(filename, ".json.gz")) {
		return false
	}
	data := make([]byte, 3)
	file.Read(data)
	if bytes.HasPrefix(data, GZ_MAGIC) {
		return true
	}
	return false
}

type GeoJSONGZExporter struct {
	jsonExporter *GeoJSONExporter
	tempFile     *os.File
	writer       io.WriteCloser
}

func newGeoJSONGZExporter(writer io.WriteCloser) Exporter {
	tempFile, err := ioutil.TempFile(os.TempDir(), "*-export.json")
	if err != nil {
		return nil
	}
	return &GeoJSONGZExporter{jsonExporter: &GeoJSONExporter{cache: &geom.FeatureCollection{Features: make([]*geom.Feature, 0, 1024)}, writer: writer}, tempFile: tempFile, writer: writer}
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
	json, err := e.jsonExporter.cache.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = e.writer.Write(json)
	defer func() {
		os.Remove(e.tempFile.Name())
		e.tempFile.Close()
	}()

	if err != nil {
		return err
	}

	writer := archiver.NewTarGz()

	defer writer.Close()

	if err := writer.Create(e.writer); err != nil {
		return err
	}

	info, err := e.tempFile.Stat()

	if err != nil {
		return err
	}

	err = writer.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   info,
			CustomName: "export.json",
		},
		ReadCloser: e.tempFile,
	})

	if err != nil {
		return err
	}

	return nil
}

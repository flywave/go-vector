package govector

import (
	"bytes"
	"compress/gzip"
	"io"
	"path/filepath"
	"strings"

	_ "github.com/flywave/go-geom"
)

var GZ_MAGIC = []byte("\x1f\x8b")

type GeoJSONGZProvider struct {
	GeoJSONProvider
}

func (p *GeoJSONGZProvider) Open(filename string, file io.Reader) error {
	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	ext := filepath.Ext(filename)
	if ext == ".gz" {
		filename = strings.TrimSuffix(filename, ".gz")
	}
	return p.GeoJSONProvider.Open(filename, reader)
}

func (p *GeoJSONGZProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".gz" || !strings.HasSuffix(filename, ".geojson.gz") {
		return false
	}
	data := make([]byte, 3)
	file.Read(data)
	if bytes.HasPrefix(data, GZ_MAGIC) {
		return true
	}
	return false
}

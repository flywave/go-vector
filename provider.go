package govector

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/flywave/go-geom"
)

type Provider interface {
	Match(filename string, file io.Reader) bool
	Open(filename string, file io.Reader) error
	Close() error
	Reset() error
	Next() bool
	Read() *geom.Feature
}

func MatchProvider(filename string, file io.Reader) Provider {
	ext := filepath.Ext(filename)
	switch ext {
	case ".geobuf":
		p := &GeoBufProvider{}
		if p.Match(filename, file) {
			return p
		}
	case ".csv":
		p := NewGeoCSVProvider()
		if p.Match(filename, file) {
			return p
		}
	case ".geojson", ".json":
		p := &GeoJSONProvider{}
		if p.Match(filename, file) {
			return p
		} else {
			p2 := &GeoJSONGSeqProvider{}
			if p2.Match(filename, file) {
				return p2
			}
		}
	case ".gpkg":
		p := &GeoPackageProvider{}
		if p.Match(filename, file) {
			return p
		}
	case ".pbf":
		p := &OSMPbfProvider{}
		if p.Match(filename, file) {
			return p
		}
	case ".gz":
		if strings.HasSuffix(filename, ".geojson.gz") || strings.HasSuffix(filename, ".json.gz") {
			p := &GeoJSONGSeqProvider{}
			if p.Match(filename, file) {
				return p
			}
		}
		if strings.HasSuffix(filename, ".tar.gz") {
			p := &ShapeProvider{}
			if p.Match(filename, file) {
				return p
			}
		}
	case ".zip":
		p := &ShapeProvider{}
		if p.Match(filename, file) {
			return p
		}
	}
	return nil
}

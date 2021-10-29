package govector

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/flywave/go-geom"
)

type Exporter interface {
	WriteFeature(feature *geom.Feature) error
	WriteFeatureCollection(feature *geom.FeatureCollection) error
	Flush() error
	Close() error
}

func NewExporter(filename string, writer io.WriteCloser) Exporter {
	ext := filepath.Ext(filename)
	switch ext {
	case ".geobuf":
		return newGeoBufExporter(writer)
	case ".gz":
		if strings.HasSuffix(filename, ".geojson.gz") || strings.HasSuffix(filename, ".json.gz") {
			return newGeoJSONGZExporter(writer)
		}
	case ".geojson", ".json":
		return newGeoJSONGSeqExporter(writer)
	case ".gpkg":
		return newGeoPackageExporter(writer)
	}
	return nil
}

package govector

import (
	"io"
	"path/filepath"

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
		return newGeoJSONGZExporter(writer)
	case ".geojson":
		return newGeoJSONGSeqExporter(writer)
	case ".gpkg":
		return newGeoPackageExporter(writer)
	}
	return nil
}

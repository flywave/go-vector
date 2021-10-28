package govector

import (
	"io"

	"github.com/flywave/go-geom"
	_ "github.com/flywave/go-geom"
)

type GeoJSONProvider struct {
	fc *geom.FeatureCollection
}

func (p *GeoJSONProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GeoJSONProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *GeoJSONProvider) Close() error {
	return nil
}

func (p *GeoJSONProvider) Next() bool {
	return false
}

func (p *GeoJSONProvider) Read() *geom.Feature {
	return nil
}

package govector

import (
	"io"

	"github.com/flywave/go-geom"
	_ "github.com/flywave/go-geom"
)

type GeoJSONGZProvider struct {
}

func (p *GeoJSONGZProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GeoJSONGZProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *GeoJSONGZProvider) Close() error {
	return nil
}

func (p *GeoJSONGZProvider) HasNext() bool {
	return false
}

func (p *GeoJSONGZProvider) NextFeature() *geom.Feature {
	return nil
}

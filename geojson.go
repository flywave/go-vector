package govector

import (
	"io"

	"github.com/flywave/go-geom"
	_ "github.com/flywave/go-geom"
)

type GeoJSONProvider struct {
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

func (p *GeoJSONProvider) HasNext() bool {
	return false
}

func (p *GeoJSONProvider) NextFeature() *geom.Feature {
	return nil
}

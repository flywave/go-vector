package govector

import (
	"io"

	_ "github.com/flywave/go-geocsv"
	"github.com/flywave/go-geom"
)

type GeoCSVProvider struct {
}

func (p *GeoCSVProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GeoCSVProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *GeoCSVProvider) Close() error {
	return nil
}

func (p *GeoCSVProvider) HasNext() bool {
	return false
}

func (p *GeoCSVProvider) NextFeature() *geom.Feature {
	return nil
}

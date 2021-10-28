package govector

import (
	"io"

	"github.com/flywave/go-geocsv"
	"github.com/flywave/go-geom"
)

type GeoCSVProvider struct {
	csv *geocsv.GeoCSV
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

func (p *GeoCSVProvider) Next() bool {
	return false
}

func (p *GeoCSVProvider) Read() *geom.Feature {
	return nil
}

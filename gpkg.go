package govector

import (
	"io"

	"github.com/flywave/go-geom"
	_ "github.com/flywave/go-gpkg"
)

type GeoPackageProvider struct {
}

func (p *GeoPackageProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GeoPackageProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *GeoPackageProvider) Close() error {
	return nil
}

func (p *GeoPackageProvider) HasNext() bool {
	return false
}

func (p *GeoPackageProvider) NextFeature() *geom.Feature {
	return nil
}

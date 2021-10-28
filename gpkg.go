package govector

import (
	"io"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-gpkg"
)

type GeoPackageProvider struct {
	db *gpkg.GeoPackage
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

func (p *GeoPackageProvider) Next() bool {
	return false
}

func (p *GeoPackageProvider) Read() *geom.Feature {
	return nil
}

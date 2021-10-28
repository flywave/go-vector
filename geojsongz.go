package govector

import (
	"io"

	_ "github.com/flywave/go-geom"
)

type GeoJSONGZProvider struct {
	GeoJSONProvider
}

func (p *GeoJSONGZProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GeoJSONGZProvider) Match(filename string, file io.Reader) bool {
	return false
}

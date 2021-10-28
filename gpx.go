package govector

import (
	"io"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-gpx"
)

type GPXProvider struct {
	gpx *gpx.GPX
}

func (p *GPXProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GPXProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *GPXProvider) Close() error {
	return nil
}

func (p *GPXProvider) Next() bool {
	return false
}

func (p *GPXProvider) Read() *geom.Feature {
	return nil
}

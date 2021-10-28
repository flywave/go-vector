package govector

import (
	"io"

	"github.com/flywave/go-geom"
	_ "github.com/flywave/go-gpx"
)

type GPXProvider struct {
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

func (p *GPXProvider) HasNext() bool {
	return false
}

func (p *GPXProvider) NextFeature() *geom.Feature {
	return nil
}

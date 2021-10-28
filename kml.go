package govector

import (
	"io"

	"github.com/flywave/go-geom"
	_ "github.com/flywave/go-geom/kml"
)

type KMLProvider struct {
}

func (p *KMLProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *KMLProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *KMLProvider) Close() error {
	return nil
}

func (p *KMLProvider) HasNext() bool {
	return false
}

func (p *KMLProvider) NextFeature() *geom.Feature {
	return nil
}

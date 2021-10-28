package govector

import (
	"io"

	"github.com/flywave/go-geom"
	_ "github.com/flywave/go-osm"
)

type OSMPbfProvider struct {
}

func (p *OSMPbfProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *OSMPbfProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *OSMPbfProvider) Close() error {
	return nil
}

func (p *OSMPbfProvider) HasNext() bool {
	return false
}

func (p *OSMPbfProvider) NextFeature() *geom.Feature {
	return nil
}

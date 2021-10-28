package govector

import (
	"io"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-osm"
)

type OSMPbfProvider struct {
	decoder *osm.Decoder
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

func (p *OSMPbfProvider) Next() bool {
	return false
}

func (p *OSMPbfProvider) Read() *geom.Feature {
	return nil
}

package govector

import (
	"io"

	_ "github.com/flywave/go-geobuf"
	"github.com/flywave/go-geom"
)

type GeoBufProvider struct {
}

func (p *GeoBufProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GeoBufProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *GeoBufProvider) Close() error {
	return nil
}

func (p *GeoBufProvider) HasNext() bool {
	return false
}

func (p *GeoBufProvider) NextFeature() *geom.Feature {
	return nil
}

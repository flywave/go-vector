package govector

import (
	"io"

	"github.com/flywave/go-geobuf"
	"github.com/flywave/go-geom"
)

type GeoBufProvider struct {
	reader *geobuf.Reader
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

func (p *GeoBufProvider) Next() bool {
	return false
}

func (p *GeoBufProvider) Read() *geom.Feature {
	return nil
}

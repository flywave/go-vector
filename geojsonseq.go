package govector

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/flywave/go-geom"
)

const resourceSep = byte('\x1E')

type chunk struct {
	endReached bool
	reader     *bufio.Reader
	buf        []byte
	err        error
}

func (ch *chunk) next() bool {
	if ch.endReached {
		return false
	}
	var err error
	ch.buf, err = ch.reader.ReadBytes(resourceSep)
	if err == io.EOF {
		ch.endReached = true
	} else if err != nil {
		ch.err = err
		return true
	}
	ch.buf = ch.buf[:len(ch.buf)-1]
	return true
}

func (ch *chunk) read() (*geom.Feature, error) {
	if ch.err != nil {
		return nil, ch.err
	}
	var fts []*geom.Feature
	err := json.Unmarshal(append(append([]byte(`[`), ch.buf...), ']'), &fts)
	if err != nil {
		return nil, err
	}
	if len(fts) == 0 || len(fts) > 1 {
		return nil, errors.New("geojson format error")
	}
	return fts[0], nil
}

type GeoJSONGSeqProvider struct {
}

func (p *GeoJSONGSeqProvider) Open(filename string, file io.Reader) error {
	return nil
}

func (p *GeoJSONGSeqProvider) Match(filename string, file io.Reader) bool {
	return false
}

func (p *GeoJSONGSeqProvider) Close() error {
	return nil
}

func (p *GeoJSONGSeqProvider) Next() bool {
	return false
}

func (p *GeoJSONGSeqProvider) Read() *geom.Feature {
	return nil
}

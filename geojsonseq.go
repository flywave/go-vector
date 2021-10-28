package govector

import (
	"bufio"
	"encoding/json"
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

func (ch *chunk) Next() bool {
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

func (ch *chunk) Scan(fc *geom.FeatureCollection) error {
	if ch.err != nil {
		return ch.err
	}
	var fts []*geom.Feature
	err := json.Unmarshal(append(append([]byte(`[`), ch.buf...), ']'), &fts)
	if err != nil {
		return err
	}
	fc.Features = append(fc.Features, fts...)
	return nil
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

func (p *GeoJSONGSeqProvider) HasNext() bool {
	return false
}

func (p *GeoJSONGSeqProvider) NextFeature() *geom.Feature {
	return nil
}

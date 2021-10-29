package govector

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"

	"github.com/flywave/go-geom"
)

const resourceSep = byte('\x1E')

type chunk struct {
	endReached bool
	reader     *bufio.Reader
	buf        []byte
	err        error
}

func newChunk(rd io.Reader) *chunk {
	return &chunk{endReached: false, reader: bufio.NewReader(rd), buf: nil, err: nil}
}

func (ch *chunk) reset(r io.Reader) {
	ch.reader.Reset(r)
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
	chunk *chunk
	file  io.Reader
}

func (p *GeoJSONGSeqProvider) Open(filename string, file io.Reader) error {
	p.file = file
	p.chunk = newChunk(file)
	return nil
}

func (p *GeoJSONGSeqProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".geojson" && ext != ".json" {
		return false
	}
	buffer := bufio.NewReader(file)

	buf, err := buffer.ReadBytes(resourceSep)
	if err != nil {
		return false
	}
	var fts []*geom.Feature
	err = json.Unmarshal(append(append([]byte(`[`), buf...), ']'), &fts)
	if err != nil {
		return false
	}
	if len(fts) == 0 || len(fts) > 1 {
		return false
	}
	return true
}

func (p *GeoJSONGSeqProvider) Reset() error {
	p.chunk.reset(p.file)
	return nil
}

func (p *GeoJSONGSeqProvider) Close() error {
	return nil
}

func (p *GeoJSONGSeqProvider) Next() bool {
	if p.chunk != nil {
		return p.chunk.next()
	}
	return false
}

func (p *GeoJSONGSeqProvider) Read() *geom.Feature {
	feat, err := p.chunk.read()
	if err != nil {
		return nil
	}
	return feat
}

type GeoJSONGSeqExporter struct {
	writer   io.WriteCloser
	bufwrite *bufio.Writer
}

func newGeoJSONGSeqExporter(writer io.WriteCloser) Exporter {
	return &GeoJSONGSeqExporter{writer: writer, bufwrite: bufio.NewWriter(writer)}
}

func (e *GeoJSONGSeqExporter) WriteFeature(feature *geom.Feature) error {
	json, err := feature.MarshalJSON()
	if err != nil {
		return err
	}
	jsonline := append([]byte(json), resourceSep)
	_, err = e.bufwrite.Write(jsonline)
	return err
}

func (e *GeoJSONGSeqExporter) WriteFeatureCollection(feature *geom.FeatureCollection) error {
	for i := range feature.Features {
		err := e.WriteFeature(feature.Features[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *GeoJSONGSeqExporter) Flush() error {
	return e.bufwrite.Flush()
}

func (e *GeoJSONGSeqExporter) Close() error {
	e.Flush()
	return e.writer.Close()
}

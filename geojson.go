package govector

import (
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

type GeoJSONProvider struct {
	fc    *geom.FeatureCollection
	index int
}

func (p *GeoJSONProvider) Open(filename string, file io.Reader) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	p.fc, err = general.UnmarshalFeatureCollection(data)
	if err != nil {
		return err
	}
	return nil
}

func (p *GeoJSONProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".geojson" && ext != ".json" {
		return false
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	p.fc, err = general.UnmarshalFeatureCollection(data)
	if err != nil {
		return false
	}
	return len(p.fc.Features) > 0
}

func (p *GeoJSONProvider) Reset() error {
	p.index = 0
	return nil
}

func (p *GeoJSONProvider) Close() error {
	return nil
}

func (p *GeoJSONProvider) Next() bool {
	if p.index < (len(p.fc.Features)) {
		p.index++
		return true
	}
	return false
}

func (p *GeoJSONProvider) Read() *geom.Feature {
	if p.fc == nil {
		return nil
	}
	return p.fc.Features[p.index-1]
}

type GeoJSONExporter struct {
	cache  *geom.FeatureCollection
	writer io.WriteCloser
}

func newGeoJSONExporter(writer io.WriteCloser) Exporter {
	return &GeoJSONExporter{cache: &geom.FeatureCollection{Features: make([]*geom.Feature, 0, 1024)}, writer: writer}
}

func (e *GeoJSONExporter) WriteFeature(feature *geom.Feature) error {
	if e.cache == nil {
		return errors.New("export not init")
	}
	e.cache.Features = append(e.cache.Features, feature)
	return nil
}

func (e *GeoJSONExporter) WriteFeatureCollection(feature *geom.FeatureCollection) error {
	if e.cache == nil {
		return errors.New("export not init")
	}
	e.cache.Features = append(e.cache.Features, feature.Features...)
	return nil
}

func (e *GeoJSONExporter) Flush() error {
	return nil
}

func (e *GeoJSONExporter) Close() error {
	json, err := e.cache.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = e.writer.Write(json)
	return e.writer.Close()
}

package govector

import (
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
	if p.index < (len(p.fc.Features) - 1) {
		p.index++
		return true
	}
	return false
}

func (p *GeoJSONProvider) Read() *geom.Feature {
	if p.fc == nil {
		return nil
	}
	return p.fc.Features[p.index]
}

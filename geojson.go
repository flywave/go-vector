package govector

import (
	"encoding/json"
	"errors"
	"io"
	"path/filepath"

	"github.com/flywave/go-geom"
)

type GeoJSONProvider struct {
	dec     *json.Decoder
	current *geom.Feature
	done    bool
	err     error
}

func (p *GeoJSONProvider) Open(filename string, file io.Reader) error {
	p.dec = json.NewDecoder(file)
	p.done = false
	p.err = nil

	t, err := p.dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return errors.New("GeoJSON must start with '{'")
	}

	for p.dec.More() {
		t, err = p.dec.Token()
		if err != nil {
			return err
		}
		key, ok := t.(string)
		if !ok {
			return errors.New("expected string key in GeoJSON object")
		}
		if key == "features" {
			t, err = p.dec.Token()
			if err != nil {
				return err
			}
			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				return errors.New("expected '[' after 'features'")
			}
			return nil
		}
		var raw json.RawMessage
		if err := p.dec.Decode(&raw); err != nil {
			return err
		}
	}
	return errors.New("no 'features' array found in GeoJSON")
}

func (p *GeoJSONProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".geojson" && ext != ".json" {
		return false
	}
	dec := json.NewDecoder(file)
	t, err := dec.Token()
	if err != nil {
		return false
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return false
	}
	hasFeatures := false
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return false
		}
		key, ok := t.(string)
		if !ok {
			return false
		}
		if key == "type" {
			var typeStr string
			if err := dec.Decode(&typeStr); err != nil {
				return false
			}
			if typeStr != "FeatureCollection" {
				return false
			}
		} else if key == "features" {
			t, err = dec.Token()
			if err != nil {
				return false
			}
			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				return false
			}
			hasFeatures = dec.More()
			break
		} else {
			var raw json.RawMessage
			if err := dec.Decode(&raw); err != nil {
				return false
			}
		}
	}
	return hasFeatures
}

func (p *GeoJSONProvider) Reset() error {
	return errors.New("Reset not supported for streaming GeoJSON provider")
}

func (p *GeoJSONProvider) Close() error {
	p.dec = nil
	p.current = nil
	p.done = true
	return nil
}

func (p *GeoJSONProvider) Next() bool {
	if p.err != nil || p.done || p.dec == nil {
		return false
	}
	if !p.dec.More() {
		p.done = true
		return false
	}
	var feat geom.Feature
	if err := p.dec.Decode(&feat); err != nil {
		p.err = err
		return false
	}
	p.current = &feat
	return true
}

func (p *GeoJSONProvider) Read() *geom.Feature {
	return p.current
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

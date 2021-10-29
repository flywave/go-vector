package govector

import (
	"io"
	"path/filepath"

	"github.com/flywave/go-geocsv"
	"github.com/flywave/go-geom"
)

type GeoCSVProvider struct {
	csv     *geocsv.GeoCSV
	options geocsv.GeoCSVOptions
	index   int
}

func NewGeoCSVProvider(options geocsv.GeoCSVOptions) Provider {
	return &GeoCSVProvider{csv: nil, options: options, index: 0}
}

func (p *GeoCSVProvider) Open(filename string, file io.Reader) error {
	var err error
	p.csv, err = geocsv.Read(file, p.options)
	if err != nil {
		return err
	}
	return nil
}

func (p *GeoCSVProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".csv" {
		return false
	}
	csv, err := geocsv.Read(file, p.options)
	if err != nil {
		return false
	}
	return csv.Valid()
}

func (p *GeoCSVProvider) Reset() error {
	p.index = 0
	return nil
}

func (p *GeoCSVProvider) Close() error {
	return nil
}

func (p *GeoCSVProvider) Next() bool {
	if p.index < (p.csv.RowCount() - 1) {
		p.index++
		return true
	}
	return false
}

func (p *GeoCSVProvider) Read() *geom.Feature {
	if p.csv == nil {
		return nil
	}
	return p.csv.Feature(p.index)
}

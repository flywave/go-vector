package govector

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-gpkg"
)

type GeoPackageProvider struct {
	dbFileName string
	db         *gpkg.GeoPackage
	layers     []gpkg.VectorLayer
	index      int
	current    *gpkg.GeoPackageReader
}

func (p *GeoPackageProvider) Open(filename string, file io.Reader) error {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(path.Base(filename), ext)

	tempGpkg, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", path.Base(name), ext))

	if err != nil {
		return err
	}

	dbFileName := tempGpkg.Name()

	tempGpkg.Sync()
	tempGpkg.Close()

	p.db = gpkg.New(dbFileName)
	p.db.Init()
	p.dbFileName = dbFileName

	p.layers, err = p.db.GetVectorLayers()

	if err != nil {
		return err
	}

	return nil
}

func (p *GeoPackageProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".gpkg" {
		return false
	}

	name := strings.TrimSuffix(path.Base(filename), ext)

	tempGpkg, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", path.Base(name), ext))

	if err != nil {
		return false
	}

	dbFileName := tempGpkg.Name()

	tempGpkg.Sync()
	tempGpkg.Close()

	defer os.Remove(dbFileName)

	db := gpkg.New(dbFileName)
	db.Init()

	layers, err := db.GetVectorLayers()

	if err != nil {
		return false
	}

	return len(layers) > 0
}

func (p *GeoPackageProvider) Reset() error {
	if p.current != nil {
		p.current = nil
	}
	p.index = 0
	return nil
}

func (p *GeoPackageProvider) Close() error {
	p.Reset()
	os.Remove(p.dbFileName)
	return p.db.Close()
}

func (p *GeoPackageProvider) Next() bool {
	if p.current == nil && p.index == 0 {
		l := p.layers[p.index]
		var err error
		p.current, err = p.db.GetFeatureReader(l.Name)
		if err != nil {
			return false
		}
	}
	if p.current.Next() {
		return true
	} else if p.index < len(p.layers)-1 {
		p.index++
		l := p.layers[p.index]
		var err error
		p.current, err = p.db.GetFeatureReader(l.Name)
		if err != nil {
			return false
		}
		return p.current.Next()
	}

	return false
}

func (p *GeoPackageProvider) Read() *geom.Feature {
	if p.current != nil {
		feat, err := p.current.Read()
		if err != nil {
			return nil
		}
		return feat
	}
	return nil
}

const EXPORT_TABLE_NAME = "export"

type GeoPackageExporter struct {
	dbFileName string
	db         *gpkg.GeoPackage
	writer     io.WriteCloser
	cache      *geom.FeatureCollection
}

func newGeoPackageExporter(writer io.WriteCloser) Exporter {
	tempGpkg, err := ioutil.TempFile(os.TempDir(), "*-export.gpkg")

	if err != nil {
		return nil
	}

	dbFileName := tempGpkg.Name()

	tempGpkg.Close()

	return &GeoPackageExporter{writer: writer, db: gpkg.Create(dbFileName), dbFileName: dbFileName, cache: &geom.FeatureCollection{Features: make([]*geom.Feature, 0, 1024)}}
}

func (e *GeoPackageExporter) WriteFeature(feature *geom.Feature) error {
	if e.cache == nil {
		return errors.New("export not init")
	}
	if e.cache != nil {
		e.cache.Features = append(e.cache.Features, feature)
	}
	return nil
}

func (e *GeoPackageExporter) WriteFeatureCollection(feature *geom.FeatureCollection) error {
	if e.cache == nil {
		return errors.New("export not init")
	}
	e.cache.Features = append(e.cache.Features, feature.Features...)
	return nil
}

func (e *GeoPackageExporter) Flush() error {
	return nil
}

func (e *GeoPackageExporter) Close() error {
	defer e.writer.Close()
	defer os.Remove(e.dbFileName)
	err := e.db.StoreFeatureCollection(EXPORT_TABLE_NAME, e.cache)
	if err != nil {
		return err
	}
	e.db.Close()
	f, err := os.Open(e.dbFileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(e.writer, f)
	if err != nil {
		return err
	}
	return nil
}

package govector

import (
	"bytes"
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

	io.Copy(tempGpkg, file)

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

var gpkgMagic = []byte("SQLite format 3\x00")

func (p *GeoPackageProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".gpkg" {
		return false
	}

	// Quick check: GPKG is SQLite, must start with magic header
	buf := make([]byte, 16)
	if _, err := io.ReadFull(file, buf); err != nil {
		return false
	}
	if !bytes.HasPrefix(buf, gpkgMagic) {
		return false
	}

	// For full validation, copy to temp and check layers
	return p.matchFull(filename, io.MultiReader(bytes.NewReader(buf), file))
}

func (p *GeoPackageProvider) matchFull(filename string, file io.Reader) bool {
	name := strings.TrimSuffix(path.Base(filename), filepath.Ext(filename))
	ext := filepath.Ext(filename)

	tempGpkg, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", path.Base(name), ext))
	if err != nil {
		return false
	}
	defer os.Remove(tempGpkg.Name())

	io.Copy(tempGpkg, file)
	tempGpkg.Sync()
	tempGpkg.Close()

	db := gpkg.New(tempGpkg.Name())
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
	for {
		if p.current == nil {
			if p.index >= len(p.layers) {
				return false
			}
			l := p.layers[p.index]
			var err error
			p.current, err = p.db.GetFeatureReader(l.Name)
			if err != nil {
				return false
			}
		}
		if p.current.Next() {
			return true
		}
		p.current = nil
		p.index++
	}
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
const EXPORT_FLUSH_BATCH = 500

type GeoPackageExporter struct {
	dbFileName    string
	db            *gpkg.GeoPackage
	writer        io.WriteCloser
	batch         []*geom.Feature
	tableReady    bool
}

func newGeoPackageExporter(writer io.WriteCloser) Exporter {
	tempGpkg, err := ioutil.TempFile(os.TempDir(), "*-export.gpkg")

	if err != nil {
		return nil
	}

	dbFileName := tempGpkg.Name()

	tempGpkg.Close()

	return &GeoPackageExporter{
		writer:     writer,
		db:         gpkg.Create(dbFileName),
		dbFileName: dbFileName,
		batch:      make([]*geom.Feature, 0, EXPORT_FLUSH_BATCH),
	}
}

func (e *GeoPackageExporter) flushBatch() error {
	if len(e.batch) == 0 {
		return nil
	}
	fc := &geom.FeatureCollection{Features: e.batch}
	if err := e.db.StoreFeatureCollection(EXPORT_TABLE_NAME, fc); err != nil {
		return err
	}
	e.batch = e.batch[:0]
	e.tableReady = true
	return nil
}

func (e *GeoPackageExporter) WriteFeature(feature *geom.Feature) error {
	if e.batch == nil {
		return errors.New("export not init")
	}
	e.batch = append(e.batch, feature)
	if len(e.batch) >= EXPORT_FLUSH_BATCH {
		return e.flushBatch()
	}
	return nil
}

func (e *GeoPackageExporter) WriteFeatureCollection(feature *geom.FeatureCollection) error {
	for _, f := range feature.Features {
		if err := e.WriteFeature(f); err != nil {
			return err
		}
	}
	return nil
}

func (e *GeoPackageExporter) Flush() error {
	return e.flushBatch()
}

func (e *GeoPackageExporter) Close() error {
	defer e.writer.Close()
	defer os.Remove(e.dbFileName)

	if err := e.flushBatch(); err != nil {
		return err
	}

	e.db.Close()
	f, err := os.Open(e.dbFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(e.writer, f)
	return err
}

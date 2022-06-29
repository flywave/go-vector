package govector

import (
	"errors"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-shp"
)

const (
	SHAPE_SHX_EXT = ".shx"
	SHAPE_SHP_EXT = ".shp"
	SHAPE_DBF_EXT = ".dbf"
)

type shapeIterator struct {
	file        *shp.ShapeFile
	currentFeat int
}

func newIterator(path string) *shapeIterator {
	return &shapeIterator{file: shp.Open(path), currentFeat: 0}
}

func (p *shapeIterator) next() bool {
	if p.file == nil || p.file.ShapeCount == 0 {
		return false
	}
	if p.currentFeat < (p.file.ShapeCount) {
		p.currentFeat++
		return true
	}
	return false
}

func (p *shapeIterator) readFeature() *geom.Feature {
	i := p.currentFeat
	return p.file.Feature(i)
}

func (p *shapeIterator) Close() error {
	p.file.Close()
	return nil
}

type ShapeProvider struct {
	archiver *Archiver
	shpfiles []string
	workDir  string
	index    int
	current  *shapeIterator
}

func (p *ShapeProvider) ShapeFiles() []string {
	return p.shpfiles
}

func (p *ShapeProvider) Open(filename string, file io.Reader) error {
	arch := NewArchiver(filename, file)

	if err := arch.Valid(); err != nil {
		return err
	}
	p.archiver = arch

	shpfiles := make(map[string]string)

	arch.Walk(func(filename string, f io.ReadCloser, size int64) error {
		ext := filepath.Ext(filename)
		if ext == SHAPE_SHP_EXT {
			shpfiles[filename] = filename
		}
		return nil
	})

	for _, f := range shpfiles {
		p.shpfiles = append(p.shpfiles, f)
	}

	return nil
}

func (p *ShapeProvider) Match(filename string, file io.Reader) bool {
	arch := NewArchiver(filename, file)

	if err := arch.Valid(); err != nil {
		return false
	}

	type shpFile struct {
		shx bool
		shp bool
		dbf bool
	}

	shpfiles := make(map[string]*shpFile)

	arch.Walk(func(filename string, f io.ReadCloser, size int64) error {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(path.Base(filename), ext)
		if _, ok := shpfiles[name]; !ok {
			shpfiles[name] = &shpFile{shx: false, shp: false, dbf: false}
		}
		if ext == SHAPE_SHX_EXT {
			shpfiles[name].shx = true
		}
		if ext == SHAPE_SHP_EXT {
			shpfiles[name].shp = true
		}
		if ext == SHAPE_DBF_EXT {
			shpfiles[name].dbf = true
		}
		return nil
	})

	if len(shpfiles) == 0 {
		return false
	}

	for _, f := range shpfiles {
		if !f.dbf || !f.shp || !f.shx {
			return false
		}
	}

	return true
}

func (p *ShapeProvider) Reset() error {
	if p.current != nil {
		p.current.Close()
	}
	p.current = nil
	if p.moveNext() {
		return nil
	}
	return errors.New("reset error")
}

func (p *ShapeProvider) Close() error {
	if p.current != nil {
		p.current.Close()
	}
	return p.archiver.Close()
}

func (p *ShapeProvider) moveNext() bool {
	if len(p.shpfiles) == 0 {
		return false
	}
	if p.index < len(p.shpfiles)-1 {
		p.index++
		filename := p.shpfiles[p.index]
		filename = path.Join(p.workDir, filename)
		if p.current != nil {
			p.current.Close()
		}
		p.current = newIterator(filename)
		return true
	}
	return false
}

func (p *ShapeProvider) Next() bool {
	if p.current != nil {
		if p.current.next() {
			return true
		} else {
			return p.moveNext()
		}
	}
	return false
}

func (p *ShapeProvider) Read() *geom.Feature {
	if p.current != nil {
		return p.current.readFeature()
	}
	return nil
}

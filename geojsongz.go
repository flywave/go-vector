package govector

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

var GZ_MAGIC = []byte("\x1f\x8b")

type GeoJSONGZProvider struct {
	GeoJSONProvider
}

func (p *GeoJSONGZProvider) Open(filename string, file io.Reader) error {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(path.Base(filename), ext)

	frr, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*-%s%s", name, ext))

	if err != nil {
		return err
	}

	defer frr.Close()

	var reader io.Reader
	var jsonname string

	err = archiver.Walk(frr.Name(), func(file archiver.File) error {
		jext := filepath.Ext(file.FileInfo.Name())

		if jext == ".json" {
			jsonname = path.Base(file.FileInfo.Name())
			reader = file.ReadCloser
		}

		return nil
	})
	if err != nil {
		return err
	}

	return p.GeoJSONProvider.Open(jsonname, reader)
}

func (p *GeoJSONGZProvider) Match(filename string, file io.Reader) bool {
	ext := filepath.Ext(filename)
	if ext != ".gz" && (!strings.HasSuffix(filename, ".geojson.gz") || !strings.HasSuffix(filename, ".json.gz")) {
		return false
	}
	data := make([]byte, 3)
	file.Read(data)
	if bytes.HasPrefix(data, GZ_MAGIC) {
		return true
	}
	return false
}

package govector

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/flywave/go-geom/general"
)

func TestGeoBufExporter(t *testing.T) {
	f, _ := os.Open("./data.json")

	data, _ := ioutil.ReadAll(f)

	fcs, _ := general.UnmarshalFeatureCollection(data)

	fileName := "./test.geobuf"
	f, err := os.Open(fileName)
	if err != nil {
		t.FailNow()
	}
	export := NewExporter(fileName, f)

	err = export.WriteFeatureCollection(fcs)
	if err != nil {
		t.FailNow()
	}

	err = export.Close()
	if err != nil {
		t.FailNow()
	}

	//os.Remove("./test.geobuf")
}

func TestGeoJsonGzExporter(t *testing.T) {
	f, _ := os.Open("./data.json")

	data, _ := ioutil.ReadAll(f)

	fcs, _ := general.UnmarshalFeatureCollection(data)

	fileName := "./test.geojson.gz"
	f, err := os.Open(fileName)
	if err != nil {
		t.FailNow()
	}
	export := NewExporter(fileName, f)

	err = export.WriteFeatureCollection(fcs)
	if err != nil {
		t.FailNow()
	}

	err = export.Close()
	if err != nil {
		t.FailNow()
	}

	//os.Remove("./test.geojson.gz")
}

func TestGeoJsonSeqExporter(t *testing.T) {
	f, _ := os.Open("./data.json")

	data, _ := ioutil.ReadAll(f)

	fcs, _ := general.UnmarshalFeatureCollection(data)

	fileName := "./test.geojson"
	f, err := os.Open(fileName)
	if err != nil {
		t.FailNow()
	}
	export := NewExporter(fileName, f)

	err = export.WriteFeatureCollection(fcs)
	if err != nil {
		t.FailNow()
	}

	err = export.Close()
	if err != nil {
		t.FailNow()
	}

	//os.Remove("./test.geojson")
}

func TestGPKGExporter(t *testing.T) {
	f, _ := os.Open("./data.json")

	data, _ := ioutil.ReadAll(f)

	fcs, _ := general.UnmarshalFeatureCollection(data)

	fileName := "./test.gpkg"
	f, err := os.Open(fileName)
	if err != nil {
		t.FailNow()
	}
	export := NewExporter(fileName, f)

	err = export.WriteFeatureCollection(fcs)
	if err != nil {
		t.FailNow()
	}

	err = export.Close()
	if err != nil {
		t.FailNow()
	}

	//os.Remove("./test.gpkg")
}

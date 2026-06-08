package govector

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

func readTestFeatures(t *testing.T) []*geom.Feature {
	t.Helper()
	data, err := ioutil.ReadFile("./testdata/data.json")
	if err != nil {
		t.Fatal(err)
	}
	fc, err := general.UnmarshalFeatureCollection(data)
	if err != nil {
		t.Fatal(err)
	}
	return fc.Features
}

func TestGeoBufExporter(t *testing.T) {
	features := readTestFeatures(t)

	outFile, err := os.CreateTemp(t.TempDir(), "*.geobuf")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	export := NewExporter(outFile.Name(), outFile)
	if export == nil {
		t.Fatal("NewExporter returned nil for .geobuf")
	}

	fc := &geom.FeatureCollection{Features: features}
	if err := export.WriteFeatureCollection(fc); err != nil {
		t.Fatal(err)
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}

	// Round-trip: read back
	outFile.Seek(0, io.SeekStart)
	provider := MatchProvider(outFile.Name(), outFile)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .geobuf round-trip")
	}
	outFile.Seek(0, io.SeekStart)
	if err := provider.Open(outFile.Name(), outFile); err != nil {
		t.Fatal(err)
	}
	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	provider.Close()
	if count == 0 {
		t.Fatal("round-trip: no features read back")
	}
	if count != len(features) {
		t.Fatalf("round-trip: got %d features, expected %d", count, len(features))
	}
}

func TestGeoBufExporter_WriteFeature(t *testing.T) {
	features := readTestFeatures(t)

	var buf bytes.Buffer
	export := NewExporter("test.geobuf", &nopWriteCloser{&buf})
	if export == nil {
		t.Fatal("NewExporter returned nil for .geobuf")
	}
	for _, f := range features {
		if err := export.WriteFeature(f); err != nil {
			t.Fatal(err)
		}
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestGeoJsonGzExporter(t *testing.T) {
	features := readTestFeatures(t)

	outFile, err := os.CreateTemp(t.TempDir(), "*.geojson.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	export := NewExporter(outFile.Name(), outFile)
	if export == nil {
		t.Fatal("NewExporter returned nil for .geojson.gz")
	}

	fc := &geom.FeatureCollection{Features: features}
	if err := export.WriteFeatureCollection(fc); err != nil {
		t.Fatal(err)
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}

	// Round-trip: read back
	outFile.Seek(0, io.SeekStart)
	provider := MatchProvider(outFile.Name(), outFile)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .geojson.gz round-trip")
	}
	outFile.Seek(0, io.SeekStart)
	if err := provider.Open(outFile.Name(), outFile); err != nil {
		t.Fatal(err)
	}
	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	provider.Close()
	if count == 0 {
		t.Fatal("round-trip: no features read back")
	}
}

func TestGeoJsonSeqExporter(t *testing.T) {
	features := readTestFeatures(t)

	outFile, err := os.CreateTemp(t.TempDir(), "*.geojson")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	export := NewExporter(outFile.Name(), outFile)
	if export == nil {
		t.Fatal("NewExporter returned nil for .geojson")
	}

	fc := &geom.FeatureCollection{Features: features}
	if err := export.WriteFeatureCollection(fc); err != nil {
		t.Fatal(err)
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}

	// Round-trip: read back
	outFile.Seek(0, io.SeekStart)
	provider := MatchProvider(outFile.Name(), outFile)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .geojson round-trip")
	}
	outFile.Seek(0, io.SeekStart)
	if err := provider.Open(outFile.Name(), outFile); err != nil {
		t.Fatal(err)
	}
	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	provider.Close()
	if count == 0 {
		t.Fatal("round-trip: no features read back")
	}
}

func TestGeoJsonGzExporter_WriteFeature(t *testing.T) {
	features := readTestFeatures(t)

	var buf bytes.Buffer
	export := NewExporter("test.geojson.gz", &nopWriteCloser{&buf})
	if export == nil {
		t.Fatal("NewExporter returned nil for .geojson.gz")
	}
	for _, f := range features {
		if err := export.WriteFeature(f); err != nil {
			t.Fatal(err)
		}
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestGeoJsonSeqExporter_WriteFeature(t *testing.T) {
	features := readTestFeatures(t)

	var buf bytes.Buffer
	export := NewExporter("test.geojson", &nopWriteCloser{&buf})
	if export == nil {
		t.Fatal("NewExporter returned nil for .geojson")
	}
	for _, f := range features {
		if err := export.WriteFeature(f); err != nil {
			t.Fatal(err)
		}
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestGPKGExporter(t *testing.T) {
	features := readTestFeatures(t)

	outFile, err := os.CreateTemp(t.TempDir(), "*.gpkg")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	export := NewExporter(outFile.Name(), outFile)
	if export == nil {
		t.Fatal("NewExporter returned nil for .gpkg")
	}

	fc := &geom.FeatureCollection{Features: features}
	if err := export.WriteFeatureCollection(fc); err != nil {
		t.Fatal(err)
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}

	// Round-trip: read back
	outFile.Seek(0, io.SeekStart)
	provider := MatchProvider(outFile.Name(), outFile)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .gpkg round-trip")
	}
	outFile.Seek(0, io.SeekStart)
	if err := provider.Open(outFile.Name(), outFile); err != nil {
		t.Fatal(err)
	}
	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	provider.Close()
	if count == 0 {
		t.Fatal("round-trip: no features read back")
	}
}

func TestGPKGExporter_WriteFeature(t *testing.T) {
	features := readTestFeatures(t)

	var buf bytes.Buffer
	export := NewExporter("test.gpkg", &nopWriteCloser{&buf})
	if export == nil {
		t.Fatal("NewExporter returned nil for .gpkg")
	}
	for _, f := range features {
		if err := export.WriteFeature(f); err != nil {
			t.Fatal(err)
		}
	}
	if err := export.Close(); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestNewExporter_UnsupportedExt(t *testing.T) {
	var buf bytes.Buffer
	export := NewExporter("test.csv", &nopWriteCloser{&buf})
	if export != nil {
		t.Fatal("expected nil for unsupported extension")
	}
	export = NewExporter("test.shp", &nopWriteCloser{&buf})
	if export != nil {
		t.Fatal("expected nil for unsupported extension")
	}
}

// nopWriteCloser wraps a bytes.Buffer to implement io.WriteCloser
type nopWriteCloser struct {
	io.Writer
}

func (n *nopWriteCloser) Close() error { return nil }

var _ io.WriteCloser = (*nopWriteCloser)(nil)

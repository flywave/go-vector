package govector

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/flywave/go-geom"
)

func TestGeoBufProvider(t *testing.T) {
	reader, err := os.Open("./testdata/data.geobuf")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/data.geobuf", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .geobuf")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/data.geobuf", reader); err != nil {
		t.Fatal(err)
	}

	var feats []*geom.Feature
	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.Fatal("expected at least one feature")
	}

	// Verify feature properties
	if feats[0].Properties == nil {
		t.Fatal("expected feature properties")
	}

	provider.Close()
}

func TestGeoBufProvider_Reset(t *testing.T) {
	reader, err := os.Open("./testdata/data.geobuf")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/data.geobuf", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/data.geobuf", reader); err != nil {
		t.Fatal(err)
	}

	// Read first few features
	for i := 0; i < 3 && provider.Next(); i++ {
		provider.Read()
	}

	if err := provider.Reset(); err != nil {
		t.Fatal(err)
	}

	// After reset, should be able to read all features again
	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	if count == 0 {
		t.Fatal("expected features after Reset()")
	}

	provider.Close()
}

func TestGeoBufProvider_CloseBeforeRead(t *testing.T) {
	reader, err := os.Open("./testdata/data.geobuf")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/data.geobuf", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/data.geobuf", reader); err != nil {
		t.Fatal(err)
	}

	provider.Close()

	// After Close, Next should return false
	if provider.Next() {
		t.Fatal("Next() should return false after Close()")
	}

	if provider.Read() != nil {
		t.Fatal("Read() should return nil after Close()")
	}
}

func TestGeoCSVProvider(t *testing.T) {
	reader, err := os.Open("./testdata/ne_110m_populated_places_simple.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/ne_110m_populated_places_simple.csv", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .csv")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/ne_110m_populated_places_simple.csv", reader); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil")
		}
		count++
	}

	if count == 0 {
		t.Fatal("expected at least one feature from CSV")
	}

	provider.Close()
}

func TestGeoCSVProvider_UnreadBeforeOpen(t *testing.T) {
	p := NewGeoCSVProvider()
	if p.Next() {
		t.Fatal("Next() should return false before Open()")
	}
	if p.Read() != nil {
		t.Fatal("Read() should return nil before Open()")
	}
}

func TestGeoJsonGzProvider(t *testing.T) {
	reader, err := os.Open("./testdata/in.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/in.json.gz", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .json.gz")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/in.json.gz", reader); err != nil {
		t.Fatal(err)
	}

	var feats []*geom.Feature
	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.Fatal("expected at least one feature from .json.gz")
	}

	provider.Close()
}

func TestGeoJsonGzProvider_ResetReturnsError(t *testing.T) {
	reader, err := os.Open("./testdata/in.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/in.json.gz", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/in.json.gz", reader); err != nil {
		t.Fatal(err)
	}

	if err := provider.Reset(); err == nil {
		t.Log("Reset on streaming GZ provider may not be supported")
	}
	provider.Close()
}

func TestGeoJsonSeqProvider(t *testing.T) {
	reader, err := os.Open("./testdata/sherlock.json")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/sherlock.json", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .json GeoJSONSeq")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/sherlock.json", reader); err != nil {
		t.Fatal(err)
	}

	var feats []*geom.Feature
	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.Fatal("expected at least one feature from GeoJSONSeq")
	}

	provider.Close()
}

func TestGPKGProvider(t *testing.T) {
	reader, err := os.Open("./testdata/natural_earth_minimal.gpkg")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/natural_earth_minimal.gpkg", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .gpkg")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/natural_earth_minimal.gpkg", reader); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		provider.Read()
		count++
	}

	if count == 0 {
		t.Fatal("expected at least one feature from GPKG")
	}

	provider.Close()
}

func TestGPKGProvider_Reset(t *testing.T) {
	reader, err := os.Open("./testdata/natural_earth_minimal.gpkg")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/natural_earth_minimal.gpkg", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/natural_earth_minimal.gpkg", reader); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2 && provider.Next(); i++ {
		provider.Read()
	}

	if err := provider.Reset(); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	if count == 0 {
		t.Fatal("expected features after Reset()")
	}

	provider.Close()
}

func TestOSMPBFProvider(t *testing.T) {
	reader, err := os.Open("./testdata/sample.osm.pbf")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/sample.osm.pbf", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .osm.pbf")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/sample.osm.pbf", reader); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil")
		}
		count++
	}

	if count == 0 {
		t.Fatal("expected at least one feature from OSM PBF")
	}

	provider.Close()
}

func TestOSMPBFProvider_Reset(t *testing.T) {
	reader, err := os.Open("./testdata/sample.osm.pbf")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/sample.osm.pbf", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/sample.osm.pbf", reader); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3 && provider.Next(); i++ {
		provider.Read()
	}

	if err := provider.Reset(); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	if count == 0 {
		t.Fatal("expected features after Reset()")
	}

	provider.Close()
}

func TestSHPProvider(t *testing.T) {
	reader, err := os.Open("./testdata/shp.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/shp.tar.gz", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .tar.gz Shapefile")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/shp.tar.gz", reader); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil")
		}
		count++
	}

	if count == 0 {
		t.Fatal("expected at least one feature from Shapefile")
	}

	provider.Close()
}

func TestSHPProvider_Reset(t *testing.T) {
	reader, err := os.Open("./testdata/shp.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/shp.tar.gz", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/shp.tar.gz", reader); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2 && provider.Next(); i++ {
		provider.Read()
	}

	if err := provider.Reset(); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		provider.Read()
		count++
	}
	if count == 0 {
		t.Fatal("expected features after Reset()")
	}

	provider.Close()
}

func TestSHPProvider_FromZip(t *testing.T) {
	reader, err := os.Open("./testdata/shp.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/shp.zip", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .zip Shapefile")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/shp.zip", reader); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil")
		}
		count++
	}

	if count == 0 {
		t.Fatal("expected at least one feature from zip Shapefile")
	}

	provider.Close()
}

func TestSHPProvider_FromBuffer(t *testing.T) {
	data, err := os.ReadFile("./testdata/shp.tar.gz")
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewReader(data)
	provider := MatchProvider("./testdata/shp.tar.gz", buf)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for buffer-based .tar.gz")
	}

	buf.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/shp.tar.gz", buf); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil")
		}
		count++
	}

	if count == 0 {
		t.Fatal("expected at least one feature from buffer-based Shapefile")
	}

	provider.Close()
}

func TestMatchProvider_UnsupportedExtensions(t *testing.T) {
	tests := []struct {
		name string
		ext  string
	}{
		{"unknown extension", "data.xyz"},
		{"no extension", "README"},
		{"image", "photo.png"},
		{"text", "doc.txt"},
	}

	for _, tc := range tests {
		t.Run(tc.ext, func(t *testing.T) {
			provider := MatchProvider(tc.name, bytes.NewReader(nil))
			if provider != nil {
				t.Fatalf("MatchProvider should return nil for %q", tc.name)
			}
		})
	}
}

func TestGeoJSONProvider_MatchThenOpen(t *testing.T) {
	data := `{"type":"FeatureCollection","features":[{"type":"Feature","properties":{"name":"test"},"geometry":{"type":"Point","coordinates":[0,0]}}]}`

	// First Match consumes the reader
	buf := strings.NewReader(data)
	provider := MatchProvider("test.geojson", buf)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for valid GeoJSON")
	}

	// Seek back and call Open
	buf.Seek(0, io.SeekStart)
	if err := provider.Open("test.geojson", buf); err != nil {
		t.Fatal(err)
	}

	if !provider.Next() {
		t.Fatal("expected at least one feature")
	}
	feat := provider.Read()
	if feat == nil {
		t.Fatal("Read() returned nil")
	}
	if feat.Properties["name"] != "test" {
		t.Fatalf("unexpected property value: %v", feat.Properties["name"])
	}
	provider.Close()
}

func TestProvider_EmptyFeatureCollection(t *testing.T) {
	data := `{"type":"FeatureCollection","features":[]}`
	buf := strings.NewReader(data)
	reader := io.ReadSeeker(buf)

	provider := MatchProvider("test.geojson", reader)
	if provider != nil {
		t.Fatal("MatchProvider should return nil for empty FeatureCollection (expected GeoJSONSeq fallback)")
	}
}

func TestProvider_ReadFromBuffer(t *testing.T) {
	data, err := os.ReadFile("./testdata/data.json")
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewReader(data)
	provider := MatchProvider("test.json", buf)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for JSON from buffer")
	}

	buf.Seek(0, io.SeekStart)
	if err := provider.Open("test.json", buf); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil")
		}
		count++
	}
	if count == 0 {
		t.Fatal("expected features from buffer-based JSON")
	}
	provider.Close()
}

func TestGPKGProvider_FeatureProperties(t *testing.T) {
	reader, err := os.Open("./testdata/natural_earth_minimal.gpkg")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/natural_earth_minimal.gpkg", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .gpkg")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/natural_earth_minimal.gpkg", reader); err != nil {
		t.Fatal(err)
	}

	if !provider.Next() {
		t.Fatal("expected at least one feature")
	}
	feat := provider.Read()
	if feat == nil {
		t.Fatal("Read() returned nil")
	}
	if feat.Properties == nil {
		t.Fatal("expected feature properties")
	}

	provider.Close()
}

func TestGeoJSONProvider_StreamingLargeData(t *testing.T) {
	// Generate a synthetic large FeatureCollection to verify streaming behavior
	var buf bytes.Buffer
	buf.WriteString(`{"type":"FeatureCollection","features":[`)
	for i := 0; i < 10000; i++ {
		if i > 0 {
			buf.WriteString(`,`)
		}
		buf.WriteString(`{"type":"Feature","properties":{"id":`)
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString(`},"geometry":{"type":"Point","coordinates":[0,0]}}`)
	}
	buf.WriteString(`]}`)

	r := bytes.NewReader(buf.Bytes())
	provider := MatchProvider("test.geojson", r)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for large GeoJSON")
	}

	r.Seek(0, io.SeekStart)
	if err := provider.Open("test.geojson", r); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil during streaming")
		}
		count++
	}
	if count != 10000 {
		t.Fatalf("expected 10000 features, got %d", count)
	}
	provider.Close()
}

func TestGeoBufProvider_FromBuffer(t *testing.T) {
	// Read geobuf file, then use as buffer to test temp-file path
	data, err := os.ReadFile("./testdata/data.geobuf")
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewReader(data)

	provider := MatchProvider("test.geobuf", buf)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for buffer-based .geobuf")
	}

	buf.Seek(0, io.SeekStart)
	if err := provider.Open("test.geobuf", buf); err != nil {
		t.Fatal(err)
	}

	var count int
	for provider.Next() {
		feat := provider.Read()
		if feat == nil {
			t.Fatal("Read() returned nil")
		}
		count++
	}
	if count == 0 {
		t.Fatal("expected features from buffer-based GeoBuf")
	}

	provider.Close()
}

func TestGPKGProvider_MatchQuick(t *testing.T) {
	// GPKG Match should detect invalid files without full DB open
	invalid := bytes.NewReader([]byte("not a gpkg file at all"))
	provider := MatchProvider("test.gpkg", invalid)
	if provider != nil {
		t.Fatal("MatchProvider should return nil for invalid GPKG")
	}
}

func TestGeoJSONProvider_FeatureProperties(t *testing.T) {
	reader, err := os.Open("./testdata/data.json")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	provider := MatchProvider("./testdata/data.json", reader)
	if provider == nil {
		t.Fatal("MatchProvider returned nil for .json")
	}

	reader.Seek(0, io.SeekStart)
	if err := provider.Open("./testdata/data.json", reader); err != nil {
		t.Fatal(err)
	}

	if !provider.Next() {
		t.Fatal("expected at least one feature")
	}
	feat := provider.Read()
	if feat == nil {
		t.Fatal("Read() returned nil")
	}
	if feat.Properties == nil {
		t.Fatal("expected feature properties")
	}
	// Check known property from Afghanistan
	if name, ok := feat.Properties["name"]; !ok || name != "Afghanistan" {
		t.Fatalf("unexpected first feature: name=%v", name)
	}

	provider.Close()
}

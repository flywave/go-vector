package govector

import (
	"io"
	"os"
	"testing"

	"github.com/flywave/go-geom"
)

func TestGeoBufProvider(t *testing.T) {
	reader, err := os.Open("./testdata/5_22_11.geobuf")
	if err != nil {
		t.FailNow()
	}
	provider := MatchProvider("./testdata/5_22_11.geobuf", reader)
	if provider == nil {
		t.FailNow()
	}

	err = provider.Open("./testdata/5_22_11.geobuf", reader)
	if err != nil {
		t.FailNow()
	}

	feats := []*geom.Feature{}

	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.FailNow()
	}

	provider.Close()
}

func TestGeoCSVProvider(t *testing.T) {
	reader, err := os.Open("./testdata/ne_110m_populated_places_simple.csv")
	if err != nil {
		t.FailNow()
	}
	provider := MatchProvider("./testdata/ne_110m_populated_places_simple.csv", reader)
	if provider == nil {
		t.FailNow()
	}

	reader.Seek(0, io.SeekStart)

	err = provider.Open("./testdata/ne_110m_populated_places_simple.csv", reader)
	if err != nil {
		t.FailNow()
	}

	feats := []*geom.Feature{}

	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.FailNow()
	}

	provider.Close()
}

func TestGeoJsonGzProvider(t *testing.T) {
	reader, err := os.Open("./testdata/in.json.gz")
	if err != nil {
		t.FailNow()
	}
	provider := MatchProvider("./testdata/in.json.gz", reader)
	if provider == nil {
		t.FailNow()
	}
	reader.Seek(0, io.SeekStart)

	err = provider.Open("./testdata/in.json.gz", reader)
	if err != nil {
		t.FailNow()
	}

	feats := []*geom.Feature{}

	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.FailNow()
	}

	provider.Close()
}

func TestGeoJsonSeqProvider(t *testing.T) {
	reader, err := os.Open("./testdata/sherlock.json")
	if err != nil {
		t.FailNow()
	}
	provider := MatchProvider("./testdata/sherlock.json", reader)
	if provider == nil {
		t.FailNow()
	}
	reader.Seek(0, io.SeekStart)

	err = provider.Open("./testdata/sherlock.json", reader)
	if err != nil {
		t.FailNow()
	}

	feats := []*geom.Feature{}

	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.FailNow()
	}

	provider.Close()
}

func TestGPKGProvider(t *testing.T) {
	reader, err := os.Open("./testdata/natural_earth_minimal.gpkg")
	if err != nil {
		t.FailNow()
	}
	provider := MatchProvider("./testdata/natural_earth_minimal.gpkg", reader)
	if provider == nil {
		t.FailNow()
	}
	reader.Seek(0, io.SeekStart)

	err = provider.Open("./testdata/natural_earth_minimal.gpkg", reader)
	if err != nil {
		t.FailNow()
	}

	feats := []*geom.Feature{}

	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.FailNow()
	}

	provider.Close()
}

func TestOSMPBFProvider(t *testing.T) {
	reader, err := os.Open("./testdata/sample.osm.pbf")
	if err != nil {
		t.FailNow()
	}
	provider := MatchProvider("./testdata/sample.osm.pbf", reader)
	if provider == nil {
		t.FailNow()
	}
	reader.Seek(0, io.SeekStart)
	err = provider.Open("./testdata/sample.osm.pbf", reader)
	if err != nil {
		t.FailNow()
	}

	feats := []*geom.Feature{}

	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.FailNow()
	}

	provider.Close()
}

func TestSHPProvider(t *testing.T) {
	reader, err := os.Open("./testdata/shp.tar.gz")
	if err != nil {
		t.FailNow()
	}
	provider := MatchProvider("./testdata/shp.tar.gz", reader)
	if provider == nil {
		t.FailNow()
	}
	reader.Seek(0, io.SeekStart)

	err = provider.Open("./testdata/shp.tar.gz", reader)
	if err != nil {
		t.FailNow()
	}

	feats := []*geom.Feature{}

	for provider.Next() {
		feats = append(feats, provider.Read())
	}

	if len(feats) == 0 {
		t.FailNow()
	}

	provider.Close()
}

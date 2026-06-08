# go-vector

[![Go Reference](https://pkg.go.dev/badge/github.com/flywave/go-vector.svg)](https://pkg.go.dev/github.com/flywave/go-vector)

Single-package Go library providing a uniform `Provider` (read) and `Exporter` (write) interface for geospatial vector formats.

- **Go 1.16+**, no sub-packages, no CLI binary
- All data flows through `io.Reader` / `io.WriteCloser` — no direct file coupling
- Streaming-first design: memory O(1) per feature for most formats
- Temp files for archive/gzip/GPKG/OSM PBF processing, cleaned up on `Close()`

## Supported formats

| Format | Read | Write | Streaming |
|--------|------|-------|-----------|
| GeoJSON (`.geojson` / `.json`) | ✅ | ✅ GeoJSONSeq | ✅ read |
| GeoJSONSeq (`.geojson` / `.json`) | ✅ | ✅ | ✅ |
| GeoJSON GZ (`.geojson.gz` / `.json.gz`) | ✅ | ✅ | ✅ read |
| GeoBuf (`.geobuf`) | ✅ | ✅ | ✅ read |
| CSV (`.csv`) | ✅ | — | depends |
| GeoPackage (`.gpkg`) | ✅ | ✅ | ✅ read, ✅ batch write |
| OSM PBF (`.osm.pbf`) | ✅ | — | ✅ (async) |
| Shapefile (`.tar.gz` / `.zip`) | ✅ | — | ✅ (file-backed) |

## Quick start

```go
package main

import (
    "bytes"
    "fmt"
    "os"

    govector "github.com/flywave/go-vector"
)

func main() {
    // --- Read ---
    f, _ := os.Open("data.geojson")
    defer f.Close()

    p := govector.MatchProvider("data.geojson", f)
    f.Seek(0, 0)
    p.Open("data.geojson", f)

    for p.Next() {
        feat := p.Read()
        fmt.Println(feat.Properties["name"])
    }
    p.Close()

    // --- Write ---
    var buf bytes.Buffer
    e := govector.NewExporter("out.geojson", &nopCloser{&buf})
    // write features...
    e.Close()
}

type nopCloser struct{ *bytes.Buffer }
func (n *nopCloser) Close() error { return nil }
```

## Provider interface

```go
type Provider interface {
    Match(filename string, file io.Reader) bool
    Open(filename string, file io.Reader) error
    Close() error
    Reset() error
    Next() bool
    Read() *geom.Feature
}
```

Dispatch is extension-based via `MatchProvider(filename, io.ReadSeeker)`:

- `.geobuf` → `GeoBufProvider`
- `.csv` → `GeoCSVProvider`
- `.geojson` / `.json` → `GeoJSONProvider` (FeatureCollection), fallback to `GeoJSONGSeqProvider`
- `.gpkg` → `GeoPackageProvider`
- `.pbf` → `OSMPbfProvider` (requires `.osm.pbf`)
- `.geojson.gz` / `.json.gz` → `GeoJSONGZProvider`
- `.tar.gz` / `.zip` → `ShapeProvider`

## Exporter interface

```go
type Exporter interface {
    WriteFeature(feature *geom.Feature) error
    WriteFeatureCollection(feature *geom.FeatureCollection) error
    Flush() error
    Close() error
}
```

Dispatch via `NewExporter(filename, io.WriteCloser)`:

- `.geobuf` → `GeoBufExporter`
- `.geojson.gz` / `.json.gz` → `GeoJSONGZExporter`
- `.geojson` / `.json` → `GeoJSONGSeqExporter`
- `.gpkg` → `GeoPackageExporter`

## Large data

All providers and exporters avoid loading the full dataset into memory:

- **GeoJSON** — `json.Decoder` stream-parses features; no `ReadAll` or full `FeatureCollection`
- **GeoJSON GZ** — decompresses to temp file, then streams via GSeqProvider
- **GeoBuf** — buffer path writes to temp file first, then file-backed read
- **GeoBuf/GZ export** — streams features to temp file, `io.Copy` on Close
- **GPKG export** — batches features (500/batch), flushes incrementally
- **GPKG/OSM PBF/SHP read** — copies to temp file on `Open`, then streams

## Format notes

- **CSV**: expects `longitude`/`latitude` or `WKT` columns
- **GeoJSON**: parsed as `FeatureCollection` first; empty arrays fallback to GeoJSONSeq (line-delimited)
- **GeoJSON GZ**: gzip-compressed GeoJSONSeq (`.json.gz`)
- **OSM PBF**: decoded asynchronously via goroutine, writing to a temp `.geobuf` intermediary
- **Shapefile**: archive must contain `.shp`+`.shx`+`.dbf` triples for each shape

## Dependencies

- C libraries required at link time (transitive from `go-shp`, `go-proj`, `go-geos`, `go-geoid`):
  `geos`, `proj`, `geographiclib`
- Go dependency `go-geobuf` requires a patched `go-geom` BoundingBox adapter.
  A local clone at `../go-geobuf` with `ToGeomBBox` is expected via `replace` directive.

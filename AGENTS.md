# go-vector

Single-package Go library (`github.com/flywave/go-vector`, Go 1.16) providing a uniform `Provider` (read) and `Exporter` (write) interface for geospatial vector formats.

## Commands

```sh
go build ./...          # build library (no tests)
go test ./...           # all tests (requires C libraries)
go test -run TestSHPProvider   # single test
go test -v ./...        # verbose
```

## Architecture

- **`provider.go`** — `Provider` interface + `MatchProvider(filename, io.ReadSeeker)` dispatcher. Ext-based dispatch: `.geobuf`, `.csv`, `.geojson`/`.json`, `.gpkg`, `.pbf`, `.gz` (`.geojson.gz`/`.json.gz`), `.tar.gz`/`.zip` (Shapefile).
- **`exporter.go`** — `Exporter` interface + `NewExporter(filename, io.WriteCloser)` factory. Supports `.geobuf`, `.geojson.gz`/`.json.gz`, `.geojson`/`.json`, `.gpkg`.
- All providers/exporters in flat package (`go vector/`). No sub-packages, no CLI binary.
- CSV expects `longitude`/`latitude` or `WKT` columns.
- `.geojson` files are parsed as `FeatureCollection` first; if empty, fallback to GeoJSONSeq (line-delimited).
- `.json` same as `.geojson`.
- `.geojson.gz`/`.json.gz` — gunzip then parse as GeoJSONSeq.
- `.tar.gz`/`.zip` — extract archive, find `.shp`+`.shx`+`.dbf` triples.
- OSM PBF (`.osm.pbf`) decodes asynchronously (goroutine), writing to a temp `.geobuf` intermediary.
- Read pattern: `Open` → loop `Next()`/`Read()` → `Close()`.

## Known bugs (fixed)

- `provider.go` `.geojson` fallback to GeoJSONSeq was missing `file.Seek(0)` after `GeoJSONProvider.Match` consumed the reader.
- `shp.go` `Close()` used `os.ReadFile` instead of `os.Remove` — temp files were never cleaned up.
- `shp.go` `Open()` ignored `ioutil.TempDir` error and `arch.Walk` error.
- `gpkg.go` `Next()` could panic with index-out-of-bounds on the last layer; didn't skip empty layers.
- `geojsongz.go` `GeoJSONGZExporter` wrote raw JSON bytes then a tar.gz to the same writer; now writes GeoJSONSeq to temp file then gzips.
- `archiver.go` `writeArchive()` never copied `a.file` reader content to the temp file.
- `geocsv.go` `Next()` lacked nil-check for `p.csv`.

## Dependency issue

`go-geobuf` requires an older `go-geom` where `BoundingBox = []float64`, but `go-shp` requires a newer `go-geom` where `BoundingBox = [2][3]float64`. Fixed via patched `../go-geobuf/` with `ToGeomBBox` adapter in `io/bbox.go`.

## Large data optimizations

### Streaming (no full file in memory)
- **GeoJSONProvider** — uses `json.Decoder` stream-parsing; no ReadAll, no full FeatureCollection in memory
- **GeoJSONGSeqProvider** — bufio line-by-line (already streaming)
- **GeoJSONGZProvider** — delegates to GSeqProvider, not GeoJSONProvider
- **GeoBufProvider** — from-reader path writes to temp file then opens file-backed (no ReadAll)
- **GeoBufExporter** — uses `geobuf.WriterFile` (temp file), streams features to disk; Close copies via `io.Copy`

### Batch-flush (bounded memory)
- **GeoPackageExporter** — batches features (500/batch), calls `StoreFeatureCollection` per batch; no unbounded FeatureCollection cache

### Match lightweight
- **GPKG Match** — checks SQLite magic header first before full DB open
- **GeoJSON Match** — uses json.Decoder to parse only as far as the `features` key; no full parse

### Still bounded (temp files on disk, not memory)
- GPKG/OSM PBF/SHP providers copy to temp files on `Open()`, then stream reads
- All temp files cleaned up on `Close()`

## Testing

- Tests live in the same package (`package govector`).
- Test fixtures in `testdata/`. All tests are integration-style (read real files).
- Exporter tests use `os.CreateTemp(t.TempDir())` for automatic cleanup and round-trip verification.
- Provider tests include streaming large-data validation (10K synthetic features).
- Some tests (`OSMPBF`, `GPKG`, `SHP`, `GeoBuf`) depend on C libraries (geos, proj, geoid) via transitive dependencies. These require the C libraries to be installed for linking.
- No CI, no Makefile, no lint config, no formatter config.

## Temp files

Many providers write to `os.TempDir()` for archive extraction, gzip decompression, GPKG/OSM PBF processing. The `Archiver` type handles archive extraction to temp dirs, cleaned up on `Close()`. All exporters use temp files for intermediate storage, streaming to the final output on `Close()`.

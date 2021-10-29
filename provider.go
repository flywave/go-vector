package govector

import (
	"io"

	"github.com/flywave/go-geom"
)

type Provider interface {
	Match(filename string, file io.Reader) bool
	Open(filename string, file io.Reader) error
	Close() error
	Reset() error
	Next() bool
	Read() *geom.Feature
}

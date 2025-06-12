// Package internal ...
package internal

import (
	"context"
)

// DirectoryReaderInterface public interface
type DirectoryReaderInterface interface {
	Read(context.Context, chan<- string) error
	ShowStats()
	AddState(uint64, uint64)
}

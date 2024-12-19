package tools

import (
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Fallbacks represents a collection of fallback sources for the tool.
// It can either be a single source type or a slice of source types, allowing flexibility
// in specifying multiple fallback methods if the primary source fails.
type Fallbacks = unmarshal.SingleOrSliceType[sources.Type]

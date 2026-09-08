// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/prometheus/prometheus/tsdb/chunkenc"

	"github.com/thanos-io/thanos/pkg/store/storepb"
)

func TestChunkEncoding_XOR2(t *testing.T) {
	// XOR2 (chunkenc.EncXOR2) is the Prometheus 3.11+ float chunk encoding used to
	// persist real start timestamps. The querier must be able to decode it back from
	// the storepb wire encoding after a store-gateway serves it.
	testutil.Equals(t, chunkenc.EncXOR2, chunkEncoding(storepb.Chunk_XOR2))
}

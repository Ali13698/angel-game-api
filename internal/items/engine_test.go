package items

import (
	"testing"
)

func TestMapToWasmItems(t *testing.T) {
	inv := []InfiniteItem{
		{ID: "cue1", Type: ItemTypeCue, Meta: map[string]any{"power": 1.5}},
	}
	out := MapToWasmItems(inv)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

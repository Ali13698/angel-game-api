package wasm

import (
	"fmt"
)

// Phase 1: Version is the observed state revision, NOT a protocol version
func validateCommandVersion(commandVersion int64, currentRev int64) error {
	if commandVersion != currentRev {
		return fmt.Errorf("state revision mismatch: command version=%d, current rev=%d (resync required)", commandVersion, currentRev)
	}
	return nil
}

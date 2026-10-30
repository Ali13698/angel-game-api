package wasm

import "fmt"

func validateCommandVersion(command, current uint64) error {
	if command == 0 || command != current {
		return fmt.Errorf("wasm: state revision mismatch: got %d, want %d; resync required", command, current)
	}
	return nil
}

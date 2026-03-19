package protocol_common

import "time"

const (
	// DefaultOperationTimeout is the per-attempt timeout for most LinMot operations.
	// Each retry attempt waits this long for a response before trying again.
	DefaultOperationTimeout = 200 * time.Millisecond

	// FlashOperationTimeout is the per-attempt timeout for flash save operations.
	// Flash saves take 75+ seconds on production LinMot drives. The timeout must
	// be long enough for the FIRST attempt to succeed without retrying, because
	// resending SaveCommandTable during an in-progress flash save restarts the
	// operation on the drive, preventing it from ever completing.
	FlashOperationTimeout = 120 * time.Second
)

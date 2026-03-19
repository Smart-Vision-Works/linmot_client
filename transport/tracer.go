package transport

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	tracerOnce sync.Once
	tracerFile *os.File
	tracerMu   sync.Mutex
	tracerPath string
)

// initTracer lazily initializes the packet tracer if LINMOT_DEBUG=1.
// Thread-safe: uses sync.Once to ensure single initialization.
func initTracer() {
	tracerOnce.Do(func() {
		if os.Getenv("LINMOT_DEBUG") != "1" {
			return
		}

		pid := os.Getpid()
		tracerPath = filepath.Join("/tmp", fmt.Sprintf("linmot_debug_trace_%d.log", pid))

		var err error
		tracerFile, err = os.OpenFile(tracerPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// Silently fail - don't break tests if file can't be created
			return
		}

		// Print exactly one line to stdout
		fmt.Printf("[LINMOT_DEBUG] trace file: %s\n", tracerPath)
	})
}

// TracePacket writes a packet trace line to the debug file (if LINMOT_DEBUG=1).
// dir: "TX" or "RX"
// local: local address (e.g., "127.0.0.1:41136")
// remote: remote address (e.g., "10.8.7.232:49360")
// data: packet bytes
// Thread-safe: uses mutex around file writes.
func TracePacket(dir, local, remote string, data []byte) {
	initTracer()

	tracerMu.Lock()
	defer tracerMu.Unlock()

	if tracerFile == nil {
		return
	}

	now := time.Now().UnixNano()
	hexStr := hex.EncodeToString(data)
	line := fmt.Sprintf("%d %s %s -> %s len=%d hex=%s\n", now, dir, local, remote, len(data), hexStr)

	_, err := tracerFile.WriteString(line)
	if err == nil {
		tracerFile.Sync() // Flush after each write
	}
	// Silently ignore write errors - don't break tests
}

// CloseTracer closes the trace file (called during cleanup).
func CloseTracer() {
	tracerMu.Lock()
	defer tracerMu.Unlock()

	if tracerFile != nil {
		tracerFile.Close()
		tracerFile = nil
	}
}

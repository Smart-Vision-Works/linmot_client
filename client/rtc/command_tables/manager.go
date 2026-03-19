package client_command_tables

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	client_common "github.com/Smart-Vision-Works/linmot_client/client/common"
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
	protocol_command_tables "github.com/Smart-Vision-Works/linmot_client/protocol/rtc/command_tables"
)

// ErrCommandTableUnchanged is returned by SetCommandTableWithOptions when the
// current entries match the desired entries byte-for-byte. No StopMC, write,
// flash save, or StartMC was performed. Callers should skip recovery/homing.
var ErrCommandTableUnchanged = errors.New("command table unchanged")

type CommandTableManager struct {
	requestManager *client_common.RequestManager
	debug          atomic.Bool
}

func NewCommandTableManager(requestManager *client_common.RequestManager) *CommandTableManager {
	return &CommandTableManager{
		requestManager: requestManager,
	}
}

// SetDebug enables or disables debug logging for command table operations.
func (manager *CommandTableManager) SetDebug(enabled bool) {
	manager.debug.Store(enabled)
}

// GetCommandTable retrieves the current command table from the drive.
func (manager *CommandTableManager) GetCommandTable(ctx context.Context) (*CommandTable, error) {
	masks, err := manager.getPresenceMasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get presence masks: %w", err)
	}

	entryIDs := manager.extractEntryIDsFromMasks(masks)

	entries := make([]Entry, 0, len(entryIDs))
	for _, id := range entryIDs {
		data, err := manager.readFullEntry(ctx, id) // No start time for GetCommandTable
		if err != nil {
			return nil, fmt.Errorf("failed to read entry %d: %w", id, err)
		}

		wireEntry, err := protocol_command_tables.DecodeEntry(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode entry %d: %w", id, err)
		}

		entry := fromWireEntry(*wireEntry)
		entry.ID = int(id)
		entries = append(entries, entry)
	}

	return &CommandTable{
		Version:    "",
		DriveModel: "",
		Entries:    entries,
	}, nil
}

// SetCommandTableOptions configures behavior for SetCommandTableWithOptions.
type SetCommandTableOptions struct {
	// RestartMC controls whether the motion controller is restarted after setting the command table.
	// If false, MC remains stopped (useful for tests that need to read back command table while MC is stopped).
	// Default: true (MC is restarted)
	RestartMC bool

	// SkipFlashSave skips the flash save step after writing entries to RAM.
	// When true, entries are written and verified in RAM only — the caller is
	// responsible for persisting to flash (e.g., via SaveCommandTableToFlash)
	// and handling the drive's recovery afterward.
	//
	// This is necessary because the LinMot drive does not send a standard RTC
	// response for SaveCommandTable — the response comes via cyclic status fields
	// that our request/response protocol layer cannot observe. Letting the caller
	// manage the flash save avoids an unavoidable timeout in the RTC layer.
	SkipFlashSave bool
}

// SetCommandTable sets the command table on the drive from a CommandTable struct.
// This is a convenience method that calls SetCommandTableWithOptions with RestartMC=true.
func (manager *CommandTableManager) SetCommandTable(ctx context.Context, ct *CommandTable) error {
	err := manager.SetCommandTableWithOptions(ctx, ct, SetCommandTableOptions{RestartMC: true})
	if errors.Is(err, ErrCommandTableUnchanged) {
		return nil // Unchanged is success for callers that don't need to distinguish
	}
	return err
}

// SetCommandTableWithOptions sets the command table on the drive from a CommandTable struct with configurable options.
func (manager *CommandTableManager) SetCommandTableWithOptions(ctx context.Context, ct *CommandTable, opts SetCommandTableOptions) error {
	startTime := time.Now() // Track start time for debug logging

	if err := ct.Validate(); err != nil {
		return fmt.Errorf("command table validation failed: %w", err)
	}

	// Encode all new entries to binary for comparison and writing
	minSize := manager.getMinEntrySize()
	newPayloads := make(map[uint16][]byte, len(ct.Entries))
	for _, entry := range ct.Entries {
		id := uint16(entry.ID)
		payload, err := entry.Encode(ct)
		if err != nil {
			return fmt.Errorf("failed to encode entry %d: %w", id, err)
		}
		if uint16(len(payload)) < minSize {
			padding := make([]byte, minSize-uint16(len(payload)))
			payload = append(payload, padding...)
		}
		newPayloads[id] = payload
	}

	// Read current entries from drive and compare with new entries.
	// If all entries match byte-for-byte, skip the entire deploy (no StopMC,
	// no write, no flash save, no StartMC). This avoids the ~39-second flash
	// save disruption AND the MC restart that clears RAM entries.
	needsWrite := false
	currentCT, readErr := manager.GetCommandTable(ctx)
	if readErr != nil {
		log.Printf("[RTC] Can't read current command table (%v) — full write needed", readErr)
		needsWrite = true
	} else if len(currentCT.Entries) != len(ct.Entries) {
		log.Printf("[RTC] Entry count changed (%d → %d) — full write needed",
			len(currentCT.Entries), len(ct.Entries))
		needsWrite = true
	} else {
		// Compare each entry's binary representation
		for _, currentEntry := range currentCT.Entries {
			id := uint16(currentEntry.ID)
			currentPayload, err := currentEntry.Encode(currentCT)
			if err != nil {
				log.Printf("[RTC] Can't encode current entry %d for comparison — full write needed", id)
				needsWrite = true
				break
			}
			if uint16(len(currentPayload)) < minSize {
				padding := make([]byte, minSize-uint16(len(currentPayload)))
				currentPayload = append(currentPayload, padding...)
			}
			newPayload, exists := newPayloads[id]
			if !exists {
				log.Printf("[RTC] Entry %d not in new table — full write needed", id)
				needsWrite = true
				break
			}
			if len(currentPayload) != len(newPayload) {
				log.Printf("[RTC] Entry %d size changed — full write needed", id)
				needsWrite = true
				break
			}
			for i := range currentPayload {
				if currentPayload[i] != newPayload[i] {
					log.Printf("[RTC] Entry %d differs at byte %d — full write needed", id, i)
					needsWrite = true
					break
				}
			}
			if needsWrite {
				break
			}
		}
	}

	if !needsWrite {
		log.Printf("[RTC] Command table unchanged — no action needed (%d entries verified in %dms)",
			len(ct.Entries), time.Since(startTime).Milliseconds())
		return ErrCommandTableUnchanged
	}

	// Entries differ — stop MC, write, flash save, restart MC.
	log.Printf("[RTC] Stopping motion controller (%d entries to write)", len(ct.Entries))
	if err := manager.stopMotionController(ctx); err != nil {
		return fmt.Errorf("failed to stop motion controller: %w", err)
	}

	// Conditionally restart MC before returning, even if context is cancelled.
	// Use background context with timeout to avoid failing cleanup.
	if opts.RestartMC {
		defer func() {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			log.Printf("[RTC] Restarting motion controller")
			if err := manager.startMotionController(cleanupCtx); err != nil {
				log.Printf("[WARN] Failed to restart motion controller during cleanup: %v", err)
			} else {
				log.Printf("[RTC] Motion controller restarted successfully (total deploy time: %dms)", time.Since(startTime).Milliseconds())
			}
		}()
	}

	// Entries differ — full write + verify + flash save
	if err := manager.deleteAllEntries(ctx); err != nil {
		return fmt.Errorf("failed to delete entries: %w", err)
	}

	for _, entry := range ct.Entries {
		id := uint16(entry.ID)
		payload := newPayloads[id]

		if err := manager.allocateEntry(ctx, id, uint16(len(payload))); err != nil {
			return fmt.Errorf("failed to allocate entry %d: %w", id, err)
		}

		for i := 0; i < len(payload); i += 4 {
			end := i + 4
			if end > len(payload) {
				end = len(payload)
			}
			if err := manager.writeEntryData(ctx, id, payload[i:end], i); err != nil {
				var writeErr *protocol_command_tables.EntryWriteError
				if errors.As(err, &writeErr) {
					return &protocol_command_tables.EntryWriteError{
						EntryID: writeErr.EntryID,
						Offset:  i,
						Status:  writeErr.Status,
					}
				}
				return fmt.Errorf("failed to write entry %d data at offset %d: %w", id, i, err)
			}
		}

		readBack, err := manager.readFullEntry(ctx, id, startTime)
		if err != nil {
			return fmt.Errorf("failed to read back entry %d: %w", id, err)
		}

		limit := len(payload)
		if len(readBack) < limit {
			limit = len(readBack)
		}
		for i := 0; i < limit; i++ {
			if readBack[i] != payload[i] {
				return fmt.Errorf("verification failed for entry %d: readback mismatch at byte %d", id, i)
			}
		}
	}

	log.Printf("[RTC] All %d entries written and verified in %dms", len(ct.Entries), time.Since(startTime).Milliseconds())

	if opts.SkipFlashSave {
		log.Printf("[RTC] Skipping flash save (caller will handle persistence)")
		return nil
	}

	// Save verified entries to flash so they persist across MC restart.
	// This takes ~39 seconds and temporarily disrupts the LinUDP interface.
	// The caller's recovery wait handles reconnection after the disruption.
	log.Printf("[RTC] Saving %d verified entries to flash (MC is stopped)...", len(ct.Entries))
	if err := manager.saveToFlash(ctx); err != nil {
		return fmt.Errorf("failed to save command table to flash: %w", err)
	}
	log.Printf("[RTC] Command table saved to flash successfully (%d entries)", len(ct.Entries))

	return nil
}

// StopMotionController stops the Motion Controller software on the drive.
func (manager *CommandTableManager) StopMotionController(ctx context.Context) error {
	return manager.stopMotionController(ctx)
}

// saveToFlash saves the current RAM command table to flash.
// MC must be stopped before calling this (enforced by caller).
func (manager *CommandTableManager) saveToFlash(ctx context.Context) error {
	request := protocol_command_tables.NewSaveCommandTableRequest()
	return manager.requestManager.SendRequest(ctx, request)
}

// SaveCommandTableToFlash sends the SaveCommandTable (0x80) RTC command.
// This is exposed for callers that need to manage the flash save lifecycle
// themselves (e.g., when using SkipFlashSave). The LinMot drive may not send
// a standard RTC response for this command, so callers should expect a timeout
// and verify completion by reading back the command table after recovery.
func (manager *CommandTableManager) SaveCommandTableToFlash(ctx context.Context) error {
	return manager.saveToFlash(ctx)
}

// stopMotionController stops the Motion Controller software on the drive.
// This must be called before command table operations.
func (manager *CommandTableManager) stopMotionController(ctx context.Context) error {
	request := protocol_command_tables.NewStopMotionControllerRequest()
	return manager.requestManager.SendRequest(ctx, request)
}

// StartMotionController restarts the Motion Controller software on the drive.
func (manager *CommandTableManager) StartMotionController(ctx context.Context) error {
	return manager.startMotionController(ctx)
}

// startMotionController restarts the Motion Controller software on the drive.
func (manager *CommandTableManager) startMotionController(ctx context.Context) error {
	request := protocol_command_tables.NewStartMotionControllerRequest()
	return manager.requestManager.SendRequest(ctx, request)
}

// deleteAllEntries clears all Command Table entries from RAM.
func (manager *CommandTableManager) deleteAllEntries(ctx context.Context) error {
	request := protocol_command_tables.NewDeleteAllEntriesRequest()
	return manager.requestManager.SendRequest(ctx, request)
}

// allocateEntry creates a CT entry with the given ID and size.
// size must be even. Returns error if allocation fails.
func (manager *CommandTableManager) allocateEntry(ctx context.Context, id uint16, size uint16) error {
	request, err := protocol_command_tables.NewAllocateEntryRequest(id, size)
	if err != nil {
		return err
	}
	response, err := client_common.SendRequestAndReceive[*protocol_command_tables.AllocateEntryResponse](manager.requestManager, ctx, request)
	if err != nil {
		return err
	}
	// Check RTC status code
	rtcStatus := response.RTCStatus()
	if !protocol_rtc.IsStatusOK(rtcStatus) {
		return &protocol_command_tables.EntryAllocationError{
			EntryID: id,
			Size:    size,
			Status:  rtcStatus,
		}
	}
	return nil
}

// writeEntryData writes up to 4 bytes to a CT entry.
// Must be called repeatedly to write all entry data.
// offset is the byte offset for error reporting.
func (manager *CommandTableManager) writeEntryData(ctx context.Context, id uint16, data []byte, offset int) error {
	request, err := protocol_command_tables.NewWriteEntryDataRequest(id, data)
	if err != nil {
		return err
	}
	response, err := client_common.SendRequestAndReceive[*protocol_command_tables.WriteEntryDataResponse](manager.requestManager, ctx, request)
	if err != nil {
		return err
	}
	rtcStatus := response.RTCStatus()
	// Accept Busy, Incomplete, and OK status codes
	if !protocol_rtc.IsStatusOK(rtcStatus) {
		return &protocol_command_tables.EntryWriteError{
			EntryID: id,
			Offset:  offset,
			Status:  rtcStatus,
		}
	}
	return nil
}

// getEntrySize returns the size of a CT entry.
func (manager *CommandTableManager) getEntrySize(ctx context.Context, id uint16) (uint16, error) {
	request := protocol_command_tables.NewGetEntrySizeRequest(id)
	response, err := client_common.SendRequestAndReceive[*protocol_command_tables.GetEntrySizeResponse](manager.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	rtcStatus := response.RTCStatus()
	if !protocol_rtc.IsStatusOK(rtcStatus) {
		return 0, &protocol_command_tables.EntryReadError{
			EntryID: id,
			Status:  rtcStatus,
		}
	}
	value := response.Value()
	// Extract size from low 16 bits (w4), not high 16 bits (w3)
	// Use uint32 masking to avoid int32 sign extension issues
	v := uint32(value)
	w4 := uint16(v & 0xFFFF)
	// Debug logging when enabled
	if manager.debug.Load() {
		w3 := uint16(v >> 16)
		fmt.Printf("[GetEntrySize] entryID=%d, raw_value=0x%08X (int32=%d), w3=0x%04X (%d), w4=0x%04X (%d), computed_size=%d\n",
			id, v, value, w3, w3, w4, w4, w4)
	}
	return w4, nil
}

// readEntryData reads up to 4 bytes from a CT entry.
// Must be called repeatedly to read all entry data.
// Wire format matches C# library: bytes come from wire as [data[0], data[1], data[2], data[3]].
// With little-endian decoding of value=(w3<<16)|w4:
//   - data[0] = w4 low byte
//   - data[1] = w4 high byte
//   - data[2] = w3 low byte
//   - data[3] = w3 high byte
func (manager *CommandTableManager) readEntryData(ctx context.Context, id uint16) ([]byte, error) {
	request := protocol_command_tables.NewReadEntryDataRequest(id)
	response, err := client_common.SendRequestAndReceive[*protocol_command_tables.ReadEntryDataResponse](manager.requestManager, ctx, request)
	if err != nil {
		return nil, err
	}
	rtcStatus := response.RTCStatus()
	if !protocol_rtc.IsStatusOK(rtcStatus) {
		return nil, &protocol_command_tables.EntryReadError{
			EntryID: id,
			Status:  rtcStatus,
		}
	}
	value := response.Value()
	w3 := uint16(value >> 16)
	w4 := uint16(value & 0xFFFF)
	// Unpack from little-endian wire order: [b0, b1, b2, b3]
	data := []byte{
		byte(w4 & 0xFF), // b0 = w4 low
		byte(w4 >> 8),   // b1 = w4 high
		byte(w3 & 0xFF), // b2 = w3 low
		byte(w3 >> 8),   // b3 = w3 high
	}
	return data, nil
}

// readFullEntry reads an entire CT entry by repeatedly calling readEntryData.
// startTime is optional; if provided (non-zero), it's used for debug logging elapsed time.
func (manager *CommandTableManager) readFullEntry(ctx context.Context, id uint16, startTime ...time.Time) ([]byte, error) {
	var start time.Time
	if len(startTime) > 0 && !startTime[0].IsZero() {
		start = startTime[0]
	}

	debugEnabled := manager.debug.Load()

	// GetEntrySizeRequest and ReadEntryDataRequest don't override OperationTimeout(),
	// so they use the default timeout
	opTimeoutGetSize := protocol_common.DefaultOperationTimeout
	opTimeoutReadData := protocol_common.DefaultOperationTimeout

	if debugEnabled {
		elapsed := time.Since(start)
		deadline, hasDeadline := ctx.Deadline()
		var remaining time.Duration
		if hasDeadline {
			remaining = time.Until(deadline)
		}
		fmt.Printf("[CMDTABLE_READ] entry=%d START elapsed=%v remaining=%v opTimeout_getSize=%v opTimeout_readData=%v\n",
			id, elapsed, remaining, opTimeoutGetSize, opTimeoutReadData)
	}

	size, err := manager.getEntrySize(ctx, id)
	if err != nil {
		if debugEnabled {
			elapsed := time.Since(start)
			deadline, hasDeadline := ctx.Deadline()
			var remaining time.Duration
			if hasDeadline {
				remaining = time.Until(deadline)
			}
			fmt.Printf("[CMDTABLE_READ] entry=%d getEntrySize FAILED elapsed=%v remaining=%v err=%v\n",
				id, elapsed, remaining, err)
		}
		return nil, err
	}

	if debugEnabled {
		elapsed := time.Since(start)
		deadline, hasDeadline := ctx.Deadline()
		var remaining time.Duration
		if hasDeadline {
			remaining = time.Until(deadline)
		}
		fmt.Printf("[CMDTABLE_READ] entry=%d getEntrySize OK size=%d elapsed=%v remaining=%v\n",
			id, size, elapsed, remaining)
	}

	buf := make([]byte, 0, size)
	chunkCount := 0
	for len(buf) < int(size) {
		if debugEnabled {
			elapsed := time.Since(start)
			deadline, hasDeadline := ctx.Deadline()
			var remaining time.Duration
			if hasDeadline {
				remaining = time.Until(deadline)
			}
			fmt.Printf("[CMDTABLE_READ] entry=%d readEntryData chunk=%d bytesRead=%d totalSize=%d elapsed=%v remaining=%v opTimeout=%v\n",
				id, chunkCount, len(buf), size, elapsed, remaining, opTimeoutReadData)
		}

		chunk, err := manager.readEntryData(ctx, id)
		if err != nil {
			if debugEnabled {
				elapsed := time.Since(start)
				deadline, hasDeadline := ctx.Deadline()
				var remaining time.Duration
				if hasDeadline {
					remaining = time.Until(deadline)
				}
				fmt.Printf("[CMDTABLE_READ] entry=%d readEntryData FAILED chunk=%d bytesRead=%d elapsed=%v remaining=%v err=%v\n",
					id, chunkCount, len(buf), elapsed, remaining, err)
			}
			return nil, err
		}
		remaining := int(size) - len(buf)
		if remaining >= 4 {
			buf = append(buf, chunk...)
		} else {
			buf = append(buf, chunk[:remaining]...)
		}
		chunkCount++
	}

	if debugEnabled {
		elapsed := time.Since(start)
		deadline, hasDeadline := ctx.Deadline()
		var remaining time.Duration
		if hasDeadline {
			remaining = time.Until(deadline)
		}
		fmt.Printf("[CMDTABLE_READ] entry=%d COMPLETE bytesRead=%d totalSize=%d chunks=%d elapsed=%v remaining=%v\n",
			id, len(buf), size, chunkCount, elapsed, remaining)
	}

	return buf, nil
}

// GetPresenceMasks retrieves the presence mask values from the drive.
// Each mask is a 32-bit bitmask indicating which entries exist in that range.
func (manager *CommandTableManager) GetPresenceMasks(ctx context.Context) ([8]uint32, error) {
	return manager.getPresenceMasks(ctx)
}

// getPresenceMasks retrieves the presence mask values from the drive.
// Each mask is a 32-bit bitmask indicating which entries exist in that range.
func (manager *CommandTableManager) getPresenceMasks(ctx context.Context) ([8]uint32, error) {
	var masks [8]uint32
	for i := uint8(0); i < protocol_command_tables.PresenceMaskCount; i++ {
		request, err := protocol_command_tables.NewPresenceMaskRequest(i)
		if err != nil {
			return masks, fmt.Errorf("failed to create presence mask request %d: %w", i, err)
		}
		response, err := client_common.SendRequestAndReceive[*protocol_command_tables.PresenceMaskResponse](manager.requestManager, ctx, request)
		if err != nil {
			return masks, fmt.Errorf("failed to read presence mask %d: %w", i, err)
		}
		rtcStatus := response.RTCStatus()
		if !protocol_rtc.IsStatusOK(rtcStatus) {
			return masks, fmt.Errorf("presence mask %d returned error status: 0x%02X", i, rtcStatus)
		}
		// Extract w3 and w4 from response value and combine into uint32.
		// Convert int32 to uint32 first to avoid sign extension.
		value := response.Value()
		valueU32 := uint32(value)
		w3 := uint32(valueU32 >> 16)
		w4 := uint32(valueU32 & 0xFFFF)
		masks[i] = (w3 << 16) | w4
	}
	return masks, nil
}

// extractEntryIDsFromMasks converts the 8 presence masks into a list of entry IDs.
// Each mask covers 32 entry IDs (mask 0 = IDs 0-31, mask 1 = IDs 32-63, etc.).
// Uses inverted logic: bit = 0 indicates entry is present.
func (manager *CommandTableManager) extractEntryIDsFromMasks(masks [8]uint32) []uint16 {
	var entryIDs []uint16
	for maskIdx, mask := range masks {
		baseID := uint16(maskIdx * 32)
		for bit := uint(0); bit < 32; bit++ {
			if mask&(1<<bit) == 0 {
				id := baseID + uint16(bit)
				if id == 0 {
					continue
				}
				entryIDs = append(entryIDs, id)
			}
		}
	}
	return entryIDs
}

// getMinEntrySize determines the minimum CT entry size for this drive.
// Returns EntryMinSize as a safe default.
func (manager *CommandTableManager) getMinEntrySize() uint16 {
	return protocol_command_tables.EntryMinSize
}

// fromWireEntry converts a protocol WireEntry to a domain Entry.
// Converts *int64 parameters to *Param with Literal values.
func fromWireEntry(we protocol_command_tables.WireEntry) Entry {
	var sequencedEntry *int
	if we.SequencedEntry != nil {
		id := int(*we.SequencedEntry)
		sequencedEntry = &id
	}
	// Helper to convert *int64 to *Param
	toParam := func(p *int64) *Param {
		if p == nil {
			return nil
		}
		return &Param{Literal: p}
	}
	return Entry{
		ID:             int(we.ID),
		Name:           we.Name,
		Type:           we.Type,
		SequencedEntry: sequencedEntry,
		Par1:           toParam(we.Par1),
		Par2:           toParam(we.Par2),
		Par3:           toParam(we.Par3),
		Par4:           toParam(we.Par4),
		Par5:           toParam(we.Par5),
		Par6:           toParam(we.Par6),
		Par7:           toParam(we.Par7),
		Par8:           toParam(we.Par8),
	}
}

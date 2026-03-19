package protocol_common

import (
	"testing"
)

func TestStatus_ActPosMM(t *testing.T) {
	tests := []struct {
		name     string
		actPos   int32
		expected float64
	}{
		{
			name:     "Zero position",
			actPos:   0,
			expected: 0.0,
		},
		{
			name:     "10mm position",
			actPos:   100000, // 10mm in 0.1 µm units
			expected: 10.0,
		},
		{
			name:     "Negative position",
			actPos:   -50000, // -5mm
			expected: -5.0,
		},
		{
			name:     "Fractional position",
			actPos:   12345, // 1.2345mm
			expected: 1.2345,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := &Status{ActualPosition: tt.actPos}
			result := status.ActualPositionMM()
			const tolerance = 0.0001
			if diff := result - tt.expected; diff < -tolerance || diff > tolerance {
				t.Errorf("ActualPositionMM() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStatus_DemPosMM(t *testing.T) {
	tests := []struct {
		name     string
		demPos   int32
		expected float64
	}{
		{
			name:     "Zero position",
			demPos:   0,
			expected: 0.0,
		},
		{
			name:     "10mm position",
			demPos:   100000, // 10mm in 0.1 µm units
			expected: 10.0,
		},
		{
			name:     "Negative position",
			demPos:   -50000, // -5mm
			expected: -5.0,
		},
		{
			name:     "Fractional position",
			demPos:   12345, // 1.2345mm
			expected: 1.2345,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := &Status{DemandPosition: tt.demPos}
			result := status.DemandPositionMM()
			if result != tt.expected {
				t.Errorf("DemandPositionMM() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRoundTrip_StatusResponse(t *testing.T) {
	// Create event
	statusOut := &Status{
		StatusWord:     0x1234,
		StateVar:       0x5678,
		ActualPosition: 100000,
		DemandPosition: 200000,
		Current:        500,
		WarnWord:       0xABCD,
		ErrorCode:      0x00,
	}
	response := NewStatusResponse(statusOut)

	// Write to packet
	packet, err := response.WritePacket()
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Read from packet using specific parser
	statusRecv, err := ReadStatusResponse(packet)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	if statusRecv.Status().ActualPosition != 100000 {
		t.Errorf("expected ActualPosition 100000, got %d", statusRecv.Status().ActualPosition)
	}
	if statusRecv.Status().DemandPosition != 200000 {
		t.Errorf("expected DemandPosition 200000, got %d", statusRecv.Status().DemandPosition)
	}
	if statusRecv.Status().StatusWord != 0x1234 {
		t.Errorf("expected StatusWord 0x1234, got 0x%04X", statusRecv.Status().StatusWord)
	}
	if statusRecv.Status().StateVar != 0x5678 {
		t.Errorf("expected StateVar 0x5678, got 0x%04X", statusRecv.Status().StateVar)
	}
	if statusRecv.Status().Current != 500 {
		t.Errorf("expected Current 500, got %d", statusRecv.Status().Current)
	}
	if statusRecv.Status().WarnWord != 0xABCD {
		t.Errorf("expected WarnWord 0xABCD, got 0x%04X", statusRecv.Status().WarnWord)
	}
	if statusRecv.Status().ErrorCode != 0x00 {
		t.Errorf("expected ErrorCode 0x00, got 0x%02X", statusRecv.Status().ErrorCode)
	}
}

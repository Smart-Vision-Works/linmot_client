package protocol_common

import (
	"math"
	"testing"
)

// ============================================================================
// Position Conversion Tests
// ============================================================================

func TestPositionConversions(t *testing.T) {
	tests := []struct {
		name      string
		mm        float64
		wantUnits int32
	}{
		{"Zero", 0.0, 0},
		{"Positive tiny", 0.0001, 1},
		{"Positive small", 0.001, 10},
		{"Positive medium", 10.5, 105000},
		{"Positive large", 100.0, 1000000},
		{"Negative tiny", -0.0001, -1},
		{"Negative small", -0.001, -10},
		{"Negative medium", -10.5, -105000},
		{"Negative large", -100.0, -1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test forward conversion (engineering → protocol)
			gotUnits := ToProtocolPosition(tt.mm)
			if gotUnits != tt.wantUnits {
				t.Errorf("ToProtocolPosition(%v) = %v, want %v", tt.mm, gotUnits, tt.wantUnits)
			}

			// Test reverse conversion (protocol → engineering)
			gotMM := FromProtocolPosition(tt.wantUnits)
			tolerance := 0.00001 // 0.01µm tolerance
			if math.Abs(gotMM-tt.mm) > tolerance {
				t.Errorf("FromProtocolPosition(%v) = %v, want %v (diff: %v)", tt.wantUnits, gotMM, tt.mm, math.Abs(gotMM-tt.mm))
			}
		})
	}
}

func TestPositionRoundTrip(t *testing.T) {
	testValues := []float64{0.0, 1.0, 10.0, 100.0, -1.0, -10.0, -100.0, 0.0001, 123.456}

	for _, mm := range testValues {
		units := ToProtocolPosition(mm)
		recovered := FromProtocolPosition(units)
		tolerance := 0.00001 // 0.01µm tolerance
		if math.Abs(recovered-mm) > tolerance {
			t.Errorf("Round-trip failed for %v mm: got %v mm (diff: %v)", mm, recovered, math.Abs(recovered-mm))
		}
	}
}

// ============================================================================
// Velocity Conversion Tests
// ============================================================================

func TestVelocityConversions(t *testing.T) {
	tests := []struct {
		name      string
		ms        float64
		wantUnits uint32
	}{
		{"Zero", 0.0, 0},
		{"Tiny", 0.000001, 1},
		{"Small", 0.001, 1000},
		{"Medium", 0.5, 500000},
		{"Large", 5.0, 5000000},
		{"Very large", 10.0, 10000000},
		{"Negative", -1.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test forward conversion
			gotUnits := ToProtocolVelocity(tt.ms)
			if gotUnits != tt.wantUnits {
				t.Errorf("ToProtocolVelocity(%v) = %v, want %v", tt.ms, gotUnits, tt.wantUnits)
			}

			// Test reverse conversion
			gotMS := FromProtocolVelocity(tt.wantUnits)
			expectedMS := tt.ms
			if expectedMS < 0 {
				expectedMS = 0
			}
			tolerance := 0.0000001 // 0.1µm/s tolerance
			if math.Abs(gotMS-expectedMS) > tolerance {
				t.Errorf("FromProtocolVelocity(%v) = %v, want %v (diff: %v)", tt.wantUnits, gotMS, expectedMS, math.Abs(gotMS-expectedMS))
			}
		})
	}
}

func TestVelocityRoundTrip(t *testing.T) {
	testValues := []float64{0.0, 0.1, 0.5, 1.0, 5.0, 10.0, 0.000001, 1.234567}

	for _, ms := range testValues {
		units := ToProtocolVelocity(ms)
		recovered := FromProtocolVelocity(units)
		tolerance := 0.0000001
		if math.Abs(recovered-ms) > tolerance {
			t.Errorf("Round-trip failed for %v m/s: got %v m/s (diff: %v)", ms, recovered, math.Abs(recovered-ms))
		}
	}
}

// ============================================================================
// Acceleration Conversion Tests
// ============================================================================

func TestAccelerationConversions(t *testing.T) {
	tests := []struct {
		name      string
		ms2       float64
		wantUnits uint32
	}{
		{"Zero", 0.0, 0},
		{"Tiny", 0.00001, 1},
		{"Small", 0.001, 100},
		{"Medium", 2.5, 250000},
		{"Large", 100.0, 10000000},
		{"Very large", 200.0, 20000000},
		{"Negative", -5.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test forward conversion
			gotUnits := ToProtocolAcceleration(tt.ms2)
			if gotUnits != tt.wantUnits {
				t.Errorf("ToProtocolAcceleration(%v) = %v, want %v", tt.ms2, gotUnits, tt.wantUnits)
			}

			// Test reverse conversion
			gotMS2 := FromProtocolAcceleration(tt.wantUnits)
			expectedMS2 := tt.ms2
			if expectedMS2 < 0 {
				expectedMS2 = 0
			}
			tolerance := 0.000001 // 0.01mm/s² tolerance
			if math.Abs(gotMS2-expectedMS2) > tolerance {
				t.Errorf("FromProtocolAcceleration(%v) = %v, want %v (diff: %v)", tt.wantUnits, gotMS2, expectedMS2, math.Abs(gotMS2-expectedMS2))
			}
		})
	}
}

func TestAccelerationRoundTrip(t *testing.T) {
	testValues := []float64{0.0, 1.0, 2.5, 10.0, 100.0, 200.0, 0.00001, 12.34567}

	for _, ms2 := range testValues {
		units := ToProtocolAcceleration(ms2)
		recovered := FromProtocolAcceleration(units)
		tolerance := 0.000001
		if math.Abs(recovered-ms2) > tolerance {
			t.Errorf("Round-trip failed for %v m/s²: got %v m/s² (diff: %v)", ms2, recovered, math.Abs(recovered-ms2))
		}
	}
}

// ============================================================================
// Jerk Conversion Tests
// ============================================================================

func TestJerkConversions(t *testing.T) {
	tests := []struct {
		name      string
		ms3       float64
		wantUnits uint32
	}{
		{"Zero", 0.0, 0},
		{"Tiny", 0.000001, 1},
		{"Small", 0.001, 1000},
		{"Medium", 10.0, 10000000},
		{"Large", 100.0, 100000000},
		{"Negative", -20.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test forward conversion
			gotUnits := ToProtocolJerk(tt.ms3)
			if gotUnits != tt.wantUnits {
				t.Errorf("ToProtocolJerk(%v) = %v, want %v", tt.ms3, gotUnits, tt.wantUnits)
			}

			// Test reverse conversion
			gotMS3 := FromProtocolJerk(tt.wantUnits)
			expectedMS3 := tt.ms3
			if expectedMS3 < 0 {
				expectedMS3 = 0
			}
			tolerance := 0.0000001
			if math.Abs(gotMS3-expectedMS3) > tolerance {
				t.Errorf("FromProtocolJerk(%v) = %v, want %v (diff: %v)", tt.wantUnits, gotMS3, expectedMS3, math.Abs(gotMS3-expectedMS3))
			}
		})
	}
}

func TestJerkRoundTrip(t *testing.T) {
	testValues := []float64{0.0, 1.0, 10.0, 100.0, 0.000001, 12.345678}

	for _, ms3 := range testValues {
		units := ToProtocolJerk(ms3)
		recovered := FromProtocolJerk(units)
		tolerance := 0.0000001
		if math.Abs(recovered-ms3) > tolerance {
			t.Errorf("Round-trip failed for %v m/s³: got %v m/s³ (diff: %v)", ms3, recovered, math.Abs(recovered-ms3))
		}
	}
}

// ============================================================================
// Signed Variant Tests
// ============================================================================

func TestSignedVelocityConversions(t *testing.T) {
	tests := []struct {
		name   string
		units  int32
		wantMS float64
	}{
		{"Zero", 0, 0.0},
		{"Positive", 500000, 0.5},
		{"Negative", -500000, -0.5},
		{"Large positive", 10000000, 10.0},
		{"Large negative", -10000000, -10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMS := FromProtocolVelocitySigned(tt.units)
			tolerance := 0.0000001
			if math.Abs(gotMS-tt.wantMS) > tolerance {
				t.Errorf("FromProtocolVelocitySigned(%v) = %v, want %v", tt.units, gotMS, tt.wantMS)
			}
		})
	}
}

func TestSignedAccelerationConversions(t *testing.T) {
	tests := []struct {
		name    string
		units   int32
		wantMS2 float64
	}{
		{"Zero", 0, 0.0},
		{"Positive", 250000, 2.5},
		{"Negative", -250000, -2.5},
		{"Large positive", 20000000, 200.0},
		{"Large negative", -20000000, -200.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMS2 := FromProtocolAccelerationSigned(tt.units)
			tolerance := 0.000001
			if math.Abs(gotMS2-tt.wantMS2) > tolerance {
				t.Errorf("FromProtocolAccelerationSigned(%v) = %v, want %v", tt.units, gotMS2, tt.wantMS2)
			}
		})
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestEdgeCases(t *testing.T) {
	t.Run("MaxInt32Position", func(t *testing.T) {
		// Maximum int32 value
		maxUnits := int32(2147483647)
		mm := FromProtocolPosition(maxUnits)
		recovered := ToProtocolPosition(mm)
		if recovered != maxUnits {
			t.Errorf("Max int32 round-trip failed: got %v, want %v", recovered, maxUnits)
		}
	})

	t.Run("MinInt32Position", func(t *testing.T) {
		// Minimum int32 value
		minUnits := int32(-2147483648)
		mm := FromProtocolPosition(minUnits)
		recovered := ToProtocolPosition(mm)
		if recovered != minUnits {
			t.Errorf("Min int32 round-trip failed: got %v, want %v", recovered, minUnits)
		}
	})

	t.Run("VerySmallPosition", func(t *testing.T) {
		// Smallest representable position (0.1µm)
		mm := 0.0001
		units := ToProtocolPosition(mm)
		if units != 1 {
			t.Errorf("Smallest position conversion failed: got %v, want 1", units)
		}
	})

	t.Run("VerySmallVelocity", func(t *testing.T) {
		// Smallest representable velocity (1µm/s)
		ms := 0.000001
		units := ToProtocolVelocity(ms)
		if units != 1 {
			t.Errorf("Smallest velocity conversion failed: got %v, want 1", units)
		}
	})

	t.Run("VerySmallAcceleration", func(t *testing.T) {
		// Smallest representable acceleration (10µm/s²)
		ms2 := 0.00001
		units := ToProtocolAcceleration(ms2)
		if units != 1 {
			t.Errorf("Smallest acceleration conversion failed: got %v, want 1", units)
		}
	})
}

// ============================================================================
// Real-World Value Tests
// ============================================================================

func TestRealWorldValues(t *testing.T) {
	t.Run("TypicalMotionProfile", func(t *testing.T) {
		// Typical values for a pick-and-place application
		position := 50.0 // 50mm stroke
		velocity := 1.0  // 1 m/s max velocity
		accel := 10.0    // 10 m/s² acceleration
		jerk := 100.0    // 100 m/s³ jerk

		posUnits := ToProtocolPosition(position)
		velUnits := ToProtocolVelocity(velocity)
		accelUnits := ToProtocolAcceleration(accel)
		jerkUnits := ToProtocolJerk(jerk)

		// Verify conversions
		if posUnits != 500000 {
			t.Errorf("Position conversion: got %v, want 500000", posUnits)
		}
		if velUnits != 1000000 {
			t.Errorf("Velocity conversion: got %v, want 1000000", velUnits)
		}
		if accelUnits != 1000000 {
			t.Errorf("Acceleration conversion: got %v, want 1000000", accelUnits)
		}
		if jerkUnits != 100000000 {
			t.Errorf("Jerk conversion: got %v, want 100000000", jerkUnits)
		}

		// Verify round-trip
		if FromProtocolPosition(posUnits) != position {
			t.Error("Position round-trip failed")
		}
		if FromProtocolVelocity(velUnits) != velocity {
			t.Error("Velocity round-trip failed")
		}
		if FromProtocolAcceleration(accelUnits) != accel {
			t.Error("Acceleration round-trip failed")
		}
		if FromProtocolJerk(jerkUnits) != jerk {
			t.Error("Jerk round-trip failed")
		}
	})
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkToProtocolPosition(b *testing.B) {
	mm := 10.5
	for i := 0; i < b.N; i++ {
		_ = ToProtocolPosition(mm)
	}
}

func BenchmarkFromProtocolPosition(b *testing.B) {
	units := int32(105000)
	for i := 0; i < b.N; i++ {
		_ = FromProtocolPosition(units)
	}
}

func BenchmarkToProtocolVelocity(b *testing.B) {
	ms := 0.5
	for i := 0; i < b.N; i++ {
		_ = ToProtocolVelocity(ms)
	}
}

func BenchmarkFromProtocolVelocity(b *testing.B) {
	units := uint32(500000)
	for i := 0; i < b.N; i++ {
		_ = FromProtocolVelocity(units)
	}
}

func BenchmarkToProtocolAcceleration(b *testing.B) {
	ms2 := 2.5
	for i := 0; i < b.N; i++ {
		_ = ToProtocolAcceleration(ms2)
	}
}

func BenchmarkFromProtocolAcceleration(b *testing.B) {
	units := uint32(250000)
	for i := 0; i < b.N; i++ {
		_ = FromProtocolAcceleration(units)
	}
}

func BenchmarkToProtocolJerk(b *testing.B) {
	ms3 := 10.0
	for i := 0; i < b.N; i++ {
		_ = ToProtocolJerk(ms3)
	}
}

func BenchmarkFromProtocolJerk(b *testing.B) {
	units := uint32(10000000)
	for i := 0; i < b.N; i++ {
		_ = FromProtocolJerk(units)
	}
}

// Benchmark all conversions for a typical motion command
func BenchmarkFullMotionConversion(b *testing.B) {
	positionMM := 50.0
	velocityMS := 1.0
	accelMS2 := 10.0
	decelMS2 := 10.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToProtocolPosition(positionMM)
		_ = ToProtocolVelocity(velocityMS)
		_ = ToProtocolAcceleration(accelMS2)
		_ = ToProtocolAcceleration(decelMS2)
	}
}

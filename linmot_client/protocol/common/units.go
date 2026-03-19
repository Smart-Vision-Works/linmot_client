package protocol_common

// ============================================================================
// Engineering Units → Protocol Units (Encoding)
// ============================================================================

// ToProtocolPosition converts millimeters to protocol position units (0.1µm).
//
// Protocol units: 1 unit = 0.1µm = 0.0001mm
// Example: 10.5mm → 105000 units
//
// Reference: C library E1100_FACTOR_POS (10000)
func ToProtocolPosition(mm float64) int32 {
	return int32(mm * float64(Factor.Position))
}

// ToProtocolVelocity converts m/s to protocol velocity units (1µm/s).
//
// Protocol units: 1 unit = 1E-6 m/s = 0.001 mm/s
// Example: 0.5 m/s → 500000 units
//
// Reference: C library E1100_FACTOR_SPEED (1000000)
func ToProtocolVelocity(ms float64) uint32 {
	if ms < 0 {
		return 0
	}
	return uint32(ms * float64(Factor.Speed))
}

// ToProtocolVelocitySigned converts m/s to signed protocol units.
func ToProtocolVelocitySigned(ms float64) int32 {
	return int32(ms * float64(Factor.Speed))
}

// ToProtocolAcceleration converts m/s² to protocol acceleration units (10µm/s²).
//
// Protocol units: 1 unit = 1E-5 m/s² = 0.01 mm/s²
// Example: 2.5 m/s² → 250000 units
//
// Reference: C library E1100_FACTOR_ACCEL (100000)
func ToProtocolAcceleration(ms2 float64) uint32 {
	if ms2 < 0 {
		return 0
	}
	return uint32(ms2 * float64(Factor.Acceleration))
}

// ToProtocolAccelerationSigned converts m/s² to signed protocol units.
func ToProtocolAccelerationSigned(ms2 float64) int32 {
	return int32(ms2 * float64(Factor.Acceleration))
}

// ToProtocolJerk converts m/s³ to protocol jerk units (1µm/s³).
//
// Protocol units: 1 unit = 1E-6 m/s³ = 0.001 mm/s³
// Example: 10.0 m/s³ → 10000000 units
//
// Note: Jerk uses same factor as velocity (Factor.Speed = 1000000)
func ToProtocolJerk(ms3 float64) uint32 {
	if ms3 < 0 {
		return 0
	}
	return uint32(ms3 * float64(Factor.Speed))
}

// ToProtocolJerkSigned converts m/s³ to signed protocol units.
func ToProtocolJerkSigned(ms3 float64) int32 {
	return int32(ms3 * float64(Factor.Speed))
}

// ============================================================================
// Protocol Units → Engineering Units (Decoding)
// ============================================================================

// FromProtocolPosition converts protocol position units to millimeters.
//
// Protocol units: 1 unit = 0.1µm = 0.0001mm
// Example: 105000 units → 10.5mm
func FromProtocolPosition(units int32) float64 {
	return float64(units) / float64(Factor.Position)
}

// FromProtocolVelocity converts protocol velocity units to m/s.
//
// Protocol units: 1 unit = 1E-6 m/s = 0.001 mm/s
// Example: 500000 units → 0.5 m/s
func FromProtocolVelocity(units uint32) float64 {
	return float64(units) / float64(Factor.Speed)
}

// FromProtocolAcceleration converts protocol acceleration units to m/s².
//
// Protocol units: 1 unit = 1E-5 m/s² = 0.01 mm/s²
// Example: 250000 units → 2.5 m/s²
func FromProtocolAcceleration(units uint32) float64 {
	return float64(units) / float64(Factor.Acceleration)
}

// FromProtocolJerk converts protocol jerk units to m/s³.
//
// Protocol units: 1 unit = 1E-6 m/s³ = 0.001 mm/s³
// Example: 10000000 units → 10.0 m/s³
func FromProtocolJerk(units uint32) float64 {
	return float64(units) / float64(Factor.Speed)
}

// ============================================================================
// Signed Variants (for streaming/negative values)
// ============================================================================

// FromProtocolVelocitySigned converts signed protocol velocity units to m/s.
// Used in streaming commands where velocity can be negative.
//
// Protocol units: 1 unit = 1E-6 m/s = 0.001 mm/s
// Example: -500000 units → -0.5 m/s
func FromProtocolVelocitySigned(units int32) float64 {
	return float64(units) / float64(Factor.Speed)
}

// FromProtocolAccelerationSigned converts signed protocol acceleration units to m/s².
// Used in streaming commands where acceleration can be negative.
//
// Protocol units: 1 unit = 1E-5 m/s² = 0.01 mm/s²
// Example: -250000 units → -2.5 m/s²
func FromProtocolAccelerationSigned(units int32) float64 {
	return float64(units) / float64(Factor.Acceleration)
}

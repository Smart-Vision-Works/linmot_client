package protocol_control_word

// Package protocol_control_word exposes helpers to build control word patterns that drive the LinMot state machine.
// While higher-level client code (see gsail-go/linmot/client) uses these builders to assemble requests, the package itself
// stays transport-agnostic and only calculates the correct bits that must be sent.
//
// The control word builder examples in the client package demonstrate how the pattern created here maps to
// the client-side EnableDrive/DisableDrive helpers; the protocol package does not depend on the client layer.

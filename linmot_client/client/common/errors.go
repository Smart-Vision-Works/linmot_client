package client_common

import "errors"

// ErrInvalidStatusTelegram indicates a status reply was rejected as invalid telemetry.
var ErrInvalidStatusTelegram = errors.New("invalid status telegram")

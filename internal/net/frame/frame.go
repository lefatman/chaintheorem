package frame

import (
	"encoding/binary"
	"errors"
)

const (
	// HeaderLen is the fixed byte length of the wire header:
	//   u16 msg_type (LE) + u32 payload_len (LE)
	HeaderLen = 2 + 4

	// DefaultMaxPayloadBytes is a conservative default aligned with config/server.json ws.read_limit_bytes (Batch 00).
	// Callers may use a different limit if configured.
	DefaultMaxPayloadBytes = 32768
)

var (
	ErrShortFrame      = errors.New("frame: too short")
	ErrPayloadTooLarge = errors.New("frame: payload too large")
	ErrLengthMismatch  = errors.New("frame: length mismatch")
)

// Encode appends a frame (header + payload) into dst (reusing capacity when possible).
// It returns the resulting slice.
func Encode(dst []byte, msgType uint16, payload []byte, maxPayload uint32) ([]byte, error) {
	if maxPayload == 0 {
		maxPayload = DefaultMaxPayloadBytes
	}
	if uint32(len(payload)) > maxPayload {
		return nil, ErrPayloadTooLarge
	}

	need := HeaderLen + len(payload)
	if cap(dst) < need {
		dst = make([]byte, 0, need)
	}
	dst = dst[:need]

	binary.LittleEndian.PutUint16(dst[0:2], msgType)
	binary.LittleEndian.PutUint32(dst[2:6], uint32(len(payload)))
	copy(dst[HeaderLen:], payload)
	return dst, nil
}

// Decode parses a full wire frame buffer and returns (msgType, payloadSlice).
// The returned payloadSlice aliases the input buffer (zero-copy).
func Decode(b []byte, maxPayload uint32) (uint16, []byte, error) {
	if len(b) < HeaderLen {
		return 0, nil, ErrShortFrame
	}
	msgType := binary.LittleEndian.Uint16(b[0:2])
	payLen := binary.LittleEndian.Uint32(b[2:6])

	if maxPayload == 0 {
		maxPayload = DefaultMaxPayloadBytes
	}
	if payLen > maxPayload {
		return 0, nil, ErrPayloadTooLarge
	}

	expected := HeaderLen + int(payLen)
	if len(b) != expected {
		return 0, nil, ErrLengthMismatch
	}

	return msgType, b[HeaderLen:expected], nil
}

// Package guid provides functions to generate globally unique identifiers (GUIDs).
// This implementation uses UUID version 7 for better sorting and uniqueness.

package guid

import (
	"crypto/rand"
	"time"

	"github.com/google/uuid"
)

// New returns uuid V7 format if err is nil and panics otherwise.
func New() uuid.UUID {
	return must(uuid.NewV7())
}

// NewString returns uuid V7 format as a string.
func NewString() string {
	return New().String()
}

// must returns uuid if err is nil and panics otherwise.
func must(uuid uuid.UUID, err error) uuid.UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}

func NewV7FromTime(t time.Time) (uuid.UUID, error) {
	var value [16]byte

	// 1. Lấy Unix timestamp tính bằng miliseconds (48 bits)
	ms := uint64(t.UnixMilli())

	// Gán 48 bit timestamp vào 6 byte đầu tiên
	value[0] = byte(ms >> 40)
	value[1] = byte(ms >> 32)
	value[2] = byte(ms >> 24)
	value[3] = byte(ms >> 16)
	value[4] = byte(ms >> 8)
	value[5] = byte(ms)

	// 2. Điền các byte còn lại bằng dữ liệu ngẫu nhiên
	if _, err := rand.Read(value[6:16]); err != nil {
		return uuid.Nil, err
	}

	// 3. Thiết lập Version (bit 4-7 của byte 6 là 0111 -> 7)
	value[6] = (value[6] & 0x0f) | 0x70

	// 4. Thiết lập Variant (bit 6-7 của byte 8 là 10 -> RFC 4122)
	value[8] = (value[8] & 0x3f) | 0x80

	return value, nil
}

// NewMinTime returns the smallest UUIDv7 for the given time instant (by ms).
// Assumes t is within the representable range for 48-bit unix milliseconds.
func NewMinTime(t time.Time) uuid.UUID {
	return uuidv7Boundary(t, false)
}

// NewMaxTime returns the largest UUIDv7 for the given time instant (by ms).
// Assumes t is within the representable range for 48-bit unix milliseconds.
func NewMaxTime(t time.Time) uuid.UUID {
	return uuidv7Boundary(t, true)
}

func uuidv7Boundary(t time.Time, max bool) uuid.UUID {
	ms := uint64(t.UTC().UnixMilli())

	var b [16]byte

	// 48-bit unix epoch milliseconds, big-endian
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)

	if !max {
		// min boundary
		b[6] = 0x70 // version=7, rand_a hi=0
		b[7] = 0x00 // rand_a lo
		b[8] = 0x80 // variant=10, rand_b hi=0
		// b[9..15] already zero
		return uuid.UUID(b)
	}

	// max boundary (no loops)
	b[6] = 0x7F // version=7 + rand_a hi=0xF
	b[7] = 0xFF // rand_a lo
	b[8] = 0xBF // variant=10 + rand_b hi=0x3F
	b[9] = 0xFF
	b[10] = 0xFF
	b[11] = 0xFF
	b[12] = 0xFF
	b[13] = 0xFF
	b[14] = 0xFF
	b[15] = 0xFF

	return uuid.UUID(b)
}

func CompareUUIDv7(a, b uuid.UUID) int {
	for i := 0; i < 6; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

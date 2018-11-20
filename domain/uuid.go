package domain

import (
	"github.com/satori/go.uuid"
)

// UUID implement for convert uuid
type UUID [16]byte

// NewUUID create new UUID with V4
func NewUUID() UUID {
	return UUID(uuid.NewV4())
}

// ToBytes to byte array
func (u *UUID) ToBytes() []byte {
	return u[:]
}

// IsZero check uuid is zero
func (u *UUID) IsZero() bool {
	if u == nil {
		return true
	}
	for _, c := range u {
		if c != 0 {
			return false
		}
	}
	return true
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// NewUUIDFromString uuid is made from string
func NewUUIDFromString(s string) (*UUID, error) {
	uuid, err := uuid.FromString(s)
	if err != nil {
		return nil, err
	}

	id := &UUID{}
	for i, c := range uuid {
		id[i] = c
	}

	return id, nil
}

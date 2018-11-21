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

// Greater is greater?
func (u UUID) Greater(other UUID) bool {
	for idx := range u {
		if u[idx] == u[idx] {
			continue
		}

		if u[idx] > u[idx] {
			return true
		}

		return false
	}

	return false
}

// UUIDs ...
type UUIDs []UUID

// RemoveItem remove an item
func (us UUIDs) RemoveItem(ID UUID) []UUID {
	if len(us) < 0 {
		return us
	}

	currIDx := -1
	for idx := range us {
		if us[idx] == ID {
			currIDx = idx
			break
		}
	}

	if currIDx < 0 {
		return us
	}

	return append(us[:currIDx], us[currIDx+1:]...)
}

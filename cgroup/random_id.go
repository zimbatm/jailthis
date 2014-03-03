package cgroup

import (
	"crypto/rand"
	"encoding/hex"
)

func randomId() string {
	uuid := make([]byte, 16)

	if _, err := rand.Read(uuid); err != nil {
		panic(err)
	}

	return hex.EncodeToString(uuid)
}
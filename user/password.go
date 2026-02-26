package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidPassword = errors.New("invalid credentials")
)

type Argon2idParams struct {
	MemoryKiB uint32 // e.g. 64*1024 = 64 MiB
	Time      uint32 // iterations
	Threads   uint8
	SaltLen   uint32
	KeyLen    uint32
}

// Solid “production default” that works on typical servers.
// Tune MemoryKiB upward if you can afford it (better).
var defaultParams = Argon2idParams{
	MemoryKiB: 64 * 1024,
	Time:      3,
	Threads:   2,
	SaltLen:   16,
	KeyLen:    32,
}

// Encoded form:
// $argon2id$v=19$m=65536,t=3,p=2$<salt_b64>$<hash_b64>
func HashPassword(password string) (string, error) {
	salt := make([]byte, defaultParams.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("rand salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.Time,
		defaultParams.MemoryKiB,
		defaultParams.Threads,
		defaultParams.KeyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		defaultParams.MemoryKiB,
		defaultParams.Time,
		defaultParams.Threads,
		b64Salt,
		b64Hash,
	)
	return encoded, nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	// ["", "argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(parts) != 6 || parts[1] != "argon2id" {
		// Treat malformed as invalid rather than “erroring” user-facing
		return false, nil
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != 19 {
		return false, nil
	}

	var mem uint32
	var time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &time, &threads); err != nil {
		return false, nil
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, nil
	}
	wantHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, nil
	}

	gotHash := argon2.IDKey([]byte(password), salt, time, mem, threads, uint32(len(wantHash)))

	return subtle.ConstantTimeCompare(gotHash, wantHash) == 1, nil
}

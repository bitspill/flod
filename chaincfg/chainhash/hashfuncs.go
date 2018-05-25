// Copyright (c) 2015 The Decred developers
// Copyright (c) 2016-2017 The btcsuite developers
// Copyright (c) 2018 The Flo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chainhash

import (
	"crypto/sha256"
	"golang.org/x/crypto/scrypt"
)

// HashB calculates hash(b) and returns the resulting bytes.
func HashB(b []byte) []byte {
	hash := sha256.Sum256(b)
	return hash[:]
}

// HashH calculates hash(b) and returns the resulting bytes as a Hash.
func HashH(b []byte) Hash {
	return Hash(sha256.Sum256(b))
}

// DoubleHashB calculates hash(hash(b)) and returns the resulting bytes.
func DoubleHashB(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:]
}

// DoubleHashH calculates hash(hash(b)) and returns the resulting bytes as a
// Hash.
func DoubleHashH(b []byte) Hash {
	first := sha256.Sum256(b)
	return Hash(sha256.Sum256(first[:]))
}

func Scrypt_1024_1_1_256(b []byte) Hash {
	dk, _ := scrypt.Key(b, b, 1024, 1, 1, 256)
	var b32 [32]byte
	copy(b32[:], dk[:32])
	return Hash(b32)
}

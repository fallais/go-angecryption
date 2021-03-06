package goangecryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

// padding
func padding(data []byte, blocklen int) ([]byte, error) {
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte("\x00"), padlen)

	return append(data, pad...), nil
}

// xorBytes
func xorBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}

	buf := make([]byte, len(a))

	for i := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}

// encryptCBC
func encryptCBC(key, iv, plaintext []byte) ([]byte, error) {
	if len(plaintext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("plaintext is not a multiple of the block size")
	}

	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error while creating the AES block: %s", err)
	}

	ciphertext := make([]byte, len(plaintext))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// decryptECB
func decryptECB(data, key []byte) ([]byte, error) {
	// Create the cipher
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error while creating the AES block: %s", err)
	}

	// Decrypt
	cipher.Decrypt(data, data)
	
	return data, nil
}

// decryptCBC
func decryptCBC(ciphertext, key, iv []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("Text is too short")
	}

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error while creating the AES block: %s", err)
	}

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}

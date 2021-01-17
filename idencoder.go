package main

import (
	"bytes"
	"fmt"
)

// Default values for encoder/decoders
const (
	DefaultAlphabet  = "3fq4rv5z7hsdamn8bpygw96j2cetxuk"
	DefaultBlockSize = 24
	DefaultChecksum  = 29
	MinLength        = 5
)

// IDEncoder contains the various values for an encoder/decoder.
type IDEncoder struct {
	alphabet  []byte
	blockSize uint64
	modulus   uint64
}

// Encode converts an integer to a unique string, using the parameters contianed in the IDEncoder
func (i *IDEncoder) Encode(n, minLength uint64) (string, bool) {

	return string(i.checksum(n)) + i.enbase(i.scramble(n), minLength), false
}

// Decode converts an string to an integer, using the parameters contianed in the IDEncoder
func (i *IDEncoder) Decode(s string) (uint64, bool) {
	b := []byte(s)
	value := i.scramble((i.debase(b[1:])))
	err := false
	if i.checksum(value) != b[0] {
		err = true
	}
	return value, err
}

func (i *IDEncoder) checksum(n uint64) byte {
	return i.alphabet[n%i.modulus]
}

func (i *IDEncoder) scramble(n uint64) uint64 {
	mask := uint64((1 << i.blockSize) - 1)
	result := n & ^mask
	for bit := uint64(0); bit < i.blockSize; bit++ {
		if n&(1<<bit) != 0 {
			result |= 1 << (i.blockSize - bit - 1)
		}
	}
	return result
}

func (i *IDEncoder) enbase(x, minLength uint64) string {
	n := uint64(len(i.alphabet))
	chars := []byte{}
	for x > 0 {
		c := x % n
		x = uint64(x / n)
		chars = append([]byte{i.alphabet[c]}, chars...)
	}
	return leftPad(string(chars), minLength, i.alphabet[0])
}

func (i *IDEncoder) debase(x []byte) uint64 {
	result := uint64(0)
	n := uint64(len(i.alphabet))
	for _, val := range x {
		result *= n
		result += uint64(bytes.IndexByte(i.alphabet, val))
	}
	return result
}

func times(c byte, n uint64) []byte {
	chars := make([]byte, n)
	for i := uint64(0); i < n; i++ {
		chars[i] = c
	}
	return chars
}

func leftPad(str string, length uint64, pad byte) string {
	return fmt.Sprintf("%v%v", string(times(pad, length-uint64(len(str)))), str)
}

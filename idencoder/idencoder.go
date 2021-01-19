/*
Package idencoder is a Golang implementation for encoding (usually sequential) integer IDs.

## Algorithm details

A bit-shuffling approach is used to avoid generating consecutive, predictable
values. However, the algorithm is deterministic and will guarantee that no
collisions will occur.

The encoding alphabet is fully customizable and may contain any number of
characters. By default, digits and lower-case letters are used, with some
characters removed to avoid confusion between characters like o, O and 0. The
default alphabet is shuffled and has a prime number of characters to further
improve the results of the algorithm.

The block size specifies how many bits will be shuffled. The lower `BlockSize`
bits are reversed. Any bits higher than `BlockSize` will remain as is.
`BlockSize` of 0 will leave all bits unaffected and the algorithm will simply
be converting your integer to a different base.

## Common usage

### URL shortening & obfuscation

The intended use is that incrementing, consecutive integers will be used as
keys to generate the encoded IDs. For example, to create a new short URL (à la
bit.ly), the unique integer ID assigned by a database could be used to generate
the last portion of the URL by using this module. Or a simple counter may be
used. As long as the same integer is not used twice, the same encoded value
will not be generated twice.

The module supports both encoding and decoding of values. The `minLength`
parameter allows you to pad the encoded value if you want it to be a specific
length.

## Provenance

Original Author (Python): [Michael Fogleman](http://code.activestate.com/recipes/576918/)
License: [MIT](https://opensource.org/licenses/MIT)

Modified Python version from which this version is ported
URL: https://github.com/brnt/idencoder

Repo: https://github.com/brnt/idencoder-go


*/
package idencoder

import (
	"bytes"
	"fmt"
)

// Alphabet is a set of characters to be used in an encoded value
type Alphabet []byte

// BlockSize determines how many bits will be shuffled
type BlockSize uint64

// Checksum validates scramble/unscramble sub-operations
type Checksum uint64

// Default values for encoder/decoders
const (
	// DefaultAlphabet SHOULD NOT be used in production!!! This value
	// can/will change at any time, which will break links to any URLs
	// with encoded IDs in them. Instead, create your own alphabets,
	// perhaps using the idencoder.RandomAlphabet() function.
	DefaultAlphabet  = "3fq4rv5z7hsdamn8bpygw96j2cetxuk"
	DefaultBlockSize = 24
	DefaultChecksum  = 29
	MinLength        = 5
)

// IDEncoder contains the various values for an encoder/decoder.
type IDEncoder struct {
	Alphabet  Alphabet
	BlockSize BlockSize
	Checksum  Checksum
}

// Encode converts an integer to a unique string, using the parameters contianed in the IDEncoder
func (i *IDEncoder) Encode(n, minLength uint64) (encoded string, ok bool) {

	return string(i.checksum(n)) + i.enbase(i.scramble(n), minLength), true
}

// Decode converts an string to an integer, using the parameters contianed in the IDEncoder
func (i *IDEncoder) Decode(s string) (decoded uint64, ok bool) {
	b := []byte(s)
	value := i.scramble((i.debase(b[1:])))
	err := true
	if i.checksum(value) != b[0] {
		err = false
	}
	return value, err
}

func (i *IDEncoder) checksum(n uint64) byte {
	return i.Alphabet[n%uint64(i.Checksum)]
}

func (i *IDEncoder) scramble(n uint64) uint64 {
	mask := uint64((1 << i.BlockSize) - 1)
	result := n & ^mask
	for bit := uint64(0); bit < uint64(i.BlockSize); bit++ {
		if n&(1<<bit) != 0 {
			result |= 1 << (uint64(i.BlockSize) - bit - 1)
		}
	}
	return result
}

func (i *IDEncoder) enbase(x, minLength uint64) string {
	n := uint64(len(i.Alphabet))
	chars := []byte{}
	for x > 0 {
		c := x % n
		x = uint64(x / n)
		chars = append([]byte{i.Alphabet[c]}, chars...)
	}
	return leftPad(string(chars), minLength, i.Alphabet[0])
}

func (i *IDEncoder) debase(x []byte) uint64 {
	result := uint64(0)
	n := uint64(len(i.Alphabet))
	for _, val := range x {
		result *= n
		result += uint64(bytes.IndexByte(i.Alphabet, val))
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

Integer ID Encoder
===================

[//]: # "insert package idencoder/idencoder.go"

Package idencoder is a Golang implementation for encoding (usually sequential) unsigned integer IDs.

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

## Common application

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

### Alphabet selection
**Do not use the default alphabet in production!** There are two reasons:

1. Your encoded values will be easy for anyone to decode, eliminating the obfuscation value.
2. The value of `idencoder.DefaultAlphabet` can/will change in the code base without warning.

One common approach is to assign a unique alphabet to each entity type in an application, and for no entities to use `idencoder.DefaultAlphabet`. It is wise to store your alphabet in the application configuration, rather than directly in your code.

Random alphabets are perfectly appropriate to use in production. They're also trivial to generate by passing the `-r` or `--random` flags to the command line application `main.go`:

```sh
$ go run main.go -r
```

Or, if you have built the executable application:

```sh
$ ./idencoder -r
```

## Provenance

Original Author (Python): [Michael Fogleman](http://code.activestate.com/recipes/576918/);
License: [MIT](https://opensource.org/licenses/MIT)

Updated module (Python): [Brent Thomson](https://github.com/brnt/idencoder) (source from which this version is ported).

Home for this project: https://github.com/brnt/idencoder-go

[//]: # "insert-end"

## Usage

A very simple example:

```go
package main

import (
    "fmt"
    "github.com/brnt/idencoder-go/idencoder"
)

alpha := "tvy7fuk4g59d6b3mhc8jqwrzexspan2"

func main() {
	ie = idencoder.IDEncoder{
		Alphabet:  []byte(alpha),
		BlockSize: idencoder.DefaultBlockSize,
		Checksum:  idencoder.DefaultChecksum,
	}

	for i := uint64(1); i <= 10; i++ {
		encoded, ok := ie.Encode(i, idencoder.MinLength)
		if !ok {
			fmt.Println("Something is broken")
		}

		decoded, ok := ie.Decode(encoded)
		if !ok || decoded != i {
			fmt.Println("Something is broken")
		}
		fmt.Println(i, encoded, decoded)
    }
}
```
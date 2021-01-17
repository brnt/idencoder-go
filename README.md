# idencoder-go

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
		decoded, ok := ie.Decode(encoded)
		if !ok || decoded != i {
			fmt.Println("Something is broken")
		}
		fmt.Println(i, encoded, decoded)
    }
}
```
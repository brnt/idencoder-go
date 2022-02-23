package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/akamensky/argparse"

	"github.com/brnt/idencoder-go/idencoder"
)

// randomAlphabet generates a random alphabet, containing the same characters as idencoder.DefaultAlphabet
func randomAlphabet() string {
	rand.Seed(time.Now().Unix())
	aRune := []rune(idencoder.DefaultAlphabet)
	rand.Shuffle(len(aRune), func(i, j int) {
		aRune[i], aRune[j] = aRune[j], aRune[i]
	})
	return string(aRune)
}

func main() {
	parser := argparse.NewParser("idencoder", "Description of my awesome program. It can be as long as I wish it to be")
	var alphabet *string = parser.String("a", "alphabet",
		&argparse.Options{
			Required: false,
			Default:  idencoder.DefaultAlphabet,
			Help:     "use ALPHA as the alphabet",
		})
	var quiet *bool = parser.Flag("q", "quiet",
		&argparse.Options{
			Required: false,
			Help:     "suppress formatting and instructional output",
		})
	var length *int = parser.Int("l", "length",
		&argparse.Options{
			Required: false,
			Default:  idencoder.MinLength,
			Help:     "set min encoded output length to NUM",
		})
	var encode *int = parser.Int("e", "encode",
		&argparse.Options{
			Required: false,
			Help:     "encode NUM",
		})
	var decode *string = parser.String("d", "decode",
		&argparse.Options{
			Required: false,
			Help:     "decode NUM",
		})
	var benchmark *int = parser.Int("b", "benchmark",
		&argparse.Options{
			Required: false,
			Help:     "run a series of NUM encode/decode cycles",
		})
	var random *bool = parser.Flag("r", "random",
		&argparse.Options{
			Required: false,
			Help:     "print a random alphabet",
		})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	ie := idencoder.IdEncoder{
		Alphabet:  []byte(*alphabet),
		BlockSize: idencoder.DefaultBlockSize,
		Checksum:  idencoder.DefaultChecksum,
	}
	switch true {
	case *encode > 0:
		encoded, err := ie.Encode(uint64(*encode), uint64(*length))
		if err != nil {
			fmt.Println("**ERROR** during encode")
		}
		fmt.Println(encoded)
	case len(*decode) >= *length:
		decoded, err := ie.Decode(*decode)
		if err != nil {
			fmt.Println("**ERROR** during decode")
		}
		fmt.Println(decoded)
	case *benchmark > 0:
		start := time.Now().UnixNano()
		for i := uint64(0); i < uint64(*benchmark); i++ {
			encoded, err := ie.Encode(i, idencoder.MinLength)
			if err != nil {
				fmt.Println("Something is weird (encode):", i, encoded)
				break
			}
			decoded, err := ie.Decode(encoded)
			if err != nil || i != decoded {
				fmt.Println("Something is weird (decode):", i, encoded, decoded, err)
				break
			}
		}
		end := time.Now().UnixNano()
		p := message.NewPrinter(language.English)
		p.Printf("BENCHMARK: Ran %d iterations in %0.3f seconds\n", *benchmark, float64(end-start)/1000000000)
	case *random:
		alpha := randomAlphabet()
		if *quiet {
			fmt.Println(alpha)
		} else {
			fmt.Println("Random alphabet:", alpha)
		}
	default:
		fmt.Print(parser.Usage("Must select one of encode, decode, random, or benchmark"))
	}

}

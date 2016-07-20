package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Hunter-Dolan/midrange/matcher"
	"github.com/Hunter-Dolan/midrange/transaction"
)

type segmentTest struct {
	size,
	noverlap int
	out [][]float64
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Creating Wav")

	transaction := transaction.NewTransaction()

	transaction.BaseFrequency = 10000
	transaction.FrameDuration = 500
	transaction.Carriers = 64
	transaction.Kilobitrate = 96
	transaction.Bandwidth = 1000
	transaction.NoiseLevel = 0
	transaction.KeyStates = 2

	s := "Hello this is a test, that is a bit longer than the other tests"

	transaction.SetData(s)

	//transaction.Build()
	process(transaction, s)
}

func process(transaction *transaction.Transaction, s string) {
	wave := transaction.Wave()

	fmt.Println("Done. Processing Wave...")

	options := matcher.Options{}
	options.GenerationOptions = transaction.FrameGenerationOptions()
	options.NFFTPower = 16

	bank := matcher.NewMatcher(&options)

	result := bank.Decode(wave)

	fmt.Println(result)

	percentMatch(s, result)
}

func percentMatch(first string, second string) {

	total := float64(len(first))
	matched := float64(0)

	for i := range first {
		firstChar := first[i]
		secondChar := second[i]

		if firstChar == secondChar {
			matched++
			fmt.Print(" ")
		} else {
			fmt.Print("^")
		}
	}

	fmt.Println("     |", matched/total*100, "%|")
}

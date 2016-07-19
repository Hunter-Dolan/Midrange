package main

import (
	"fmt"
	"math/rand"
	"strings"
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

	transaction.BaseFrequency = 1000
	transaction.FrameDuration = 1000
	transaction.Carriers = 32
	transaction.Kilobitrate = 96
	transaction.Bandwidth = 1000
	transaction.NoiseLevel = 80

	s := "Everyone is talking about the car, I'm much more impressed by how sturdy that bike is."

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

	fmt.Println(strings.Contains(result, s))
}

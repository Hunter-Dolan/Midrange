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

	transaction.BaseFrequency = 2000
	transaction.FrameDuration = 500
	transaction.Kilobitrate = 96 * 2
	transaction.Bandwidth = 1000
	transaction.NoiseLevel = 0
	transaction.OMFSKConstant = 2

	s := `Overview: Argentina, rich in natural resources, benefits also from a highly literate population, an export-oriented agricultural sector, and a diversified industrial base. Nevertheless, following decades of mismanagement and statist policies, the economy in the late 1980s was plagued with huge external debts and recurring bouts of hyperinflation. Elected in 1989, in the depths of recession, President MENEM has implemented a comprehensive economic restructuring program that shows signs of putting Argentina on a path of stable, sustainable growth. Argentina's currency has traded at par with the US dollar since April 1991, and inflation has fallen to its lowest level in 20 years. Argentines have responded to the relative price stability by repatriating flight capital and investing in domestic industry. The economy registered an impressive 6% advance in 1994, fueled largely by inflows of foreign capital and strong domestic consumption spending.`

	transaction.SetData(s)

	transaction.Build()
	process(transaction, s)
}

func process(transaction *transaction.Transaction, s string) {
	wave := transaction.Wave()

	fmt.Println("Done. Processing Wave...")

	options := matcher.Options{}
	options.GenerationOptions = transaction.FrameGenerationOptions()
	options.NFFTPower = 17

	bank := matcher.NewMatcher(&options)

	result := bank.Decode(wave)

	fmt.Println(result)

	percentMatch(s, result)
}

func stringToBin(s string) string {
	binaryString := ""

	for _, c := range s {
		binary := fmt.Sprintf("%b", c)
		binaryLen := len(binary)

		padAmount := 8 - binaryLen

		for i := 0; i < padAmount; i++ {
			binary = "0" + binary
		}

		binaryString += binary
	}
	return binaryString
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

	firstBin := stringToBin(first)
	secondBin := stringToBin(second)

	totalBin := float64(len(firstBin))

	matchedBin := float64(0)

	for i := range firstBin {
		firstChar := firstBin[i]
		secondChar := secondBin[i]

		if firstChar == secondChar {
			matchedBin++
		}
	}

	fmt.Println("\n\n|", matched/total*100, "byte %", matchedBin/totalBin*100, "bin %|")
}

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

	transaction := transaction.NewTransaction()

	transaction.SetData("Lorem ipsum set dolar amit consectetur adipiscing elit. Nulla sollicitudin dui eu est dictum, vel faucibus mi blandit. Nunc malesuada vel ipsum vitae varius. Integer dignissim, nulla at volutpat tincidunt, odio ligula mollis ipsum, eget convallis justo ante in risus. Aenean nec nunc sit amet turpis condimentum scelerisque quis ut augue. Proin porta congue ex at ornare. Fusce fermentum nisi libero, ac laoreet dui finibus a. Morbi tincidunt magna ac maximus condimentum. Phasellus porttitor nulla sed dapibus cursus. Suspendisse potenti. Ut placerat dui eu risus suscipit condimentum. Integer vitae fringilla risus, et rhoncus mauris. In eget mi libero. In hac habitasse platea dictumst. Nulla pellentesque at mauris eu consectetur. Interdum et malesuada fames ac ante ipsum primis in faucibus.")

	//transaction.Build()
	//*
	options := matcher.Options{}
	options.GenerationOptions = transaction.FrameGenerationOptions()
	options.NFFTPower = 16

	bank := matcher.NewMatcher(&options)

	result := bank.Decode(transaction.Wave())

	fmt.Println(result)
	//*/
	/*
		transaction := transaction.NewTransaction()



		wave := transaction.Wave()

		//fmt.Println(wave)

		fs := float64(2)
		opts := &spectral.PwelchOptions{}

		p, freq := spectral.Pwelch(wave, fs, opts)

		for i, _p := range p {
			_f := freq[i]

			fmt.Println(_f, ",", _p)
		}
	*/
}

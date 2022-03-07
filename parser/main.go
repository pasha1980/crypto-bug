package parser

import (
	"crypto-bug/parser/src/algorithm"
)

var Algorithms = []algorithm.Algorithm{
	algorithm.AbnormallyPriceAlgorithm{},
	algorithm.ExchangeDifferenceAlgorithm{},
	algorithm.StableDifferenceAlgorithm{},
}

func Init() {
	for _, algo := range Algorithms {
		algo.Analyze()
	}
}

package parser

import (
	"crypto-bug/parser/src/algorithm"
)

var Algorithms = []algorithm.Algorithm{
	algorithm.AbnormallyPriceAlgorithm{},
	algorithm.ExchangeDifferenceAlgorithm{},
}

func Init() {
	for _, algo := range Algorithms {
		algo.Analyze()
	}
}

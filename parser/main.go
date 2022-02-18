package parser

import "crypto-bug/parser/config"

func Init() {
	for _, algo := range config.Algorithms {
		algo.Analyze()
	}
}

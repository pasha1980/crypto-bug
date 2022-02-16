package exchages

import "fmt"

type Binance struct {
}

func (binance Binance) Save(track string, base string) { // todo
	fmt.Println(track, base)
}

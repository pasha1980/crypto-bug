package test

import (
	"crypto-bug/model"
	"crypto-bug/service/utils"
	"fmt"
)

type UtilsTest struct {
}

func (utilsTest UtilsTest) Do() bool {

	// Array test
	var array = []string{
		"abs",
		"qwe",
		"yui",
	}
	exist, index := utils.InArray("abs", array)
	if !exist || index != 0 {
		fmt.Println("Not found value in array")
		return false
	}

	exist, index = utils.InArray("ghgh", array)
	if exist || index != nil {
		fmt.Println("Found non-existing value in array")
		return false
	}

	// Map test
	var mapVal = map[string]interface{}{
		"int":    453,
		"string": "string",
		"float":  324.2,
		"struct": model.Quote{
			Exchange: "Test exchange",
		},
	}
	exist, index = utils.InArray("string", mapVal)
	if !exist {
		fmt.Println("Not found string in map", exist, index)
		return false
	}

	exist, index = utils.InArray(453, mapVal)
	if !exist {
		fmt.Println("Not found int in map", exist, index)
		return false
	}

	exist, index = utils.InArray(324.2, mapVal)
	if !exist {
		fmt.Println("Not found float in map", exist, index)
		return false
	}

	exist, index = utils.InArray(model.Quote{Exchange: "Test exchange"}, mapVal)
	if !exist {
		fmt.Println("Not found struct in map", exist, index)
		return false
	}

	exist, index = utils.InArray(2542525, mapVal)
	if exist {
		fmt.Println("Found non-existing value in array")
		return false
	}

	return true
}

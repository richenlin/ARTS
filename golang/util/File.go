package util

import (
	"fmt"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadFile(name string) {
	b, err := ioutil.ReadFile(name)
	check(err)
	fmt.Println(b)
	str := string(b)
	fmt.Println(str)
}

func WriteFile(name string, data []byte) {
	err := ioutil.WriteFile(name, data, 0755)
	check(err)
}

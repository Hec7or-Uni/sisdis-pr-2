package main

import (
	"io/ioutil"
	"sisdis-pr-2/cmd"
)

func LeerFichero() string {
	data, err := ioutil.ReadFile("file.txt")
	cmd.CheckError(err)
	return string(data)
}

func main() {
	LeerFichero()
}

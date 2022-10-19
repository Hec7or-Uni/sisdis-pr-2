package main

import (
	"io/ioutil"
	"os"
	"sisdis-pr-2/cmd"
	"sisdis-pr-2/ra"
	"strconv"
)


func LeerFichero() string {
	data, err := ioutil.ReadFile("../../file.txt")
	cmd.CheckError(err)
	return string(data)
}

func main() {
	const ACTOR = "lector"
	const ITERACIONES = 20

	args := os.Args[1:]
	PID, _ := strconv.Atoi(args[0])
	
	ra := ra.New(PID, args[1], cmd.LECTOR)

	for i := 1; i < ITERACIONES; i++ {
		ra.PreProtocol()
		LeerFichero()
		ra.PostProtocol()
	}
	ra.Stop()
}

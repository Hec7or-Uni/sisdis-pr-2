package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sisdis-pr-2/cmd"
	"sisdis-pr-2/ra"
	"strconv"
	"time"
)

func LeerFichero() string {
	data, err := ioutil.ReadFile("file.txt")
	cmd.CheckError(err)
	return string(data)
}

func main() {
	const ACTOR = "lector"
	const ITERACIONES = 20

	args := os.Args[1:]
	PID, _ := strconv.Atoi(args[0])
	
	ra := ra.New(PID, args[1], cmd.LECTOR)
	time.Sleep(1 * time.Second)
	for i := 0; i < ITERACIONES; i++ {
		ra.PreProtocol()
		LeerFichero()
		fmt.Printf("%d esta leyendo...\n", PID)
		ra.PostProtocol()
		time.Sleep(3 * time.Millisecond)
	}
	ra.Stop()
}

package main

import (
	"fmt"
	"os"
	"sisdis-pr-2/cmd"
	"sisdis-pr-2/ra"
	"strconv"
	"time"
)

func EscribirFichero(fragmento string) {
	file, err := os.OpenFile("file.txt", os.O_APPEND|os.O_RDWR, 0600)
	cmd.CheckError(err)
	defer file.Close()
		
	len, err := file.WriteString(fragmento)
	cmd.CheckError(err)

	fmt.Printf("\nLength: %d bytes", len)
}

func main() {
	const ACTOR = "escritor"
	const ITERACIONES = 20

	args := os.Args[1:]
	PID, _ := strconv.Atoi(args[0])
	
	ra := ra.New(PID, args[1], cmd.ESCRITOR)
	time.Sleep(3 * time.Second)
	for i := 0; i < ITERACIONES; i++ {
		ra.PreProtocol()
		EscribirFichero("Hola mundo\n")
		ra.PostProtocol()
	}
	ra.Stop()
}

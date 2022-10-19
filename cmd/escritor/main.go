package main

import (
	"fmt"
	"os"
	"sisdis-pr-2/cmd"
	"sisdis-pr-2/ra"
	"strconv"
)

func EscribirFichero(path, fragmento string) {
	file, err := os.OpenFile(path, os.O_APPEND, 0600)
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

	for i := 1; i < ITERACIONES; i++ {
		ra.PreProtocol()
		EscribirFichero("file.txt", "Hola mundo")
		ra.PostProtocol()
	}
	ra.Stop()
}

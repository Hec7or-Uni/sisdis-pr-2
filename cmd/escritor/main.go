package main

import (
	"fmt"
	"os"
	"sisdis-pr-2/cmd"
)

func EscribirFichero(fragmento string) {
	file, err := os.OpenFile("file.txt", os.O_APPEND, 0600)
	cmd.CheckError(err)
	defer file.Close()
		
	len, err := file.WriteString(fragmento)
	cmd.CheckError(err)

	fmt.Printf("\nLength: %d bytes", len)
}


func main() {
	EscribirFichero("Hola mundo")
}

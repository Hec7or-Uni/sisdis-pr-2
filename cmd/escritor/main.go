package main

import (
	"fmt"
	"os"
	"sisdis-pr-2/cmd"
	"sisdis-pr-2/ra"
	"strconv"
	"time"
)

func main() {
	const ACTOR = "escritor"
	const ITERACIONES = 20

	args := os.Args[1:]
	PID, _ := strconv.Atoi(args[0])
	ra := ra.New(PID, args[1], cmd.ESCRITOR)
	time.Sleep(1 * time.Second)
	for i := 0; i < ITERACIONES; i++ {
		ra.PreProtocol()
		fragmento := fmt.Sprintf("Escritor %d escribiendo en el fichero.\n", PID)
		cmd.EscribirFichero(args[2], fragmento)
		fmt.Println(fragmento)
		ra.PostProtocol(fragmento)
		time.Sleep(20 * time.Millisecond)
	}
	ra.Stop()
}

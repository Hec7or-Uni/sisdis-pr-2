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
	const ACTOR = "lector"
	const ITERACIONES = 20

	args := os.Args[1:]
	PID, _ := strconv.Atoi(args[0])
	
	ra := ra.New(PID, args[1], cmd.LECTOR)
	time.Sleep(1 * time.Second)
	for i := 0; i < ITERACIONES; i++ {
		ra.PreProtocol()
		cmd.LeerFichero(args[2])
		fmt.Printf("Lector %d leyendo el fichero.\n", PID)
		ra.PostProtocol("")
		time.Sleep(20 * time.Millisecond)
	}
	ra.Stop()
}

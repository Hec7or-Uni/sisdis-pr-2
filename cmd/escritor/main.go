package main

import (
	"fmt"
	"os"
	"sisdis-pr-2/cmd"
	"sisdis-pr-2/ra"
	"strconv"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
)

func EscribirFichero(fragmento string) {
	file, err := os.OpenFile("file.txt", os.O_APPEND|os.O_RDWR, 0600)
	cmd.CheckError(err)
	defer file.Close()
		
	file.WriteString(fragmento)
}

func main() {
	const ACTOR = "escritor"
	const ITERACIONES = 20

	args := os.Args[1:]
	PID, _ := strconv.Atoi(args[0])
	logger := govec.InitGoVector(fmt.Sprintf("P%d", PID), fmt.Sprintf("logs/r%d", PID), govec.GetDefaultConfig())
	ra := ra.New(PID, args[1], cmd.ESCRITOR, logger)
	time.Sleep(1 * time.Second)
	for i := 0; i < ITERACIONES; i++ {
		ra.PreProtocol(logger)
		EscribirFichero("Hola mundo\n")
		fmt.Printf("%d esta escribiendo...\n", PID)
		ra.PostProtocol(logger)
		time.Sleep(3 * time.Millisecond)
	}
	ra.Stop()
}

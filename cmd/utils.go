package cmd

import (
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}


type ACTOR string

const (
	LECTOR	ACTOR	=	"lector"
	ESCRITOR			= "escritor"
)

// RD && RD -> FALSE
// RD && WR -> TRUE
// WR && RD -> TRUE
// WR && WR -> TRUE
func Exclude(A1 ACTOR, A2 ACTOR) bool {
	return A1 == ESCRITOR || A2 == ESCRITOR
}

func Max(a, b int) int {
	if a > b { return a}
	return b
}
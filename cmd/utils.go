package cmd

import (
	"fmt"
	"io/ioutil"
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
	LECTOR   ACTOR = "lector"
	ESCRITOR       = "escritor"
)

// RD && RD -> FALSE
// RD && WR -> TRUE
// WR && RD -> TRUE
// WR && WR -> TRUE
func Exclude(A1 ACTOR, A2 ACTOR) bool {
	return A1 == ESCRITOR || A2 == ESCRITOR
}

func MaxArray(a []int, b []int) []int {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			a[i] = b[i]
		}
	}
	return a
}

func Max(a []int) int {
	max := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}
	
	return max
}

func EscribirFichero(path, fragmento string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0600)
	CheckError(err)
	defer file.Close()

	file.WriteString(fragmento)
}

func LeerFichero(path string) string {
	data, err := ioutil.ReadFile(path)
	CheckError(err)
	return string(data)
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"go.mozilla.org/mozlog"
)

func main() {
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	entry, err := mozlog.FromJSON(string(line))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n")
	entry.Evaluate()
	entry.ToString()
}

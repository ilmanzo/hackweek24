package main

import (
	"bufio"
	"log"
	"os/exec"

	"github.com/fatih/color"
)

func main() {
	cmd := exec.Command("./myscript.sh")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		color.Red(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

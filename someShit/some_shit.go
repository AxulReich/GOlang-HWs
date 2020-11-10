package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var x int

	for scanner.Scan() {
		txt := scanner.Text()
		n, err := strconv.Atoi(txt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing number: '%v', err: %v\n", txt, err)
			continue
		}
		x ^= n

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v:", err)
	}

	res := int(0 ^ x)

	fmt.Printf("%v\n", res)

}






package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func inputString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func inputInt(prompt string) int64 {
	var input int64
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

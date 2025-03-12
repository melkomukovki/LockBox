package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func inputString(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(inputStyle.Render(prompt))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func inputInt(prompt string) (int64, error) {
	var input int64
	fmt.Print(inputStyle.Render(prompt))
	_, err := fmt.Scanln(&input)
	if err != nil {
		return 0, err
	}
	return input, nil
}

func displayError(err error) {
	fmt.Println(errorStyle.Render(err.Error()))
}

func displayHeader(s string) {
	fmt.Println(headerStyle.Render(s))
}

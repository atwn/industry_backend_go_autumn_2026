package main

import (
	"fmt"
	"strings"
)

func greet(name string) string {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) > 0 {
		return fmt.Sprintf("Hello, %v!", trimmed)
	}

	return "Hello, World!"
}

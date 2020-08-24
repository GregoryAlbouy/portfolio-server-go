package clog

import (
	"runtime"
)

type color struct {
	reset string
	red   string
	green string
	blue  string
}

var c = color{}

func init() {
	if runtime.GOOS != "windows" {
		c.reset = "\033[0m"
		c.red = "\033[31m"
		c.green = "\033[32m"
		c.blue = "\033[34m"
	}
}

// Red outputs the input string in red color on compatible systems
func Red(s string) string {
	return c.red + s + c.reset
}

// Green outputs the input string in green color on compatible systems
func Green(s string) string {
	return c.green + s + c.reset
}

// Blue outputs the input string in blue color on compatible systems
func Blue(s string) string {
	return c.blue + s + c.reset
}

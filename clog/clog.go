// Package clog stands for custom-log
package clog

import (
	"fmt"
	"log"
)

// Printlb (print-lalbel) prints the input preceded by a label.
// Can be handy for debugging.
func Printlb(v interface{}, label string) {
	log.Printf("%s\n%v\n\n", label, v)
}

// Errorlb prints any value preceded by a label with a red "ERROR" tag.
func Errorlb(v interface{}, label string) {
	Printlb(v, fmt.Sprintf("%s: %s", Red("ERROR"), label))
}

// Fatallb (fatal-label) prints the input preceded by a label and exits
func Fatallb(v interface{}, label string) {
	log.Fatalf("%s:\n%v\n\n", label, v)
}

// Package clog stands for custom-log
package clog

import (
	"log"
)

// Printlb (print-lalbel) prints the input preceded by a label.
// Can be handy for debugging.
func Printlb(v interface{}, label string) {
	log.Printf("%s:\n%v\n\n", label, v)
}

// Fatallb (fatal-label) prints the input preceded by a label and exits
func Fatallb(v interface{}, label string) {
	log.Fatalf("%s:\n%v\n\n", label, v)
}

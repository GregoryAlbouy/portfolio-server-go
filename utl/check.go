package utl

import "log"

// Fatal prints an error and exits the process
func Fatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Check prints an error and panics
func Check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// Print prints an error and continue
func Print(err error) {
	if err != nil {
		log.Println(err)
	}
}

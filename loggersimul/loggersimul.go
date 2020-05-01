package loggersimul

import (
	"bytes"
	"log"
	"os"
)

// SimulateLogger simulates logger - allows catching errors
// https://stackoverflow.com/questions/44119951/how-to-check-a-log-output-in-go-test
// need to check the actual log message as well
func SimulateLogger() {
	var buf bytes.Buffer
	log.SetOutput(&buf)
}

// UnsimulateLogger restores logger back
func UnsimulateLogger() {
	log.SetOutput(os.Stderr)
}

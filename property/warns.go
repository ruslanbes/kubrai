package property

import (
	"fmt"
	"log"
)

func warnFileReadError(err error) {
	log.Println(fmt.Errorf("WARN : File read error: %s", err))
}

func warnConvertToIntError(err error) {
	log.Println(fmt.Errorf("WARN : Convert to int error: %s", err))
}

func warnNonBool(b string) {
	log.Println(fmt.Errorf("WARN : Non-bool: %s", b))
}

package property

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// PropertyDir is the location where to look for property files
var PropertyDir = "./properties"

// AsString gets property as string.
// On error: warn and return empty string
func AsString(name string) string {
	bs, err := ioutil.ReadFile(PropertyDir + "/" + name)
	if err != nil {
		warnFileReadError(err)
		return ""
	}
	str := strings.Trim(string(bs), " \n")

	return str
}

// AsByte gets property as byte.
// On error: warn and return 0
func AsByte(name string) byte {
	return byte(AsInt(name))
}

// AsInt gets property as int.
// On error: warn and return 0
func AsInt(name string) int {
	num, err := strconv.ParseInt(AsString(name), 10, 64)
	if err != nil {
		warnConvertToIntError(err)
		return 0
	}
	return int(num)
}

// AsBool gets property as bool.
// On error: warn and return false
func AsBool(name string) bool {
	b := strings.Trim(AsString(name), " ")
	if b == "ON" {
		return true
	} else if b == "OFF" {
		return false
	} else {
		warnNonBool(b)
		return false
	}
}

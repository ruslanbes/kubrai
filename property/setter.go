package property

import "github.com/ruslanbes/kubrai/fileutils"

// SetProperties sets properties from the map
func SetProperties(props map[string]string) {
	for n, v := range props {
		setStringProperty(n, v)
	}
}

func setStringProperty(name, value string) {
	fileutils.FilePutContents(PropertyDir+"/"+name, value)
}

package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FilePutContents saves string to file
func FilePutContents(filename string, data string) {
	os.Mkdir(filepath.Dir(filename), 0777)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(fmt.Errorf("Error creating testfile %s", err))
		return
	}
	defer file.Close()

	file.WriteString(data)
}

// FileRemove removes file
func FileRemove(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println(fmt.Errorf("Error deleting testfile %s", err))
		return
	}
}

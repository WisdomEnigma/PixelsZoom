package dir

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const INTERNAL_PATH string = "app_data/"

func Chdir(path string) (*os.File, error) {

	if _, err := os.Stat(path); os.IsExist(err) {
		return &os.File{}, err
	}

	paths, err := os.Stat(INTERNAL_PATH)
	if err != nil {
		return &os.File{}, err
	}

	// Application storage path
	if !paths.IsDir() {
		return &os.File{}, err
	}

	// Create new file or open image file

	file, err := ioutil.TempFile(filepath.Dir(INTERNAL_PATH), "scale_*-"+path)
	if err != nil {
		log.Fatalln("Error :", err)
		return &os.File{}, err
	}

	return file, nil
}

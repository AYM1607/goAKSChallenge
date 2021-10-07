package common

import (
	"io/ioutil"
	"path/filepath"
)

func GetAllFilesInDir(dir string) ([]string, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	filePaths := []string{}
	for _, fi := range fis {
		filePaths = append(filePaths, filepath.Join(dir, fi.Name()))
	}
	return filePaths, nil
}

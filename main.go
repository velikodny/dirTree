////////////////////////////////////////////////////////
// Showing tree of directories
// params:
//		-f	 with a files and their total sizes
//

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func printTree(out io.Writer, path string, printFiles bool) error {

	err := readDir(out, path, printFiles, 0)

	if err != nil {
		return fmt.Errorf("Was error in ReadDir: %s", err)
	}

	return nil
}

func readDir(out io.Writer, path string, printFiles bool, finishName int) error {

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return fmt.Errorf("Error in readDir: %s", err)
	}

	for index, file := range files {

		var countFiles int

		if !file.IsDir() && !printFiles {
			continue
		}

		if !printFiles {
			countFiles = countDir(files)
		} else {
			countFiles = len(files) - 1
		}

		if file.Name() == ".DS_Store" {
			continue
		}

		spaceFormat := ""

		for spaceCount := strings.Count(path, "/") - finishName; spaceCount > 0; spaceCount-- {
			spaceFormat += "│\t"
		}

		for index := finishName; index > 0; index-- {
			spaceFormat += "\t"
		}

		sizeFormat := ""
		if !file.IsDir() {
			sizeFormat += fmt.Sprintf("%v", fileSizeToStr(file.Size()))
		}

		if index >= countFiles {
			fmt.Fprintf(out, "%s└───%s%s\n", spaceFormat, file.Name(), sizeFormat)
			finishName++
		} else {
			fmt.Fprintf(out, "%s├───%s%s\n", spaceFormat, file.Name(), sizeFormat)
		}

		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())
			readDir(out, newPath, printFiles, finishName)
		}

	}

	return nil
}

func countDir(files []os.FileInfo) (result int) {

	for _, file := range files {
		if file.IsDir() {
			result++
		}
	}

	return
}

func fileSizeToStr(size int64) (result string) {

	if size == 0 {
		result = " (empty)"
	} else {
		result = fmt.Sprintf(" (%vb)", size)
	}

	return
}

func main() {

	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := printTree(out, path, printFiles)

	if err != nil {
		panic(err.Error())
	}
}


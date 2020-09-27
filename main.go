package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")

	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	dirsFileList, err := listDir(path)
	fmt.Println(dirsFileList)
	if err != nil {
		return fmt.Errorf("")
	}
	return nil
}

func dirWalk(root string, dirList []string, prefix string, printFiles bool) error {
	if printFiles {
		var filePrefix string
		//if len(dirs) > 0 {
		//	filePrefix = prefix + "|"
		//} else {
		//	filePrefix = prefix + " "
		//}
		filePrefix = filePrefix + "\t"
	}

	return nil
}

func listDir(path string) (dirFileList []string, err error)  {
	absPath,err := filepath.Abs(path)
	fmt.Println(absPath)
	if err != nil {
		return dirFileList, fmt.Errorf("Directory not found")
	}
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return dirFileList, fmt.Errorf("Directory not found")
	}

	for _, file := range files {
		dirFileList = append(dirFileList, file.Name())
	}

	return dirFileList, nil
}
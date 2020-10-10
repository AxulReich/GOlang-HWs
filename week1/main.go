package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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
	absPath,err := filepath.Abs(path)
	fmt.Println(absPath)
	if err != nil {
		return err
	}
	dirFilesLits, err := listDir(absPath, printFiles)
	if err != nil {
		return err
	}

	err = dirWalk(dirFilesLits, output, path, "", printFiles)
	if err != nil {
		return err
	}
	return nil
}

func dirWalk(dirFileList []os.FileInfo, output io.Writer, root string, prefix string, printFiles bool) error {
	var dirfilePrefix string = prefix + "├───"
	var walkPrefix string = prefix + "│\t"

	for idx, file := range dirFileList {
		if idx == len(dirFileList) - 1 {
			dirfilePrefix = prefix + "└───"
			walkPrefix = prefix + "\t"
		}
		if file.IsDir() {
			_, err := fmt.Fprintln(output, dirfilePrefix + file.Name())

			if err != nil {
				return err
			}
			newPath := filepath.Join(root, file.Name())
			dirFileListInternal, err := listDir(newPath, printFiles)

			err = dirWalk(dirFileListInternal, output, newPath, walkPrefix, printFiles)

			if err != nil {
				return err
			}
		} else if printFiles {
			size := strconv.FormatInt(file.Size(), 10)
			if size == "0" {
				size = "empty"
			} else {
				size = size + "b"
			}

			_, err := fmt.Fprintln(output, dirfilePrefix + file.Name() + " (" + size + ")")

			if err != nil {
				return err
			}

		}
	}
	return nil
}


func listDir(path string, printFiles bool) (dirFileList []os.FileInfo, err error)  {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return dirFileList, err
	}

	if printFiles {
		for _, file := range files {
			dirFileList = append(dirFileList, file)
		}
	} else {
		for _, file := range files {
			if file.IsDir() {
				dirFileList = append(dirFileList, file)
			}
		}
	}


	return dirFileList, nil
}
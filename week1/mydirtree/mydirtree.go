package mydirtree

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func DirTree(out io.Writer, path string, printFiles bool) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dirFilesLists, err := listDir(absPath, printFiles)
	if err != nil {
		return err
	}

	return dirWalk(dirFilesLists, out, path, "", printFiles)
}

func dirWalk(dirFileList []os.FileInfo, out io.Writer, root string, prefix string, printFiles bool) error {
	dirfilePrefix := prefix + "├───"
	walkPrefix := prefix + "│\t"

	for idx, file := range dirFileList {
		if idx == len(dirFileList)-1 {
			dirfilePrefix = prefix + "└───"
			walkPrefix = prefix + "\t"
		}

		if file.IsDir() {
			if _, err := fmt.Fprintln(out, dirfilePrefix+file.Name()); err != nil {
				return err
			}

			newPath := filepath.Join(root, file.Name())
			dirFileListInternal, err := listDir(newPath, printFiles)
			if err != nil {
				return err
			}

			err = dirWalk(dirFileListInternal, out, newPath, walkPrefix, printFiles)
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

			if _, err := fmt.Fprintln(out, dirfilePrefix+file.Name()+" ("+size+")"); err != nil {
				return err
			}
		}
	}
	return nil
}

func listDir(path string, printFiles bool) (dirFileList []os.FileInfo, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return dirFileList, err
	}

	if printFiles {
		return files, nil
	}

	for _, file := range files {
		if file.IsDir() {
			dirFileList = append(dirFileList, file)
		}
	}

	return
}

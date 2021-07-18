package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func sortDirs(paths []os.DirEntry) []os.DirEntry {

	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i].Name() < paths[j].Name()
	})

	return paths
}

func removeNotDirs(paths []os.DirEntry) []os.DirEntry {
	arr := []os.DirEntry{}

	for _, node := range paths {
		if node.IsDir() {
			arr = append(arr, node)
		}
	}

	return arr
}

func getLineAndPrefix(file os.DirEntry, currentPrefix string, isLast bool) (string, string, error) {
	var nodeSymbol string
	var prefix string
	var err error
	if !isLast {
		nodeSymbol = "├───"
		prefix = currentPrefix + "│\t"
	} else {
		nodeSymbol = "└───"
		prefix = currentPrefix + "\t"
	}

	lineToPrint := currentPrefix + nodeSymbol + file.Name()

	if !file.IsDir() {
		fsize, fileErr := file.Info()

		if fileErr != nil {
			err = fileErr
			return "", "", err
		}

		if fsize.Size() == 0 {
			lineToPrint += " (empty)"
		} else {
			lineToPrint += fmt.Sprintf(" (%db)", fsize.Size())
		}

	}

	lineToPrint += "\n"

	return lineToPrint, prefix, nil

}

func printDirs(out io.ReadWriter, path string, printFiles bool, currentPrefix string) error {
	var err error
	files, fileErr := os.ReadDir(path)

	if fileErr != nil {
		err = fileErr
		return err
	}

	files = sortDirs(files)
	if !printFiles {
		files = removeNotDirs(files)
	}

	for idx, file := range files {
		lineToPrint, prefix, lineErr := getLineAndPrefix(file, currentPrefix, idx == len(files)-1)

		if lineErr != nil {
			err = lineErr
			return err
		}
		out.Write([]byte(lineToPrint))

		if file.IsDir() {
			nextPath := path + string(os.PathSeparator) + file.Name()
			dirErr := printDirs(out, nextPath, printFiles, prefix)

			if dirErr != nil {
				err = dirErr
				return err
			}

		}
	}
	return err
}

func dirTree(out io.ReadWriter, path string, printFiles bool) error {
	var err error

	printErr := printDirs(out, path, printFiles, "")

	if printErr != nil {
		err = printErr
	}
	return err
}

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

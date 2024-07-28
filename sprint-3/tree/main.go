package main

/*
Course `Web services on Go`, week 1, homework, `tree` program.
See: week_01\materials.zip\week_1\99_hw\tree

mkdir -p week01_homework/tree
pushd week01_homework/tree
go mod init tree
go mod tidy
pushd ..
go work init
go work use ./tree/
go vet tree
gofmt -w tree
go test -v tree
go run tree . -f
go run tree ./tree/testdata
cd tree && docker build -t mailgo_hw1 .

https://en.wikipedia.org/wiki/Tree_(command)
https://mama.indstate.edu/users/ice/tree/
https://stackoverflow.com/questions/32151776/visualize-tree-in-bash-like-the-output-of-unix-tree

*/

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

/*
	Example output:

	├───project
	│	└───gopher.png (70372b)
	├───static
	│	├───a_lorem
	│	│	├───dolor.txt (empty)
	│	├───css
	│	│	└───body.css (28b)
	...
	│			└───gopher.png (70372b)

	- path should point to a directory,
	- output all dir items in sorted order, w/o distinction file/dir
	- last element prefix is `└───`
	- other elements prefix is `├───`
	- nested elements aligned with one tab `	` for each level
*/

const (
	EOL             = "\n"
	BRANCHING_TRUNK = "├───"
	LAST_BRANCH     = "└───"
	TRUNC_TAB       = "│\t"
	LAST_TAB        = "\t"
	EMPTY_FILE      = "empty"
	ROOT_PREFIX     = ""

	USE_RECURSION_ENV_KEY = "RECURSIVE_TREE"
	USE_RECURSION_ENV_VAL = "YES"
)

func main() {
	// This code is given
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage: go run main.go . [-f]")
	}

	out := os.Stdout
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

// dirTree: `tree` program implementation, top-level function, signature is fixed.
// Write `path` dir listing to `out`. If `prinFiles` is set, files is listed along with directories.
func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeOneLevel(out, path, ROOT_PREFIX, printFiles)
}

func dirTreeOneLevel(out io.Writer, path, prefix string, printFiles bool) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("dirTreeOneLevel: %v", err)
	}

	files = preprocessingFiles(files, !printFiles)
	indexLastBranch := getLastBranch(files, printFiles)

	var branch, tab string

	for i, file := range files {
		if i == indexLastBranch {
			branch = LAST_BRANCH
			tab = LAST_TAB
		} else {
			branch = BRANCHING_TRUNK
			tab = TRUNC_TAB
		}

		record, err := getRecord(file, prefix, branch)
		if err != nil {
			return fmt.Errorf("dirTreeOneLevel: %v", err)
		}

		out.Write([]byte(record))

		if file.IsDir() {
			nextPath := path + string(os.PathSeparator) + file.Name()
			if err = dirTreeOneLevel(out, nextPath, prefix+tab, printFiles); err != nil {
				return fmt.Errorf("dirTreeOneLevel: %v", err)
			}
		}
	}

	return nil
}

func preprocessingFiles(files []os.DirEntry, deleteRegularFiles bool) []os.DirEntry {
	if deleteRegularFiles {
		var dirs []os.DirEntry
		for _, file := range files {
			if file.IsDir() {
				dirs = append(dirs, file)
			}
		}

		files = dirs
	}

	slices.SortFunc(files, func(a, b os.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return files
}

func getLastBranch(files []os.DirEntry, printFiles bool) int {
	if printFiles {
		return len(files) - 1
	}

	lastBranch := len(files) - 1
	for i, file := range files {
		if file.IsDir() {
			lastBranch = i
		}
	}
	return lastBranch
}

func getRecord(file os.DirEntry, prefix, branch string) (string, error) {
	if file.IsDir() {
		return strings.Join([]string{prefix, branch, file.Name(), EOL}, ""), nil
	}

	record := prefix + branch + file.Name()

	info, err := file.Info()
	if err != nil {
		return "", fmt.Errorf("getRecord: %v", err)
	}

	if info.Size() == 0 {
		record = fmt.Sprintf("%s (%s)", record, EMPTY_FILE)
	} else {
		record = fmt.Sprintf("%s (%db)", record, info.Size())
	}

	record += EOL

	return record, nil
}

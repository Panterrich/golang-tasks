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

type GraphicDirWorker struct {
	out        io.Writer
	printFiles bool

	indexLastBranch int
	prefix          string
	branch          string
	tab             string
}

var _ DirWorker = (*GraphicDirWorker)(nil)

func NewGraphicDirWorker(out io.Writer, printFiles bool) GraphicDirWorker {
	return GraphicDirWorker{
		out:        out,
		printFiles: printFiles,
	}
}

func (w *GraphicDirWorker) Clone() DirWorker {
	return &GraphicDirWorker{
		out:        w.out,
		printFiles: w.printFiles,

		indexLastBranch: w.indexLastBranch,
		prefix:          w.prefix,
		branch:          w.branch,
		tab:             w.tab,
	}
}

func (w *GraphicDirWorker) ProcessFiles(files []os.DirEntry) ([]os.DirEntry, error) {
	if !w.printFiles {
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

	w.indexLastBranch = getLastBranch(files, w.printFiles)
	w.prefix += w.tab

	return files, nil
}

func (w *GraphicDirWorker) Work(index int, file os.DirEntry) error {
	if index == w.indexLastBranch {
		w.branch = LAST_BRANCH
		w.tab = LAST_TAB
	} else {
		w.branch = BRANCHING_TRUNK
		w.tab = TRUNC_TAB
	}

	record, err := w.getRecord(file)
	if err != nil {
		return fmt.Errorf("get record: %v", err)
	}

	if _, err := w.out.Write([]byte(record)); err != nil {
		return fmt.Errorf("write: %v", err)
	}

	return nil
}

func (w *GraphicDirWorker) getRecord(file os.DirEntry) (string, error) {
	if file.IsDir() {
		return strings.Join([]string{w.prefix, w.branch, file.Name(), EOL}, ""), nil
	}

	record := strings.Join([]string{w.prefix, w.branch, file.Name()}, "")

	info, err := file.Info()
	if err != nil {
		return "", fmt.Errorf("file %s info : %v", file.Name(), err)
	}

	if info.Size() == 0 {
		record = fmt.Sprintf("%s (%s)", record, EMPTY_FILE)
	} else {
		record = fmt.Sprintf("%s (%db)", record, info.Size())
	}

	record += EOL

	return record, nil
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
// Write `path` dir listing to `out`. If `printFiles` is set, files is listed along with directories.
func dirTree(out io.Writer, path string, printFiles bool) error {
	walker := NewPreOrderDirWalker()
	worker := NewGraphicDirWorker(out, printFiles)

	return walker.Walk(path, &worker)
}

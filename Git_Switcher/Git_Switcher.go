package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	// "os/exec"
	"strconv"
	"strings"
)

const ()

// Use %CD% as .cmd argument

func main() {
	args := os.Args
	cwd := args[1]

	gitFolders := allGitFolders(cwd)
	fmt.Println(cwd)
	fmt.Println(gitFolders)

	b, err := ioutil.ReadFile(cwd + "/.git/gitname.txt") // just pass the file name
	checkError(err)

	fmt.Println("Current git project: ", string(b))
	fmt.Println("All git projects in this repo: ")
	if len(gitFolders) == 0 {
		fmt.Println("There are no other git projects in this directory.")
		return
	}
	for index, folder := range gitFolders {
		name, err := ioutil.ReadFile(folder + "/gitname.txt")
		checkError(err)
		fmt.Println(strconv.Itoa(index) + ". " + string(name))
	}

	fmt.Println("Type the number of git project you'd like to switch to. Type anything else to exit.")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	input, err := strconv.Atoi(text)
	if err != nil {
		return
	}

	if input >= 0 && input < len(gitFolders) {
		oldName := gitFolders[input]
		os.Rename(oldName, "temp")
		os.Rename(".git", oldName)
		os.Rename("temp", ".git")
		copyFileContents(cwd+"/.git/ReadMe.txt", cwd+"/ReadMe.txt")
		copyFileContents(cwd+"/.git/.gitignore", cwd+"/.gitignore")
		copyFileContents(cwd+"/.git/.gitattributes", cwd+"/.gitattributes")

	}
}

// ------------------------------------------- Private ------------------------------------------- //

func allGitFolders(cwd string) []string {
	dir, err := ioutil.ReadDir(cwd)
	checkError(err)
	output := make([]string, 0)
	for _, item := range dir {
		name := item.Name()
		if name[:4] == ".git" && name != ".gitignore" && name != ".gitattributes" {
			output = append(output, name)
		}
	}
	return output
}

// ------------------------------------------- Utilities ------------------------------------------- //

// Respond to errors
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// https://stackoverflow.com/a/21067803/9463878
func copyFileContents(src string, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

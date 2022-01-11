package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Checks if the file is a valid c or cpp file
func isValidFile(f os.FileInfo) bool {
	return !f.IsDir() && strings.Contains(f.Name(), ".c") && strings.Count(f.Name(), ".") == 1
}

func getAllCodeFiles() {
	files_compiled := 0

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		// Run if file is not a directory
		if isValidFile(f) {
			var cmd *exec.Cmd
			var output_name string

			if strings.Contains(f.Name(), "cpp") {
				output_name = strings.TrimSuffix(f.Name(), ".cpp")
				cmd = exec.Command("g++", f.Name(), "-o", output_name)
			} else {
				output_name = strings.TrimSuffix(f.Name(), ".c")
				cmd = exec.Command("gcc", f.Name(), "-o", output_name)
			}

			fmt.Println("Executing:", cmd)
			cmd.Stderr = os.Stderr
			cmd.Run()

			exec.Command("mv", output_name, "output/").Run()

			files_compiled++
		}
	}

	println("Compiled", files_compiled, "files")
}

func createOutputFolder() {
	exec.Command("mkdir", "output").Run()
}

func main() {
	createOutputFolder()
	getAllCodeFiles()
	fmt.Print("The end")

}

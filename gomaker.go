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

func compileFiles() {
	files_compiled := 0

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if isValidFile(file) {
			var cmd *exec.Cmd
			var output_name string

			if strings.Contains(file.Name(), "cpp") {
				output_name = strings.TrimSuffix(file.Name(), ".cpp")
				cmd = exec.Command("g++", file.Name(), "-o", output_name)
			} else {
				output_name = strings.TrimSuffix(file.Name(), ".c")
				cmd = exec.Command("gcc", file.Name(), "-o", output_name)
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
	compileFiles()
	fmt.Print("Compilation complete.")

}

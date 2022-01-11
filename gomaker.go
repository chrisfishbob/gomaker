package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	// Import go routines
	"sync"
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

	// Create a waitgroup
	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		go func(file os.FileInfo) {
			// We defer wg.Done() to decrement waitgroup regarless of validity of name
			defer wg.Done()
			
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
		}(file)
	}

	wg.Wait()

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

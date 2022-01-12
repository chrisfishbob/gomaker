package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Checks if the file is a valid c or cpp file
func isValidFile(f os.FileInfo) bool {
	return !f.IsDir() && strings.Contains(f.Name(), ".c") && strings.Count(f.Name(), ".") == 1
}

func processFiles() {
	start := time.Now()
	files_compiled := 0
	files_skipped := 0
	files_skipped_string := ""

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
				runCompileCommand(file, &files_compiled)
			}else{
				files_skipped_string += file.Name() + "\n"
				files_skipped++
			}
		}(file)
	}

	// Gathering all the goroutines
	wg.Wait()  
	end := time.Now()
	fmt.Println("\nCompiled", files_compiled, "files in", end.Sub(start).Seconds(), "seconds")
	fmt.Println("Skipped", files_skipped, "files: ")
	fmt.Println(files_skipped_string)
}

func runCompileCommand(file os.FileInfo, files_compiled *int) {
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
	// Ensures that warnings and errors are printed
	cmd.Stderr = os.Stderr
	cmd.Run()

	cwd, _ := os.Getwd()
	exec.Command("mv", output_name, cwd + "/output").Run()

	(*files_compiled)++
}


func createOutputFolder() {
	os.Mkdir("output", os.ModePerm)
}

func main() {
	createOutputFolder()
	processFiles()
	fmt.Print("Compilation complete.")
}

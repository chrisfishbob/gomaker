package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {

			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func unzipToCurrentDirectory() {
	files, err := ioutil.ReadDir(".")
	pwd, _ := os.Getwd()

	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".zip") {
			Unzip(file.Name(), pwd)
		}
	}
}

func extractFolders() {
	cmd_string := "find  . -mindepth 2 -type f -exec mv {} . \\;"
	//call the cmd_string command with bash
	cmd := exec.Command("bash", "-c", cmd_string)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// Checks if the file is a valid c or cpp file
func isValidFile(f os.FileInfo) bool {
	return !f.IsDir() && strings.Contains(f.Name(), ".c") && strings.Count(f.Name(), ".") == 1
}

func printExitInformation(string_slice *[]string,
	stderr_slice *[]string,
	files_skipped int,
	files_compiled int,
	start time.Time,
	end time.Time,
	files_skipped_string string) {

	// Use the slices to print out the files that did not compile perfectly
	for i := 0; i < len(*string_slice); i++ {
		fmt.Println("\n\n"+(*string_slice)[i], "compiled with the following warnings:")
		fmt.Println((*stderr_slice)[i])
	}

	fmt.Println("The following compiled with warnings / errors:")

	for i := 0; i < len(*string_slice); i++ {
		fmt.Println((*string_slice)[i])
	}

	fmt.Println("")
	fmt.Println("Skipped", files_skipped, "files: ")
	fmt.Println(files_skipped_string)
	fmt.Println("Compiled", files_compiled, "files in", end.Sub(start).Seconds(), "seconds")
}

func processFiles(additional_flags string) {
	start := time.Now()
	files_compiled := 0
	files_skipped := 0
	files_skipped_string := ""
	string_slice := make([]string, 0)
	stderr_slice := make([]string, 0)

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
				runCompileCommand(file, &files_compiled, &string_slice, &stderr_slice, additional_flags)
			} else {
				files_skipped_string += file.Name() + "\n"
				files_skipped++
			}
		}(file)
	}

	// Gathering all the goroutines
	wg.Wait()
	end := time.Now()

	printExitInformation(
		&string_slice,
		&stderr_slice,
		files_skipped,
		files_compiled,
		start, end,
		files_skipped_string)
}

func runCompileCommand(file os.FileInfo, files_compiled *int, string_slice *[]string, stderr_slice *[]string, additional_flags string) {
	var cmd *exec.Cmd
	var output_name string

	if strings.Contains(file.Name(), "cpp") {
		output_name = strings.TrimSuffix(file.Name(), ".cpp")
		if additional_flags == "none" {
			cmd = exec.Command("g++", file.Name(), "-fdiagnostics-color=always", "-o", output_name)
		} else {
			cmd = exec.Command("g++", file.Name(), "-fdiagnostics-color=always", "-o", output_name, additional_flags)
		}
	} else {
		output_name = strings.TrimSuffix(file.Name(), ".c")

		if additional_flags == "none" {
			cmd = exec.Command("gcc", file.Name(), "-fdiagnostics-color=always", "-o", output_name)
		} else {
			cmd = exec.Command("gcc", file.Name(), "-fdiagnostics-color=always", "-o", output_name, additional_flags)
		}
	}

	// Ensures that warnings and errors are printed
	// Override stderr to a buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Run()
	// Check if the buffer is empty
	if stderr.String() == "" {
		fmt.Println("Compiled", file.Name(), "with no compiler warnings")
	} else {
		*string_slice = append(*string_slice, file.Name())
		*stderr_slice = append(*stderr_slice, stderr.String())
	}

	cwd, _ := os.Getwd()
	exec.Command("mv", output_name, cwd+"/output").Run()

	(*files_compiled)++
}

func createOutputFolder() {
	os.Mkdir("output", os.ModePerm)
}

func confirmRun() {
	cwd, _ := os.Getwd()
	fmt.Println("gomaker is about to execute at", cwd, "are you sure you want to continue? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

}

func removeEmptyDirectories() {
	cmd_string := "find . -type d -empty -delete"
	//call the cmd_string command with bash
	cmd := exec.Command("bash", "-c", cmd_string)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	additional_flags := "none"

	//take in the command line flags
	var frFlag = flag.Bool("fr", false, "Flatten folders recursively")
	var yFlag = flag.Bool("y", false, "Skip confirmation prompt")
	var zFlag = flag.Bool("z", false, "Unzips all .zip files")
	var fFlag = flag.Bool("f", false, "Additional flags for compilation")

	flag.Parse()

	if !*yFlag {
		confirmRun()
	}
	if *zFlag {
		unzipToCurrentDirectory()
	}
	if *frFlag {
		extractFolders()
	}

	if *fFlag {
		// Prompt user for additonal flag
		fmt.Println("Please enter additional flags for compilation:")
		fmt.Scanln(&additional_flags)
	}

	createOutputFolder()
	processFiles(additional_flags)
	removeEmptyDirectories()
	fmt.Println("Compilation complete.")
}

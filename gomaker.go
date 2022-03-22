package main

import (
	"archive/zip"
	"bufio"
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

// Global mutex for printing with threads
var mutex sync.Mutex

// Returns true if the user
func usedBannedKeyword(filename string, banned_words []string) bool {
	current_line_number := 0
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		for _, banned_word := range banned_words {
			if strings.Contains(scanner.Text(), banned_word) {
				mutex.Lock()
				fmt.Println("\n", filename, "FAILED style test")
				fmt.Print("\tBanned keyword: ", banned_word, " used on line ",
					current_line_number, "\n\n")
				mutex.Unlock()
				return true
			}
		}

		current_line_number++
	}

	return false
}

func functionLengthUnderLimit(filename string, limit int) bool {
	var currentLine string
	openBraceCount := 0
	closeBraceCount := 0
	functionLength := 0
	blockStartLine := 0
	currentLineNumber := 0

	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		currentLineNumber++
		currentLine = scanner.Text()
		// trim all trailing whitespace
		currentLine = strings.TrimRight(currentLine, " \t")

		if strings.Contains(currentLine, "{") {
			// If it's the first open brace, it marks the beginning of a block
			if openBraceCount == 0 {
				blockStartLine = currentLineNumber
			}

			openBraceCount++
		}

		if strings.Contains(currentLine, "}") {
			closeBraceCount++
		}

		// If open == close, then we're at the end of a function
		if openBraceCount == closeBraceCount {
			// We check limit - 2 because we don't count the braces
			if functionLength-2 > limit {
				// Lock mutex
				mutex.Lock()
				fmt.Println("\n", filename, "FAILED style test")
				fmt.Println("\tBlock at line", blockStartLine, "is too long")
				fmt.Println("\tThe block is", functionLength-2, "long")
				fmt.Print("\tThe limit is ", limit, "\n\n")
				mutex.Unlock()
				return false
			} else {
				functionLength = 0
				closeBraceCount = 0
				openBraceCount = 0
				blockStartLine = 0
			}
		}

		functionLength++
	}

	file.Close()

	return true
}

func underLineLimit(filename string, limit int) bool {
	var currentLine string
	currentLineNumber := 1
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		currentLine = scanner.Text()
		// trim all trailing whitespace
		currentLine = strings.TrimRight(currentLine, " \t")

		if len([]rune(currentLine)) > limit {
			mutex.Lock()
			fmt.Println("\n", filename, "FAILED style test")
			fmt.Println("\tLine", currentLineNumber, "has", len([]rune(currentLine)), "characters.")
			fmt.Print("\tThe limit is ", limit, "\n\n")
			mutex.Unlock()
			return false
		}

		currentLineNumber++
	}

	file.Close()
	return true
}

// Does not return any value, prints out style errors in terminal
// Unlike the style check during compilation, this does not terminate
// once an error is found
func runStyleCheckOnly(function_lines_limit int, line_char_limit int, pedantic bool, banned_words []string) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if isValidFile(file) {
			underLineLimit(file.Name(), line_char_limit)
			functionLengthUnderLimit(file.Name(), function_lines_limit)

			if pedantic {
				usedBannedKeyword(file.Name(), banned_words)
			}

		}
	}
}

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

func shouldCompile(file os.FileInfo,
	check_for_style bool,
	function_length_limit int,
	line_length_limit int,
	banned_words []string) bool {

	// If the file is a valid C/C++ file and we need to check for style and
	// the file passes the style check, then we should compile
	if isValidFile(file) && 
		check_for_style && 
		functionLengthUnderLimit(file.Name(), function_length_limit) && 
		underLineLimit(file.Name(), line_length_limit) && 
		!usedBannedKeyword(file.Name(), banned_words) {

		return true
	} 

	// If the file is a valid C/C++ file and we do not need to check for style,
	// we should also compile
	if isValidFile(file) && !check_for_style {
		return true
	}

	return false		
}

func processFiles(additional_flags string, check_for_style bool,
	function_line_limit int, chars_per_line_limit int, banned_words []string) {
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

			// Check if the file should be compiled
			if shouldCompile(file, check_for_style, function_line_limit, chars_per_line_limit, banned_words) {
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
	check_for_style := false
	function_lenth_limit := 0
	characters_per_line_limit := 0
	var banned_words string
	var banned_words_slice []string
	bufio := bufio.NewReader(os.Stdin)

	//take in the command line flags
	var frFlag = flag.Bool("fr", false, "Flatten folders recursively")
	var yFlag = flag.Bool("y", false, "Enable confirmation prompt")
	var zFlag = flag.Bool("z", false, "Unzips all .zip files")
	var fFlag = flag.Bool("f", false, "Additional flags for compilation")
	var sFlag = flag.Bool("s", false, "Enable style check")
	var styleOnlyFlag = flag.Bool("styleonly", false, "Only do style check, no compilation")
	var pedanticFlag = flag.Bool("pedantic", false, "Bans certain feaures of C/C++")

	flag.Parse()

	// If pedantic, then we need to ask for the banned words
	if *pedanticFlag {
		// TODO: Enable pedantic flag with normal compilation
		fmt.Print("Please enter the banned words: ")
		banned_words, _ = bufio.ReadString('\n')
		banned_words = strings.TrimSuffix(banned_words, "\n")
		banned_words_slice = strings.Split(banned_words, ", ")
	}

	if *styleOnlyFlag {
		fmt.Print("Please enter the function length limit: ")
		fmt.Scanln(&function_lenth_limit)
		fmt.Print("Please enter the characters per line limit: ")
		fmt.Scanln(&characters_per_line_limit)

		runStyleCheckOnly(function_lenth_limit, characters_per_line_limit, *pedanticFlag, banned_words_slice)

		return
	}

	if *yFlag {
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
		fmt.Print("Please enter additional flags for compilation: ")
		fmt.Scanln(&additional_flags)
	}

	if *sFlag {
		check_for_style = true
		fmt.Print("Please enter the function length limit: ")
		fmt.Scanln(&function_lenth_limit)
		fmt.Print("Please enter the characters per line limit: ")
		fmt.Scanln(&characters_per_line_limit)
	}

	createOutputFolder()
	processFiles(additional_flags, check_for_style, function_lenth_limit, characters_per_line_limit, banned_words_slice)
	removeEmptyDirectories()
	fmt.Println("Compilation complete.")
}

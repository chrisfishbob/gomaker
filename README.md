<img src="https://github.com/chrisfishbob/gomaker/blob/main/gophers.png" width="430" height="430" align="left"/> 

# **Gomaker**
A lightning-fast multi-threaded concurrent compiling tool with customizable style checking

## Purpose
Gomaker provides the ability to rapidly compile hundreds of C/C++ source files rapidly through concurrency. <br/><br/>
Gomaker also provides basic style checking, either before compilation or independently. <br/><br/>
The main use case of Gomaker is to compile large quantities of independent source files that all do a similar task. (e.x. A university assignment) <br/>

## Dependencies:
Go and Python3 is required for the compilation and install process. GCC/G++ is required during runtime. <br/><br/>
```
sudo apt install python3 gcc g++ go
```
<br/><br/>

## Build
### Using the provided install script
```
python3 install.py
```
Enter the password if prompted. <br/><br/>
The compiled binary will be installed at /usr/local/bin. <br/><br/><br/>

## Basic Usage
### To simply compile all the C/C++ files in the current directory:
```
./gomaker
```
This command will compile all the C/C++ files and place it in the "output" folder. <br/><br/>
Files that compiled smoothly, compiled with warning, and failed to compile are displayed to the terminal in their own sections. Along with
the warning/error message, if applicable. <br/><br/>
Files that were skipped (i.e. not a valid C/C++ file) will be skipped and will also be displayed to the terminal. <br/><br/><br/>

### To compile with basic style check:
```
./gomaker -s
```
The program will prompt the user to enter two parameters: How long each function block can be, and character limit per line. <br/><br/>
After the prompt, compiling with style check behaves almost identically as the basic use case, with the exception that files are checked for style compliance before compilation begins. <br/><br/>
The style checker will stop checking a given file when one violation is found and move on to the next file. <br/><br/><br/>

### To compile with stict style check:
```
./gomaker -s -pedantic
```
The `-pedantic` flag adds an additional prompt for the user to enter banned keywords/features (e.x. using namespace std) <br/><br/>
Note that the pedantic option cannot be used without `-s`. <br/><br/><br/>


### To compile with additional compiler flags:
```
./gomaker -f
```
The `-f` flag adds an additional prompt for the user to enter additonal compiler flags, most often used for linking. <br/><br/>
The additonal flag only applies when the source file compiles, the `-s` and `-pedantic` flags are unaffected and can be used in conjunction with
`-f` to combine their effects.<br/><br/><br/>



### To simply check for style:
```
./gomaker -styleonly
```
Checks for style compliance and reports violations in terminal. No files will be compiled. <br/><br><br/>


## Advanced Usage:
### To flatten all folders recursively before compiling:
```
./gomake -fr -y
```
If the target directory contains folders that encloses the source files (e.x. when some student submit souce files directly, while others submit their files in a folder), adding the `-fr` flag will retract all the files before compilation begins like normal. <br/><br/>

The `-y` flag enables a confirmation prompt before the program executes.<br/><br/>

Warning: It is highly recommend that the `-rf` flag is always accompanied by `-y` as flattening directories could be destructive if done in the wrong path. <br/><br/> <br/>


### To unzip all zipfiles before compiling:
```
./gomake -z
```
The `-z` flag will extract all the contents of any .zip files before compilation begins like normal. <br/><br/>
Combining the `-fr` and `-z` flag is possible, the archives will first be unzipped before source extraction begins.<br/><br/><br/>

## Features and Progress:
- :white_check_mark: Compiles all source files in one thread
- :white_check_mark: Compiles all source files concurrently via Goroutines
- :white_check_mark: Option to extract source files recursively
- :white_check_mark: Option to unzip .zip files
- :white_check_mark: Option to add additional compiler flags
- :white_check_mark: Character line limit check
- :white_check_mark: Function line limit check
- :white_check_mark: Banned feature usage check
- :white_check_mark: Install script
- 🔷 Colored results for style check (in progress)


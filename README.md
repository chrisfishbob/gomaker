<img src="https://github.com/chrisfishbob/gomaker/blob/main/gophers.png" width="430" height="430" align="left"/> 

# **Gomaker**
A lightning-fast multi-threaded concurrent compiling tool with customizable style checking

## Purpose
Gomaker provides the ability to rapidly compile hundreds of C/C++ source files rapidly through concurrency. <br/><br/>
Gomaker also provides basic style checking, either before compilation or independently. <br/><br/>
The main use case of Gomaker is to compile large quantities of independent source files that all do a similar task. (e.x. A university assignment) <br/>


## Build
### Using the provided Makefile
```
make
```
<br/>

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

### To simply check for style:
```
./gomaker -styleonly
```
Checks for style compliance and reports violations in terminal. No files will be compiled. <br/><br>


### 

# substr: string substitution utility ![build](https://img.shields.io/badge/build-pass-green?style=flat-square&logo=okay)

A commandline utility for replacing text in files. It was designed to quickly and easily find and replace strings of text,
making it an ideal tool for finding and replacing text in large amounts of files. 
The utility supports both single and multiple file replacements, as well as replacements in entire directories. Substr is
written in Golang and is open-sourced under the MIT License.

### Prerequisites

- Go 1.14 or higher must be installed on your system.

### Using go install

1. Open a terminal window and navigate to the directory where you want to install `substr`.
2. Run the following command to install `substr`:
  ```
  go install github.com/oyamo/substr@latest
  ```
3. Verify installation
 ```
 $ substr -version
 ```


## Usage 
```shell
Usage:
  substr  --original-text <text> --substitute-text <text> [FILES]

Options:
  -h, --help            Display this help message and exit
  -v, --version         Display version information and exit
  -t, --original-text   Text to replace
  -s, --substitute-text Final text
  -o, --output-file     File path to redirect output
  -d, --output-dir      Dir path to redirect output from all files
```

## Examples
```shell
Examples:
  substr -t foo -s bar example.txt   # replaces all occurences of 'foo' with 'bar' in example.txt
  substr -t foo -s bar *.txt         # replaces all occurences of 'foo' with 'bar' in all .txt files in the current directory
  substr -t foo -s bar -o output.txt  # replaces all occurences of 'foo' with 'bar' in all input files and writes the output to output.txt
  substr -t foo -s bar -d output/    # replaces all occurences of 'foo' with 'bar' in all input files and writes the output to the specified directory

```


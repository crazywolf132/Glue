# Glue

**Glue** is a CLI tool that combines multiple project files into a single text file and can also reverse the process by extracting the files back from the combined text file. This is particularly useful when you need to share your project's code with AI assistants or other tools that require a single input file.

```
  ▄████  ██▓     █    ██ ▓█████ 
 ██▒ ▀█▒▓██▒     ██  ▓██▒▓█   ▀ 
▒██░▄▄▄░▒██░    ▓██  ▒██░▒███   
░▓█  ██▓▒██░    ▓▓█  ░██░▒▓█  ▄ 
░▒▓███▀▒░██████▒▒▒█████▓ ░▒████▒
 ░▒   ▒ ░ ▒░▓  ░░▒▓▒ ▒ ▒ ░░ ▒░ ░
  ░   ░ ░ ░ ▒  ░░░▒░ ░ ░  ░ ░  ░
░ ░   ░   ░ ░    ░░░ ░ ░    ░   
      ░     ░  ░   ░        ░  ░
                                 
```

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Building from Source](#building-from-source)
- [Usage](#usage)
    - [Combining Files](#combining-files)
    - [Extracting Files](#extracting-files)
    - [Flags and Options](#flags-and-options)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Introduction

When working with AI models or other tools that require code input, it can be cumbersome to copy and paste multiple files or manage large codebases. **Glue** simplifies this process by:

- Combining all your project files into a single text file, with clear separation between files.
- Respecting your `.gitignore` file to exclude unnecessary files.
- Allowing custom inclusion and exclusion patterns via glob patterns.
- Reconstructing the original project structure from the combined file.

## Features

- **Combine Files**: Merge multiple files into a single `combined.txt` file.
- **Reverse Operation**: Extract files from a combined file to recreate the project structure.
- **Glob Patterns**: Use glob-like patterns to include or exclude specific files.
- **.gitignore Support**: Automatically ignores files specified in `.gitignore`.
- **Customizable Output**: Specify the output file name.
- **ASCII Logo**: Displays a cool ASCII logo when running the tool.

## Installation

### Prerequisites

- **Go**: You need to have Go installed (version 1.16 or later). You can download it from [the official website](https://golang.org/dl/).

### Building from Source

1. **Clone the Repository**

   ```bash
   git clone https://github.com/crazywolf132/glue.git
   cd glue
   ```

2. **Fetch Dependencies**

   ```bash
   go mod tidy
   ```

3. **Build the Application**

   ```bash
   go build -o glue
   ```

   Or, use the provided `Makefile`:

   ```bash
   make build
   ```

4. **Install the Application (Optional)**

   ```bash
   make install
   ```

   This will install `glue` into your `$GOPATH/bin` directory.

## Usage

### Combining Files

By default, `glue` will:

- Find all files in the current directory and subdirectories.
- Exclude files specified in `.gitignore`.
- Combine the files into `combined.txt`.

```bash
./glue
```

### Extracting Files

To reverse the operation and recreate the project from a combined file:

```bash
./glue -r
```

This will read from `combined.txt` and reconstruct the files.

### Flags and Options

- `-a`, `--all`: Include all files, ignoring `.gitignore`.

  ```bash
  ./glue -a
  ```

- `-i`, `--ignore`: Glob patterns to ignore.

  You can specify multiple ignore patterns.

  ```bash
  ./glue -i "**/*.test.js" -i "node_modules/**"
  ```

- `-o`, `--output`: Specify the output file name. Default is `combined.txt`.

  ```bash
  ./glue -o my_combined_file.txt
  ```

- `-r`, `--reverse`: Reverse the operation, recreate project from the combined file.

  ```bash
  ./glue -r -o my_combined_file.txt
  ```

- **Inclusion Patterns**: You can specify glob patterns to include specific files. Provide these patterns as arguments.

  ```bash
  ./glue "**/*.go" "**/*.md"
  ```

### Examples

#### Combine All Files (Default)

```bash
./glue
```

#### Combine All Files, Ignoring `.gitignore`

```bash
./glue -a
```

#### Combine Specific Files Using Glob Patterns

Include only `.go` and `.md` files:

```bash
./glue "**/*.go" "**/*.md"
```

#### Combine Files, Ignoring Specific Patterns

Ignore test files:

```bash
./glue -i "**/*_test.go"
```

#### Specify Output File Name

```bash
./glue -o project_code.txt
```

#### Reverse Operation with Custom Output File

```bash
./glue -r -o project_code.txt
```

## How It Works

When combining files, `glue` walks through the file system and collects files based on the inclusion patterns provided (or all files if none are provided). It then excludes any files that match the patterns in `.gitignore` (unless `-a` is specified) and any additional patterns provided via `-i`.

Each file is written to the output file in the following format:

```
-- path/to/file.ext
```ext
<file contents>
```
```

When reversing the operation, `glue` reads the combined file, parses it into sections based on the format above, and writes each section back to its original file path.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

1. **Fork the repository.**
2. **Create a new branch:**

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes.**
4. **Commit your changes:**

   ```bash
   git commit -am 'Add some feature'
   ```

5. **Push to the branch:**

   ```bash
   git push origin feature/your-feature-name
   ```

6. **Submit a pull request.**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Repository:** [https://github.com/crazywolf132/glue](https://github.com/crazywolf132/glue)

If you have any questions or need further assistance, feel free to open an issue or contact the maintainer.
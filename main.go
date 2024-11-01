package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/urfave/cli/v2"
)

const logo = `
  ▄████  ██▓     █    ██ ▓█████ 
 ██▒ ▀█▒▓██▒     ██  ▓██▒▓█   ▀ 
▒██░▄▄▄░▒██░    ▓██  ▒██░▒███   
░▓█  ██▓▒██░    ▓▓█  ░██░▒▓█  ▄ 
░▒▓███▀▒░██████▒▒▒█████▓ ░▒████▒
 ░▒   ▒ ░ ▒░▓  ░░▒▓▒ ▒ ▒ ░░ ▒░ ░
  ░   ░ ░ ░ ▒  ░░░▒░ ░ ░  ░ ░  ░
░ ░   ░   ░ ░    ░░░ ░ ░    ░   
      ░     ░  ░   ░        ░  ░
                                 
`

type section struct {
	path    string
	content string
}

func main() {
	app := &cli.App{
		Name:  "glue",
		Usage: "Combine and extract project files into/from a single text file",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Include all files, do not use .gitignore",
			},
			&cli.StringSliceFlag{
				Name:    "ignore",
				Aliases: []string{"i"},
				Usage:   "Glob patterns to ignore",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "combined.txt",
				Usage:   "Output file name",
			},
			&cli.BoolFlag{
				Name:    "reverse",
				Aliases: []string{"r"},
				Usage:   "Reverse operation, recreate project from combined file",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(logo)
			start := time.Now()

			all := c.Bool("all")
			output := c.String("output")
			reverse := c.Bool("reverse")
			ignorePatterns := c.StringSlice("ignore")
			inclusionPatterns := c.Args().Slice()

			var numFiles int
			var err error
			if reverse {
				numFiles, err = reverseOperation(output)
				if err != nil {
					return err
				}
			} else {
				numFiles, err = combineFiles(all, output, inclusionPatterns, ignorePatterns)
				if err != nil {
					return err
				}
			}

			duration := time.Since(start)
			fmt.Printf("Processed %d files in %v\n", numFiles, duration)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func combineFiles(all bool, output string, inclusionPatterns, ignorePatterns []string) (int, error) {
	if len(inclusionPatterns) == 0 {
		inclusionPatterns = []string{"**/*"}
	}

	var exclusionPatterns []string
	if !all {
		gitignorePatterns, err := readGitignore()
		if err != nil && !os.IsNotExist(err) {
			return 0, err
		}
		exclusionPatterns = append(exclusionPatterns, gitignorePatterns...)
	}
	exclusionPatterns = append(exclusionPatterns, ignorePatterns...)

	var includedFiles []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Exclude the output file itself
		if path == output {
			return nil
		}

		// Check inclusion patterns
		matchedInclusion := false
		for _, pattern := range inclusionPatterns {
			matched, err := doublestar.PathMatch(pattern, path)
			if err != nil {
				return err
			}
			if matched {
				matchedInclusion = true
				break
			}
		}
		if !matchedInclusion {
			return nil
		}

		// Check exclusion patterns
		for _, pattern := range exclusionPatterns {
			matched, err := doublestar.PathMatch(pattern, path)
			if err != nil {
				return err
			}
			if matched {
				return nil // Exclude this file
			}
		}

		// Include this file
		includedFiles = append(includedFiles, path)
		return nil
	})
	if err != nil {
		return 0, err
	}

	outFile, err := os.Create(output)
	if err != nil {
		return 0, err
	}
	defer outFile.Close()

	for _, path := range includedFiles {
		fmt.Fprintf(outFile, "-- %s\n", path)
		fileType := getFileType(path)
		fmt.Fprintf(outFile, "```%s\n", fileType)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return 0, err
		}
		outFile.Write(content)
		fmt.Fprintf(outFile, "\n```\n\n")
	}

	return len(includedFiles), nil
}

func reverseOperation(output string) (int, error) {
	data, err := ioutil.ReadFile(output)
	if err != nil {
		return 0, err
	}

	content := string(data)
	sections := parseSections(content)

	for _, section := range sections {
		err := writeFile(section.path, section.content)
		if err != nil {
			return 0, err
		}
	}

	return len(sections), nil
}

func parseSections(content string) []section {
	var sections []section

	lines := strings.Split(content, "\n")
	var currentSection *section
	inCodeBlock := false
	var codeLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "-- ") {
			if currentSection != nil && inCodeBlock {
				// End of previous section
				currentSection.content = strings.Join(codeLines, "\n")
				sections = append(sections, *currentSection)
				codeLines = nil
				inCodeBlock = false
			}
			path := strings.TrimSpace(strings.TrimPrefix(line, "-- "))
			currentSection = &section{path: path}
		} else if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				// Start of code block
				inCodeBlock = true
			} else {
				// End of code block
				inCodeBlock = false
				if currentSection != nil {
					currentSection.content = strings.Join(codeLines, "\n")
					sections = append(sections, *currentSection)
					currentSection = nil
					codeLines = nil
				}
			}
		} else if inCodeBlock {
			codeLines = append(codeLines, line)
		}
	}
	// Handle any remaining section
	if currentSection != nil && len(codeLines) > 0 {
		currentSection.content = strings.Join(codeLines, "\n")
		sections = append(sections, *currentSection)
	}

	return sections
}

func writeFile(path string, content string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(content), 0644)
}

func readGitignore() ([]string, error) {
	data, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, nil
}

func getFileType(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		return ext[1:] // Remove the dot
	}
	return ""
}

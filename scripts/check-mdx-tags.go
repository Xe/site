package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TagInfo struct {
	Name      string
	Line      int
	IsClosing bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: check-mdx-tags <file_or_directory>")
		fmt.Println("Example: check-mdx-tags src/blog/")
		os.Exit(1)
	}

	path := os.Args[1]

	stat, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if stat.IsDir() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, ".mdx") || strings.HasSuffix(path, ".jsx") || strings.HasSuffix(path, ".tsx") {
				checkFile(path)
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error walking directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		checkFile(path)
	}
}

func checkFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	// Regex to match JSX tags
	// Matches: <TagName>, </TagName>, <TagName/>, <TagName ...>, </TagName ...>
	// More robust pattern to handle all self-closing variations
	tagRegex := regexp.MustCompile(`<(/?)([A-Z][a-zA-Z0-9]*)((?:[^>]*?)?)(/?)>`)

	var tagStack []TagInfo
	var hasErrors bool

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Find all tags in this line
		matches := tagRegex.FindAllStringSubmatch(line, -1)

		for _, match := range matches {
			isClosingSlash := match[1] == "/" // </TagName>
			tagName := match[2]
			innerContent := strings.TrimSpace(match[3])
			trailingSlash := match[4] == "/"

			// Check if it's self-closing: either ends with /> or has / at the end of content
			isSelfClosing := trailingSlash || strings.HasSuffix(innerContent, "/")

			if isSelfClosing {
				// Self-closing tags don't need to be tracked
				continue
			}

			if isClosingSlash {
				// This is a closing tag
				if len(tagStack) == 0 {
					fmt.Printf("%s:%d: Unexpected closing tag </%s> with no matching opening tag\n", filename, lineNum, tagName)
					hasErrors = true
					continue
				}

				// Check if it matches the most recent opening tag
				lastTag := tagStack[len(tagStack)-1]
				if lastTag.Name != tagName {
					fmt.Printf("%s:%d: Mismatched closing tag </%s>, expected </%s> (opened at line %d)\n",
						filename, lineNum, tagName, lastTag.Name, lastTag.Line)
					hasErrors = true
				} else {
					// Pop the matched tag from stack
					tagStack = tagStack[:len(tagStack)-1]
				}
			} else {
				// This is an opening tag
				tagStack = append(tagStack, TagInfo{
					Name:      tagName,
					Line:      lineNum,
					IsClosing: false,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s: %v\n", filename, err)
		return
	}

	// Check for unclosed tags
	for _, tag := range tagStack {
		fmt.Printf("%s:%d: Unclosed opening tag <%s>\n", filename, tag.Line, tag.Name)
		hasErrors = true
	}

	if !hasErrors {
		fmt.Printf("%s: ✓ All tags properly closed\n", filename)
	}
}

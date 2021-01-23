package blogger

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//In this module I want to parse frontmatter from markdown files, that means I will also have to parse markdown files first.
//We will be reading and parsing each file in the posts directory.

//return the absolute path to the posts directory
func postsDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(wd, "posts")
}

//retrieves all markdown files in the posts directory
func getPosts() ([]string, error) {
	pDir := postsDir()
	files, err := ioutil.ReadDir(pDir)
	if err != nil {
		return nil, fmt.Errorf("readposts: %v", err)
	}

	filenames := []string{}
	for _, file := range files {
		//check that file is a markdown file
		if !strings.HasSuffix(file.Name(), ".md") {
			return nil, fmt.Errorf("readposts: '.%s' unknown file extension", strings.Split(file.Name(), ".")[1])
		}

		filenames = append(filenames, file.Name())
	}

	return filenames, nil
}

//parses each post(markdown) file's contents
func readPosts() ([]string, []string, error) {
	//first we read the file
	files, err := getPosts()
	if err != nil {
		return nil, nil, err
	}

	frontMatter := []string{}
	mdContent := []string{}

	for _, file := range files {

		f, err := os.Open(filepath.Join(postsDir(), file))
		if err != nil {
			return nil, nil, fmt.Errorf("readposts: %v", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		var delimCount int

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "---" {
				delimCount++
				continue
			}

			if delimCount == 2 {
				if len(line) != 0 {
					mdContent = append(mdContent, line)
				}
			} else {
				if len(scanner.Text()) != 0 {
					frontMatter = append(frontMatter, line)
				}
			}
		}
	}

	return frontMatter, mdContent, nil
}

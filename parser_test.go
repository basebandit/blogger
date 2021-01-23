package blogger

import (
	"reflect"
	"testing"
)

func TestGetPosts(t *testing.T) {
	want := []string{
		"golang-arrays.md",
	}

	files, err := getPosts()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(files, want) {
		t.Fatalf("want %s got %s\n", want, files)
	}

}

func TestReadPosts(t *testing.T) {
	wantFrontMatter := []string{
		"title: \"Golang Arrays\"",
		"date: 2019-05-11T16:32:44+03:00",
		"draft: false",
		"tags: ['Arrays','golang']",
	}

	wantBody := []string{
		"## Go Arrays",
		"In Go, an array is a fixed length, ordered collection of values of the same type stored in contiguous memory locations.The number of elements is the array's length and it is never negative.",
	}

	fm, body, err := readPosts()
	if err != nil {
		t.Fatal(err)
	}

	if len(fm) != len(wantFrontMatter) {
		t.Fatalf("want %d lines of frontmatter got %d\n", len(wantFrontMatter), len(fm))
	}

	if !reflect.DeepEqual(fm, wantFrontMatter) {
		t.Fatalf("want %s got %s", wantFrontMatter, fm)
	}

	if !reflect.DeepEqual(body, wantBody) {
		t.Fatalf("want %v got %v", wantBody, body)
	}
}

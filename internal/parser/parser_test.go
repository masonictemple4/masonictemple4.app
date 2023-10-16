package parser

import (
	"os"
	"testing"

	"github.com/masonictemple4/masonictemple4.app/internal/dtos"
)

func TestStandaloneParser(t *testing.T) {

	var tests = []struct {
		name string
		fp   string
	}{
		{"Test parse with just frontmatter in the file", "/home/masonictemple4/personal/blogs/TestBlog.md"},
		{"Test with excess content after frontmatter", "/home/masonictemple4/personal/blogs/TestFullBlog.md"},
		// TODO: This test case could pass in the case that there is no additonal page content
		// after the frontmatter content and the parser reaches the end of file.
		// We'll probably want to check for both the begin and end precence otherwise just call
		// it invalid
		{"Test with valid open but no present close", "/home/masonictemple4/personal/blogs/TestWithNoFormatterClose.md"},
		{"Test with no formatter open, but with a valid close.. So technically is it an open???", "/home/masonictemple4/personal/blogs/TestNoFormatterOpen.md"},
		{"Test with no fromatter precense at all.", "/home/masonictemple4/personal/blogs/TestWithNoFrontMatterAtAll.md"},
		{"Test plain blog post inside of the frontmatter.", "/home/masonictemple4/personal/blogs/TestWithPlainContentInFrontmatter.md"},
	}

	// TODO: Come back and refactor this.
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result dtos.PostInput
			err := ParseFile(tt.fp, &result)
			if err != nil {
				t.Errorf("there was an error parsing the file: %v", err)
			}
			t.Logf("Test %d: %+v\n", i, result)
		})
	}

}

func TestParserObject(t *testing.T) {
	t.Run("Test just frontmatter in file with parser object", func(t *testing.T) {
		var result dtos.PostInput
		fp := "/home/masonictemple4/personal/blogs/TestBlog.md"
		f, err := os.Open(fp)
		if err != nil {
			t.Errorf("there was an error opening the file: %v", err)
		}
		defer f.Close()
		ymlFormat := NewYamlFrontMatterFormat()
		psr := New(f, ymlFormat)
		err = psr.Parse(&result)
	})

	t.Run("Test with the parser object", func(t *testing.T) {
		var result dtos.PostInput
		fp := "/home/masonictemple4/personal/blogs/TestFullBlog.md"
		f, err := os.Open(fp)
		if err != nil {
			t.Errorf("there was an error opening the file: %v", err)
		}
		defer f.Close()
		ymlFormat := NewYamlFrontMatterFormat()
		psr := New(f, ymlFormat)
		err = psr.Parse(&result)
	})

}

package utils

import (
	"fmt"
	"testing"
)

func TestGrabFileName(t *testing.T) {
	paths := []string{
		"/home/masonictemple4/personal/notes",
		"/home/masonictemple4/personal/dotfiles/README.md",
	}

	for _, p := range paths {
		res := GrabFileName(p)
		fmt.Printf("The result: %s\n", res)
	}
}

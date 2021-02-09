package main

import (
	"github.com/kobeHub/sendbetween/pkg/cmd/sendbetween"
)

func main() {
	err := sendbetween.NewCommand().Execute()
	CheckError(err)
}

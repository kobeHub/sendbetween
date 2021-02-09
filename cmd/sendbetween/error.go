package main

import (
	"context"
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		if err != context.Canceled {
			fmt.Fprintf(os.Stderr, "An error occured: %s\n", err)
		}
		os.Exit(3)
	}
}

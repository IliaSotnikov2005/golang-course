package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/IliaSotnikov2005/golang-course/task1/internal/client"
)

func main() {
	repoInput := flag.String("repo", "", "GitHub repository (full URL or 'owner/repo' format)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <repository>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://github.com/user/repo.git\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s user/repo\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -repo user/repo\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()

	var input string

	if *repoInput != "" {
		input = *repoInput
	} else {
		args := flag.Args()
		if len(args) == 0 {
			flag.Usage()
			fmt.Fprintln(os.Stderr, "Error: repository information is missing")
			os.Exit(1)
		}

		input = args[0]
	}

	client := client.NewClient()
	repo, err := client.GetRepositoryInfo(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(repo.String())
}

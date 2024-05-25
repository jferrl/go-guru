// Package main.
package main

import (
	"fmt"

	"github.com/google/go-github/v62/github"
	"github.com/gregjones/httpcache"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	action := githubactions.New()

	githubToken := action.GetInput("GITHUB_TOKEN")
	if githubToken == "" {
		githubactions.Fatalf("missing input 'GITHUB_TOKEN'")
	}

	repository := action.GetInput("GITHUB_REPOSITORY")
	if repository == "" {
		githubactions.Fatalf("missing input 'GITHUB_REPOSITORY'")
	}

	pullRequest := action.GetInput("GITHUB_REF")
	if pullRequest == "" {
		githubactions.Fatalf("missing input 'GITHUB_REF'")
	}

	action.AddMask(githubToken)

	fmt.Println("Hello World From GitHub Action")

	_ = github.NewClient(
		httpcache.NewMemoryCacheTransport().Client(),
	).WithAuthToken(githubToken)

}

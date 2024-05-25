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

	repository := action.GetInput("GITHUB_REPOSITORY_NAME")
	if repository == "" {
		githubactions.Fatalf("missing input 'GITHUB_REPOSITORY_NAME'")
	}

	owner := action.GetInput("GITHUB_REPOSITORY_OWNER")
	if owner == "" {
		githubactions.Fatalf("missing input 'GITHUB_REPOSITORY_OWNER'")
	}

	action.AddMask(githubToken)

	fmt.Println("Hello World From GitHub Action")

	client := github.NewClient(
		httpcache.NewMemoryCacheTransport().Client(),
	).WithAuthToken(githubToken)

}

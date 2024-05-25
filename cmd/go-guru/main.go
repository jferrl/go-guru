// Package main.
package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
	"github.com/gregjones/httpcache"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	ctx := context.Background()

	action := githubactions.New()

	_, err := action.Context()
	if err != nil {
		action.Fatalf("failed to get action context: %v", err)
	}

	githubToken := action.GetInput("GITHUB_TOKEN")
	if githubToken == "" {
		action.Fatalf("missing input 'GITHUB_TOKEN'")
	}

	action.AddMask(githubToken)

	fmt.Println("Hello World From GitHub Action")

	githubClient := github.NewClient(
		httpcache.NewMemoryCacheTransport().Client(),
	).WithAuthToken(githubToken)

	orgs, _, err := githubClient.Organizations.List(ctx, "jferrl", nil)
	if err != nil {
		action.Fatalf("failed to get PR: %v", err)
	}

	action.Infof("PR: %s", orgs)
}

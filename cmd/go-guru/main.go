// Package main.
package main

import (
	"context"

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
		action.Fatalf("missing input GITHUB_TOKEN")
	}

	action.AddMask(githubToken)

	githubClient := github.NewClient(
		httpcache.NewMemoryCacheTransport().Client(),
	).WithAuthToken(githubToken)

	user, _, err := githubClient.Users.Get(ctx, "jferrl")
	if err != nil {
		action.Fatalf("failed to get PR: %v", err)
	}

	pr, _, err := githubClient.PullRequests.ListFiles(ctx, "google", "go-github", 1, nil)
	if err != nil {
		action.Fatalf("failed to get PR: %v", err)
	}

	action.Infof("PR: %v", pr)

	action.Infof("User: %s", user)
}

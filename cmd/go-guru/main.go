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

	actionCtx, err := action.Context()
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

	action.Infof("Context: %v", actionCtx)

	action.Infof("User: %s", user)

}

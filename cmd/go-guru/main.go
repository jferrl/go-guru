// Package main.
package main

import (
	"context"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/gregjones/httpcache"
	"github.com/sashabaranov/go-openai"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	ctx := context.Background()

	action := githubactions.New()

	_, err := action.Context()
	if err != nil {
		action.Fatalf("failed to get action context: %v", err)
	}

	// githubToken := action.GetInput("GITHUB_TOKEN")
	// if githubToken == "" {
	// 	action.Fatalf("missing input GITHUB_TOKEN")
	// }

	githubToken := ""
	openaiToken := ""

	action.AddMask(githubToken)

	githubClient := github.NewClient(
		httpcache.NewMemoryCacheTransport().Client(),
	).WithAuthToken(githubToken)

	openaiConf := openai.DefaultConfig(openaiToken)
	openaiConf.OrgID = "org-ZAWXMY7d7r4KRVpLWWHE0GYf"

	openaiClient := openai.NewClientWithConfig(openaiConf)

	assistant, err := openaiClient.RetrieveAssistant(ctx, "")
	if err != nil {
		action.Fatalf("failed to get assistant: %v", err)
	}

	if err != nil {
		action.Fatalf("failed to create thread and run: %v", err)
	}

	action.Infof("Assistant: %s", assistant.ID)

	files, _, err := githubClient.PullRequests.ListFiles(ctx, "google", "go-github", 3174, nil)
	if err != nil {
		action.Fatalf("failed to get PR: %v", err)
	}

	createdRun, err := openaiClient.CreateThreadAndRun(ctx, openai.CreateThreadAndRunRequest{
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    openai.ThreadMessageRoleUser,
					Content: files[0].GetPatch(),
				},
			},
		},
		RunRequest: openai.RunRequest{
			AssistantID: assistant.ID,
		},
	})
	if err != nil {
		action.Fatalf("failed to create thread and run: %v", err)
	}

	action.Infof("Run: %s", createdRun.ID)

	// wait until the run is completed

	for {
		run, err := openaiClient.RetrieveRun(ctx, createdRun.ThreadID, createdRun.ID)
		if err != nil {
			action.Fatalf("failed to get run: %v", err)
		}

		if run.Status == openai.RunStatusCompleted {
			action.Infof("Run completed: %s", run.ID)
			break
		}

		action.Infof("Run status: %s", run.Status)

		time.Sleep(2 * time.Second)
	}

	steps, err := openaiClient.ListRunSteps(ctx, createdRun.ThreadID, createdRun.ID, openai.Pagination{})
	if err != nil {
		action.Fatalf("failed to get run: %v", err)
	}

	for _, step := range steps.RunSteps {
		action.Infof("Step: %s", step.ID)

		if step.Status == openai.RunStepStatusCompleted {
			action.Infof("Step completed: %s", step.ID)

			message, err := openaiClient.RetrieveMessage(ctx, createdRun.ThreadID, step.StepDetails.MessageCreation.MessageID)
			if err != nil {
				action.Fatalf("failed to get message: %v", err)
			}

			for _, content := range message.Content {
				if content.Text != nil {
					action.Infof("Message: %s", content.Text.Value)
				}
			}
		}
	}

}

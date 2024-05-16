package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

// Change represents a file change in a pull request
type Change struct {
	Filename string `json:"filename"`
	Patch    string `json:"patch"`
}

// ReviewRequest represents a request to the GPT API
type ReviewRequest struct {
	Changes []Change `json:"changes"`
}

// Comment represents a comment to be posted on a pull request
type Comment struct {
	Text     string `json:"text"`
	CommitID string `json:"commit_id"`
	Path     string `json:"path"`
	Position int    `json:"position"`
}

// ReviewResponse represents a response from the GPT API
type ReviewResponse struct {
	Comments []Comment `json:"comments"`
}

func main() {
	// Load environment variables
	githubToken := os.Getenv("GITHUB_TOKEN")
	gptAPIKey := os.Getenv("GPT_API_KEY")
	repo := os.Getenv("GITHUB_REPOSITORY")
	prNumber := os.Getenv("GITHUB_REF")

	// Split the repository owner and name
	repoParts := strings.Split(repo, "/")
	if len(repoParts) != 2 {
		fmt.Println("Invalid repository format")
		return
	}
	repoOwner := repoParts[0]
	repoName := repoParts[1]

	// Set up GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get pull request details
	prNum := strings.Split(prNumber, "/")[2]

	prID, err := strconv.Atoi(prNum)
	if err != nil {
		fmt.Printf("Error converting PR number to int: %s\n", err)
		return
	}

	_, _, err = client.PullRequests.Get(ctx, repoOwner, repoName, prID)
	if err != nil {
		fmt.Printf("Error fetching PR: %s\n", err)
		return
	}

	// Get pull request files
	files, _, err := client.PullRequests.ListFiles(ctx, repoOwner, repoName, prID, nil)
	if err != nil {
		fmt.Printf("Error fetching PR files: %s\n", err)
		return
	}

	// Prepare changes for GPT API
	var changes []Change
	for _, file := range files {
		if file.GetStatus() != "removed" {
			changes = append(changes, Change{
				Filename: file.GetFilename(),
				Patch:    file.GetPatch(),
			})
		}
	}

	// Call GPT API
	reviewRequest := ReviewRequest{Changes: changes}
	reqBody, err := json.Marshal(reviewRequest)
	if err != nil {
		fmt.Printf("Error marshalling request: %s\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.your-gpt-service.com/review", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+gptAPIKey)
	req.Header.Set("Content-Type", "application/json")

	gptClient := &http.Client{}

	resp, err := gptClient.Do(req)
	if err != nil {
		fmt.Printf("Error calling GPT API: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err)
		return
	}

	var reviewResponse ReviewResponse
	err = json.Unmarshal(body, &reviewResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %s\n", err)
		return
	}

	// Post comments on the pull request
	for _, comment := range reviewResponse.Comments {
		_, _, err := client.PullRequests.CreateComment(ctx, repoOwner, repoName, prID, &github.PullRequestComment{
			Body:     github.String(comment.Text),
			CommitID: github.String(comment.CommitID),
			Path:     github.String(comment.Path),
			Position: github.Int(comment.Position),
		})
		if err != nil {
			fmt.Printf("Error creating comment: %s\n", err)
		}
	}
}

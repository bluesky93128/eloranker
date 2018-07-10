package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
)

type requestMessageDataImportVariant struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type importTypeGithubIssues struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

func (c *Client) importVariant(message requestMessageDataImportVariant) {
	canWrite, err := c.hasWriteAccess("")
	if err != nil {
		c.Error(err.Error(), "variant:import")
		return
	}
	if !canWrite {
		c.Error("not enough permissions to import variants", "variant:import")
		return
	}

	if message.Type != "github-issues" {
		c.Error("Unsupported import type", "variant:import")
		return
	}

	var data importTypeGithubIssues
	json.Unmarshal(message.Data, &data)

	githubRequestOptions :=  &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	client := github.NewClient(nil)
	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Issues.ListByRepo(context.Background(), data.Owner, data.Repo, githubRequestOptions)
		if err != nil {
			c.Error(err.Error(), "variant:import")
			return
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		githubRequestOptions.Page = resp.NextPage
	}

	for _, issue := range allIssues {
		if issue.IsPullRequest() {
			continue
		}

		title := issue.GetTitle()
		id := issue.GetNumber()

		variantText := fmt.Sprintf("#%d: %s", id, title)
		UUID, err := c.room.AllocateNewVariant(c)
		if err != nil {
			c.Error(err.Error(), "variant:import")
			return
		}

		err = c.room.UpdateVariant(nil, requestMessageDataUpdateVariant{
			UUID: UUID,
			Text: variantText,
		})
		if err != nil {
			c.Error(err.Error(), "variant:import")
		}
	}
}

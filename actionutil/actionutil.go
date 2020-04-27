package actionutil

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/posener/goaction"
	"github.com/posener/script"
	"golang.org/x/oauth2"
)

func GitConfig(name, email string) error {
	err := script.Exec("git", "config", "user.name", name).ToStdout()
	if err != nil {
		return err
	}
	return script.Exec("git", "config", "user.email", email).ToStdout()
}

// PRComment adds a comment to the curerrent pull request. If the comment already exists
// it updates the exiting comment with the new content.
func PRComment(ctx context.Context, token string, actionID string, content string) error {
	var (
		own = goaction.Owner()
		prj = goaction.Project()
		num = goaction.PrNum()

		hiddenSignature = fmt.Sprintf("<!-- comment by %s -->", actionID)
	)

	oauthClient := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	gh := github.NewClient(oauthClient)

	reviewID := int64(-1)
	comments, _, err := gh.PullRequests.ListComments(ctx, own, prj, num, nil)
	if err != nil {
		return err
	}
	for _, c := range comments {
		if strings.HasPrefix(c.GetBody(), hiddenSignature) {
			reviewID = c.GetID()
			break
		}
	}

	commentBody := hiddenSignature + "\n\n" + content

	if reviewID >= 0 {
		log.Printf("Updating existing review: %d\n", reviewID)
		_, _, err = gh.PullRequests.UpdateReview(ctx, own, prj, num, reviewID, commentBody)
	} else {
		log.Printf("Creating new review")
		_, _, err = gh.PullRequests.CreateReview(ctx, own, prj, num, &github.PullRequestReviewRequest{
			Body:  github.String(commentBody),
			Event: github.String("COMMENT"),
		})
	}
	return err
}

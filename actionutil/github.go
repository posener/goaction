package actionutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v31/github"
	"github.com/posener/goaction"
	"github.com/posener/goaction/log"
	"golang.org/x/oauth2"
)

// PRComment adds a comment to the curerrent pull request. If the comment already exists
// it updates the exiting comment with the new content.
func PRComment(ctx context.Context, token string, content string) error {
	var (
		own = goaction.Owner()
		prj = goaction.Project()
		num = goaction.PrNum()

		// Hidden signature is added to the review comment body and is used in following runs to
		// identify which comment to update.
		hiddenSignature = fmt.Sprintf(
			"<!-- comment by %s (%s) -->",
			goaction.Workflow, goaction.ActionID)
	)

	oauthClient := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	gh := github.NewClient(oauthClient)

	// Look for an existing review. reviewID<0 means that we didn't find a matching review.
	reviewID := int64(-1)
	reviews, _, err := gh.PullRequests.ListReviews(ctx, own, prj, num, nil)
	if err != nil {
		return err
	}
	for _, review := range reviews {
		if strings.HasPrefix(review.GetBody(), hiddenSignature) {
			reviewID = review.GetID()
			break
		}
	}

	// Update or post a new review.
	commentBody := hiddenSignature + "\n\n" + content
	if reviewID >= 0 {
		log.Printf("Updating existing review: %d\n", reviewID)
		_, _, err = gh.PullRequests.UpdateReview(ctx, own, prj, num, reviewID, commentBody)
	} else {
		log.Printf("Creating new review")
		_, _, err = gh.PullRequests.CreateReview(ctx, own, prj, num,
			&github.PullRequestReviewRequest{
				Body:  github.String(commentBody),
				Event: github.String("COMMENT"),
			})
	}
	return err
}

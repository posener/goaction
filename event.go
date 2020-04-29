package goaction

// Code auto generated with `go run ./internal/genevents/main.go`. DO NOT EDIT

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/v31/github"
)

//go:generate go run ./internal/genevents/main.go

// A Github action triggering event.
// See https://help.github.com/en/actions/reference/events-that-trigger-workflows.
type EventType string

// All Github action event types.
const (
	EventCheckRun                 EventType = "check_run"
	EventCheckSuite               EventType = "check_suite"
	EventCreate                   EventType = "create"
	EventDelete                   EventType = "delete"
	EventDeployment               EventType = "deployment"
	EventFork                     EventType = "fork"
	EventGollum                   EventType = "gollum"
	EventIssueComment             EventType = "issue_comment"
	EventIssues                   EventType = "issues"
	EventLabel                    EventType = "label"
	EventMilestone                EventType = "milestone"
	EventPageBuild                EventType = "page_build"
	EventProject                  EventType = "project"
	EventProjectCard              EventType = "project_card"
	EventPublic                   EventType = "public"
	EventPullRequest              EventType = "pull_request"
	EventPullRequestReview        EventType = "pull_request_review"
	EventPullRequestReviewComment EventType = "pull_request_review_comment"
	EventPush                     EventType = "push"
	EventRegistryPackage          EventType = "registry_package"
	EventRelease                  EventType = "release"
	EventStatus                   EventType = "status"
	EventWatch                    EventType = "watch"
	EventSchedule                 EventType = "schedule"
	EventRepositoryDispatch       EventType = "repository_dispatch"
)

// GetCheckRun returns information about a current check run.
func GetCheckRun() (*github.CheckRunEvent, error) {
	if Event != EventCheckRun {
		return nil, fmt.Errorf("not 'check_run' event")
	}
	var i github.CheckRunEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetCheckSuite returns information about a current check suite.
func GetCheckSuite() (*github.CheckSuiteEvent, error) {
	if Event != EventCheckSuite {
		return nil, fmt.Errorf("not 'check_suite' event")
	}
	var i github.CheckSuiteEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetCreate returns information about a current create.
func GetCreate() (*github.CreateEvent, error) {
	if Event != EventCreate {
		return nil, fmt.Errorf("not 'create' event")
	}
	var i github.CreateEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetDelete returns information about a current delete.
func GetDelete() (*github.DeleteEvent, error) {
	if Event != EventDelete {
		return nil, fmt.Errorf("not 'delete' event")
	}
	var i github.DeleteEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetDeployment returns information about a current deployment.
func GetDeployment() (*github.DeploymentEvent, error) {
	if Event != EventDeployment {
		return nil, fmt.Errorf("not 'deployment' event")
	}
	var i github.DeploymentEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetFork returns information about a current fork.
func GetFork() (*github.ForkEvent, error) {
	if Event != EventFork {
		return nil, fmt.Errorf("not 'fork' event")
	}
	var i github.ForkEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetGollum returns information about a current gollum.
func GetGollum() (*github.GollumEvent, error) {
	if Event != EventGollum {
		return nil, fmt.Errorf("not 'gollum' event")
	}
	var i github.GollumEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetIssueComment returns information about a current issue comment.
func GetIssueComment() (*github.IssueCommentEvent, error) {
	if Event != EventIssueComment {
		return nil, fmt.Errorf("not 'issue_comment' event")
	}
	var i github.IssueCommentEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetIssues returns information about a current issues.
func GetIssues() (*github.IssuesEvent, error) {
	if Event != EventIssues {
		return nil, fmt.Errorf("not 'issues' event")
	}
	var i github.IssuesEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetLabel returns information about a current label.
func GetLabel() (*github.LabelEvent, error) {
	if Event != EventLabel {
		return nil, fmt.Errorf("not 'label' event")
	}
	var i github.LabelEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetMilestone returns information about a current milestone.
func GetMilestone() (*github.MilestoneEvent, error) {
	if Event != EventMilestone {
		return nil, fmt.Errorf("not 'milestone' event")
	}
	var i github.MilestoneEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetPageBuild returns information about a current page build.
func GetPageBuild() (*github.PageBuildEvent, error) {
	if Event != EventPageBuild {
		return nil, fmt.Errorf("not 'page_build' event")
	}
	var i github.PageBuildEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetProject returns information about a current project.
func GetProject() (*github.ProjectEvent, error) {
	if Event != EventProject {
		return nil, fmt.Errorf("not 'project' event")
	}
	var i github.ProjectEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetProjectCard returns information about a current project card.
func GetProjectCard() (*github.ProjectCardEvent, error) {
	if Event != EventProjectCard {
		return nil, fmt.Errorf("not 'project_card' event")
	}
	var i github.ProjectCardEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetPublic returns information about a current public.
func GetPublic() (*github.PublicEvent, error) {
	if Event != EventPublic {
		return nil, fmt.Errorf("not 'public' event")
	}
	var i github.PublicEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetPullRequest returns information about a current pull request.
func GetPullRequest() (*github.PullRequestEvent, error) {
	if Event != EventPullRequest {
		return nil, fmt.Errorf("not 'pull_request' event")
	}
	var i github.PullRequestEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetPullRequestReview returns information about a current pull request review.
func GetPullRequestReview() (*github.PullRequestReviewEvent, error) {
	if Event != EventPullRequestReview {
		return nil, fmt.Errorf("not 'pull_request_review' event")
	}
	var i github.PullRequestReviewEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetPullRequestReviewComment returns information about a current pull request review comment.
func GetPullRequestReviewComment() (*github.PullRequestReviewCommentEvent, error) {
	if Event != EventPullRequestReviewComment {
		return nil, fmt.Errorf("not 'pull_request_review_comment' event")
	}
	var i github.PullRequestReviewCommentEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetPush returns information about a current push.
func GetPush() (*github.PushEvent, error) {
	if Event != EventPush {
		return nil, fmt.Errorf("not 'push' event")
	}
	var i github.PushEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetRelease returns information about a current release.
func GetRelease() (*github.ReleaseEvent, error) {
	if Event != EventRelease {
		return nil, fmt.Errorf("not 'release' event")
	}
	var i github.ReleaseEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetStatus returns information about a current status.
func GetStatus() (*github.StatusEvent, error) {
	if Event != EventStatus {
		return nil, fmt.Errorf("not 'status' event")
	}
	var i github.StatusEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetWatch returns information about a current watch.
func GetWatch() (*github.WatchEvent, error) {
	if Event != EventWatch {
		return nil, fmt.Errorf("not 'watch' event")
	}
	var i github.WatchEvent
	err := decodeEventInfo(&i)
	return &i, err
}

// GetRepositoryDispatch returns information about a current repository dispatch.
func GetRepositoryDispatch() (*github.RepositoryDispatchEvent, error) {
	if Event != EventRepositoryDispatch {
		return nil, fmt.Errorf("not 'repository_dispatch' event")
	}
	var i github.RepositoryDispatchEvent
	err := decodeEventInfo(&i)
	return &i, err
}

func decodeEventInfo(i interface{}) error {
	f, err := os.Open(eventPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(i)
}

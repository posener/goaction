// Creates action files for Go code
package main

import (
	"context"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v31/github"
	"github.com/posener/goaction"
	"github.com/posener/goaction/actionutil"
	"github.com/posener/goaction/metadata"
	"github.com/posener/script"
	"golang.org/x/oauth2"
)

var (
	path  = flag.String("path", "", "Path to main package.")
	name  = flag.String("name", "", "Override action name, the default name is the package name.")
	desc  = flag.String("desc", "", "Override action description, the default description is the package synopsis.")
	icon  = flag.String("icon", "", "Set branding icon.")
	color = flag.String("color", "", "Set branding color.")

	email       = os.Getenv("EMAIL")
	githubToken = os.Getenv("GITHUB_TOKEN") // Optional
)

const (
	goMain     = "main.go"
	action     = "action.yml"
	dockerfile = "Dockerfile"
)

func init() {
	if email == "" {
		email = "posener@gmail.com"
	}

	log.SetFlags(log.Lshortfile)
}

func main() {
	flag.Parse()

	// Load go code.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join(*path, goMain), nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Parse Go code to Github actions metadata.
	m, err := metadata.New(f)
	if err != nil {
		log.Fatal(err)
	}
	if *name != "" {
		m.Name = *name
	}
	if *desc != "" {
		m.Desc = *desc
	}

	m.Branding.Icon = *icon
	m.Branding.Color = *color

	// Applying changes.

	// Create action file.
	err = script.Writer("yml", func(w io.Writer) error { return yaml.NewEncoder(w).Encode(m) }).
		ToFile(action)
	if err != nil {
		log.Fatal(err)
	}

	// Create dockerfile
	err = script.Writer("template", func(w io.Writer) error { return tmpls.Execute(w, struct{ Path string }{Path: *path}) }).
		ToFile(dockerfile)
	if err != nil {
		log.Fatal(err)
	}

	if !goaction.CI {
		log.Println("Skipping commit stage.")
		os.Exit(0)
	}

	// Runs only in Github CI mode.

	err = actionutil.GitConfig("goaction", email)
	if err != nil {
		log.Fatal(err)
	}

	// Add files to git
	err = script.Exec("git", "add", action, dockerfile).ToStdout()
	if err != nil {
		log.Fatal(err)
	}

	diff := stagedDiff()

	if diff == "" {
		log.Println("No changes were made. Aborting")
		os.Exit(0)
	}

	switch {
	case goaction.IsPush():
		push()
	case goaction.IsPR():
		pr(diff)
	default:
		log.Fatalf("unexpected action mode.")
	}
}

func stagedDiff() string {
	var diff strings.Builder
	for _, path := range []string{action, dockerfile} {
		d, err := script.Exec("git", "diff", "--staged", "--no-color", path).Head(-5).ToString()
		if err != nil {
			log.Fatalf("git diff for %s: %s", path, err)
		}
		if d != "" {
			diff.WriteString(fmt.Sprintf("Path: %s\n\n", path))
			diff.WriteString(fmt.Sprintf("```diff\n%s\n```\n\n", d))
		}
	}
	return diff.String()
}

// Commit and push chnages to upstream branch.
func push() {
	err := script.Exec("git", "commit", "-m", "Update readme according to Go doc").ToStdout()
	if err != nil {
		log.Fatal(err)
	}
	err = script.Exec("git", "push", "origin", "HEAD:"+goaction.Branch()).ToStdout()
	if err != nil {
		log.Fatal(err)
	}
}

const commentHeader = "[GoAction](https://github.com/posener/goaction) diff:"

func pr(diff string) {
	if githubToken == "" {
		log.Println("In order to add request comment, set the GITHUB_TOKEN input.")
		return
	}

	var (
		own = goaction.Owner()
		prj = goaction.Project()
		num = goaction.PrNum()
	)

	ctx := context.Background()
	oauthClient := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken}))

	gh := github.NewClient(oauthClient)

	exitingReviewID := int64(-1)
	comments, _, err := gh.PullRequests.ListComments(ctx, own, prj, num, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range comments {
		if strings.HasPrefix(c.GetBody(), commentHeader) {
			exitingReviewID = c.GetID()
			break
		}
	}

	commentBody := commentHeader + "\n\n" + diff

	if exitingReviewID > 0 {
		log.Printf("Updating existing review: %d\n", exitingReviewID)
		_, _, err = gh.PullRequests.UpdateReview(ctx, own, prj, num, exitingReviewID, commentBody)
	} else {
		log.Printf("Creating new review")
		_, _, err = gh.PullRequests.CreateReview(ctx, own, prj, num, &github.PullRequestReviewRequest{
			Body:  github.String(commentBody),
			Event: github.String("COMMENT"),
		})
	}
	if err != nil {
		log.Fatal(err)
	}
}

var tmpls = template.Must(template.New("dockerfile").Parse(`
FROM golang:1.14.1-alpine3.11
RUN apk add git

ADD . /home/src/
WORKDIR /home/src
RUN go build {{ .Path }} -o /bin/action

FROM alpine:3.11
RUN apk add git
COPY --from=0 /bin/action /bin/action

ENTRYPOINT [ "/bin/action" ]
`))

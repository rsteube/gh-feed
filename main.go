package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2/pkg/api"
)

type event struct {
	Type  string
	Actor struct {
		Login string
		Url   string
	}
	Repo struct {
		Name string
		Url  string
	}
	Payload struct {
		Action string
		Forkee struct {
			FullName string `json:"full_name"`
			Url      string
		}
		Issue struct {
			Number int
			Title  string
		}
		PullRequest struct {
			Number int
			Title  string
		} `json:"pull_request"`
		Comment struct {
			Body string
		}
		Ref     string
		RefType string `json:"ref_type"`
		Release struct {
			TagName string `json:"tag_name"`
		}
	}
	CreatedAt time.Time
}

func main() {
	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	login := struct{ Login string }{}
	err = client.Get("user", &login)
	if err != nil {
		fmt.Println(err)
		return
	}

	var events []event
	err = client.Get(fmt.Sprintf("users/%v/received_events", login.Login), &events)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO can be different events based on Payload.Action
	for _, e := range events {
		switch e.Type {
		case "CreateEvent":
			fmt.Printf("%v %v %v %v on %v\n", e.Actor.Login, format("created"), e.Payload.RefType, faint(e.Payload.Ref), e.Repo.Name)
		case "DeleteEvent":
			fmt.Printf("%v %v %v %v from %v\n", e.Actor.Login, format("deleted"), e.Payload.RefType, faint(e.Payload.Ref), e.Repo.Name)
		case "ForkEvent":
			fmt.Printf("%v %v %v from %v\n", e.Actor.Login, format("forked"), e.Payload.Forkee.FullName, e.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf("%v %v issue %v \"%v\" on %v\n", e.Actor.Login, format("commented"), e.Payload.Issue.Number, faint(e.Payload.Issue.Title), e.Repo.Name)
			// fmt.Printf("%v\n", e.Payload.Comment.Body) // TODO limit and format markdown
		case "PullRequestEvent":
			fmt.Printf("%v %v pull request %v \"%v\" on %v\n", e.Actor.Login, format(e.Payload.Action), e.Payload.PullRequest.Number, faint(e.Payload.PullRequest.Title), e.Repo.Name)
		case "PushEvent":
			fmt.Printf("%v %v to %v\n", e.Actor.Login, format("pushed"), e.Repo.Name)
		case "ReleaseEvent":
			fmt.Printf("%v %v %v of %v\n", e.Actor.Login, format("released"), e.Payload.Release.TagName, e.Repo.Name)
		case "WatchEvent":
			fmt.Printf("%v %v %v\n", e.Actor.Login, format("starred"), e.Repo.Name) // TODO might not be a starred event
		default:
			fmt.Printf("unknown event %#v\n", e.Type)
		}
	}
}

func faint(s string) string {
	return lipgloss.NewStyle().Faint(true).Render(s)
}

func format(s string) string {
	color := map[string]string{
		"closed":    "1",
		"commented": "9",
		"created":   "4",
		"deleted":   "1",
		"forked":    "5",
		"opened":    "2",
		"pushed":    "6",
		"released":  "4",
		"starred":   "3",
	}[s]
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(s)
}

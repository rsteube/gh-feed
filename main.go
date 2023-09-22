package main

import (
	"fmt"
	"strings"
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
		Member struct {
			Login string
		}
		PullRequest struct {
			Number int
			Title  string
		} `json:"pull_request"`
		Repo struct {
			Name string
		}
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

func (e event) FormatActor() string {
	if strings.HasSuffix(e.Actor.Login, "[bot]") {
		return lipgloss.NewStyle().Faint(true).Render(e.Actor.Login)
	}
	return e.Actor.Login
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
			switch e.Payload.RefType {
			case "repository":
				fmt.Printf("%v %v %v %v\n", e.FormatActor(), format("created"), e.Payload.RefType, e.Repo.Name)
			case "branch":
				fmt.Printf("%v %v %v %v on %v\n", e.FormatActor(), format("created"), e.Payload.RefType, faint(e.Payload.Ref), e.Repo.Name)
			default:
				fmt.Printf("unknown reftype %#v\n", e.Payload.RefType)
			}
		case "DeleteEvent":
			fmt.Printf("%v %v %v %v from %v\n", e.FormatActor(), format("deleted"), e.Payload.RefType, faint(e.Payload.Ref), e.Repo.Name)
		case "ForkEvent":
			fmt.Printf("%v %v %v from %v\n", e.FormatActor(), format("forked"), e.Payload.Forkee.FullName, e.Repo.Name)
		case "IssuesEvent":
			fmt.Printf("%v %v issue %v \"%v\" on %v\n", e.FormatActor(), format(e.Payload.Action), e.Payload.Issue.Number, faint(e.Payload.Issue.Title), e.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf("%v %v issue %v \"%v\" on %v\n", e.FormatActor(), format("commented"), e.Payload.Issue.Number, faint(e.Payload.Issue.Title), e.Repo.Name)
			// fmt.Printf("%v\n", e.Payload.Comment.Body) // TODO limit and format markdown
		case "MemberEvent":
			fmt.Printf("%v %v %v to %v\n", e.FormatActor(), format(e.Payload.Action), e.Payload.Member.Login, e.Repo.Name)
		case "PullRequestEvent":
			fmt.Printf("%v %v pull request %v \"%v\" on %v\n", e.FormatActor(), format(e.Payload.Action), e.Payload.PullRequest.Number, faint(e.Payload.PullRequest.Title), e.Repo.Name)
		case "PushEvent":
			fmt.Printf("%v %v to %v\n", e.FormatActor(), format("pushed"), e.Repo.Name)
		case "ReleaseEvent":
			fmt.Printf("%v %v %v of %v\n", e.FormatActor(), format("released"), e.Payload.Release.TagName, e.Repo.Name)
		case "WatchEvent":
			fmt.Printf("%v %v %v\n", e.FormatActor(), format("starred"), e.Repo.Name) // TODO might not be a starred event
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
		"added":     "4",
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

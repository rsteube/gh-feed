package main

import (
	"fmt"
	"time"

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
			fmt.Printf("%v created %v %v on %v\n", e.Actor.Login, e.Payload.RefType, e.Payload.Ref, e.Repo.Name)
		case "DeleteEvent":
			fmt.Printf("%v deleted %v %v from %v\n", e.Actor.Login, e.Payload.RefType, e.Payload.Ref, e.Repo.Name)
		case "ForkEvent":
			fmt.Printf("%v forked %v from %v\n", e.Actor.Login, e.Payload.Forkee.FullName, e.Repo.Name)
		case "WatchEvent":
			fmt.Printf("%v starred %v\n", e.Actor.Login, e.Repo.Name)
		case "PushEvent":
			fmt.Printf("%v pushed to %v\n", e.Actor.Login, e.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf("%v commented issue %v \"%v\" on %v\n", e.Actor.Login, e.Payload.Issue.Number, e.Payload.Issue.Title, e.Repo.Name)
			// fmt.Printf("%v\n", e.Payload.Comment.Body) // TODO limit and format markdown
		case "PullRequestEvent":
			fmt.Printf("%v %v pull request %v \"%v\" on %v\n", e.Actor.Login, e.Payload.Action, e.Payload.PullRequest.Number, e.Payload.PullRequest.Title, e.Repo.Name)
		default:
			fmt.Printf("unknown event %#v\n", e.Type)
		}
	}
}

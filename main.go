package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v45/github"
	"gopkg.in/yaml.v2"
)

type config struct {
	User       string `yaml:"user"`
	OAuthToken string `yaml:"oauth_token"`
}

type response struct {
	Permission string `json:"permission"`
	RoleName   string `json:"role_name"`
	User       struct {
		AvatarURL        string `json:"avatar_url"`
		EventsURL        string `json:"events_url"`
		FollowersURL     string `json:"followers_url"`
		FollowingURL     string `json:"following_url"`
		GistsURL         string `json:"gists_url"`
		GravatarID       string `json:"gravatar_id"`
		HTMLURL          string `json:"html_url"`
		ID               int64  `json:"id"`
		Login            string `json:"login"`
		NodeID           string `json:"node_id"`
		OrganizationsURL string `json:"organizations_url"`
		Permissions      struct {
			Admin    bool `json:"admin"`
			Maintain bool `json:"maintain"`
			Pull     bool `json:"pull"`
			Push     bool `json:"push"`
			Triage   bool `json:"triage"`
		} `json:"permissions"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		RoleName          string `json:"role_name"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"user"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [user or organization] [repository]\n", os.Args[0])
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadFile(filepath.Join(u.HomeDir, ".config", "hub"))

	var cfg map[string][]config

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	gcs, ok := cfg["github.com"]
	if !ok || len(gcs) == 0 {
		log.Fatal("GitHub token not found")
	}
	gc := gcs[0]

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gc.OAuthToken},
	)
	client := github.NewClient(oauth2.NewClient(context.Background(), ts))

	resp, _, err := client.Repositories.GetPermissionLevel(
		context.Background(),
		flag.Arg(0),
		flag.Arg(1),
		gc.User)
	if err != nil {
		fmt.Printf("You don't have permission for %s/%s: %v\n", flag.Arg(0), flag.Arg(1), err)
	} else {
		fmt.Printf("You are %s on %s/%s\n", *resp.Permission, flag.Arg(0), flag.Arg(1))
	}
}

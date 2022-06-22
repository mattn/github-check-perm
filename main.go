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
	if err != nil {
		log.Fatal(err)
	}

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
		if err2, ok := err.(*github.ErrorResponse); ok && err2.Response.StatusCode != 404 {
			fmt.Printf("You don't have permission for %s/%s\n", flag.Arg(0), flag.Arg(1))
		} else {
			log.Fatal(err2.Message)
		}
	} else {
		fmt.Printf("You are %s on %s/%s\n", *resp.Permission, flag.Arg(0), flag.Arg(1))
	}
}

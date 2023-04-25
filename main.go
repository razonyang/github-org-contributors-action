package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"

	"github.com/google/go-github/v52/github"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

var org = ""
var output = ""

func init() {
	flag.StringVar(&org, "org", "", "organization name")
	flag.StringVar(&output, "output", "", "output filename")
}

func main() {
	flag.Parse()
	if org == "" || output == "" {
		log.Fatal("org and output are required")
	}

	ctx := context.Background()
	contributors, err := getContributors(ctx, org)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]*contributor)
	for _, c := range contributors {
		if *c.Type != "User" {
			continue
		}

		v, ok := data[*c.Login]
		if !ok {
			data[*c.Login] = newContribution(c)
		} else {
			v.Contributions += *c.Contributions
			v.Repos++
		}
	}

	var bs []byte
	switch path.Ext(output) {
	case ".yaml", ".yml":
		bs, err = yaml.Marshal(data)
	case ".toml":
		bs, err = toml.Marshal(data)
	default:
		bs, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		log.Fatal(err)
	}

	if err = os.MkdirAll(path.Dir(output), 0755); err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(output, bs, 0644); err != nil {
		log.Fatal(err)
	}
}

type contributor struct {
	ID            int64  `json:"id" toml:"id" yaml:"id"`
	Login         string `json:"login" toml:"login" yaml:"login"`
	AvatarURL     string `json:"avatar_url" toml:"avatar_url" yaml:"avatar_url"`
	Contributions int    `json:"contributions" toml:"contributions" yaml:"contributions"`
	Repos         int    `json:"repos" toml:"repos" yaml:"repos"`
}

func newContribution(c *github.Contributor) *contributor {
	return &contributor{
		ID:            *c.ID,
		Login:         *c.Login,
		AvatarURL:     *c.AvatarURL,
		Contributions: *c.Contributions,
		Repos:         1,
	}
}

func newClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func getContributors(ctx context.Context, owner string) (contributors []*github.Contributor, err error) {
	repos, err := getRepos(ctx, owner)
	if err != nil {
		return
	}

	for _, repo := range repos {
		page := 1
		for {
			log.Printf("fetching contributors of %q, page: %d\n", *repo.Name, page)
			v, resp, err := getContributorsByRepo(ctx, owner, *repo.Name, page)
			if err != nil {
				return nil, err
			}
			contributors = append(contributors, v...)
			if resp.NextPage > 0 {
				page = resp.NextPage
			} else {
				break
			}
		}
	}

	return
}

func getContributorsByRepo(ctx context.Context, owner, repo string, page int) (contributors []*github.Contributor, resp *github.Response, err error) {
	client := newClient(ctx)
	contributors, resp, err = client.Repositories.ListContributors(ctx, owner, repo, &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    page,
		},
	})
	if err != nil {
		return
	}

	return
}

func getRepos(ctx context.Context, owner string) (repos []*github.Repository, err error) {
	page := 1
	for {
		v, resp, err := getReposByPage(ctx, owner, page)
		log.Printf("fetching repos of %q, page: %d\n", owner, page)
		if err != nil {
			return nil, err
		}

		repos = append(repos, v...)
		if resp.NextPage > 0 {
			page = resp.NextPage
		} else {
			break
		}
	}

	return
}

func getReposByPage(ctx context.Context, owner string, page int) (repos []*github.Repository, resp *github.Response, err error) {
	client := newClient(ctx)
	results, resp, err := client.Repositories.ListByOrg(ctx, owner, &github.RepositoryListByOrgOptions{
		Type: "public",
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    page,
		},
	})
	if err != nil {
		return
	}
	for _, repo := range results {
		if !*repo.Fork {
			repos = append(repos, repo)
		}
	}

	return
}

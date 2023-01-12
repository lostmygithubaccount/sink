/*
TODO: refactor further
*/
package utils

// imports
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cli/go-gh/pkg/api"
	"golang.org/x/exp/slices"
)

// types -- must match GitHub web API types (currently only REST)
type User struct {
	Login string `json:"login"`
	URL   string `json:"html_url"`
}

type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	URL    string `json:"html_url"`
}

type Label struct {
	Name        string `json:"name"`
	Repo        string `json:"repo"`
	Color       string `json:"color"`
	Description string `json:"description"`
	NewName     string `json:"new_name"`
}

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type Labels map[string]Label // map of labels by name
type Repos []Repo            // slice of repos

func (user User) String() string {
	return user.URL
}

func (issue Issue) String() string {
	return issue.URL
}

func (label Label) String() string {
	return fmt.Sprintf("%*s | %*s | %s", 20, label.Name, 7, label.Color, label.Description)
}

func (labels Labels) String() string {
	out := "\n"
	for _, label := range labels {
		out += fmt.Sprintf("\t%s\n", label)
	}

	return out
}

func (repo Repo) String() string {
	return repo.Name
}

func (repos Repos) String() string {
	out := "\n"
	for _, repo := range repos {
		out += fmt.Sprintf("\t%s\n", repo)
	}

	return out
}

// public helpers
func GetUser(c api.RESTClient) (user User) {

	// construct REST endpoint
	endpoint := fmt.Sprintf("user")

	// make the request
	err := c.Get(endpoint, &user)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func GetTargetRepos(org, team string, targetRepos, excludeTeamRepos, extraRepos Repos, c api.RESTClient) (targets Repos) {

	if len(targetRepos) == 0 {
		log.Println("target-repos is empty! using team repositories...")
		if team == "" {
			log.Fatal("one of target-repos or team is required")
		} else {
			log.Println("team:", team)

			// get the team repositories
			teamRepos := getTeamRepos(org, team, c)
			log.Println("team repositories:", teamRepos)

			for _, repo := range teamRepos {
				if !slices.Contains(excludeTeamRepos, repo) {
					targets = append(targets, repo)
				}
			}
		}
	} else {
		targets = targetRepos
	}

	// extra repos
	for _, repo := range extraRepos {
		if !slices.Contains(targets, repo) && !slices.Contains(excludeTeamRepos, repo) {
			targets = append(targets, repo)
		}
	}

	return
}

func GetIssue(org string, repo Repo, issueNumber int, c api.RESTClient) (issue Issue) {

	// construct REST endpoint
	endpoint := fmt.Sprintf("repos/%s/%s/issues/%d", org, repo, issueNumber)

	// get the issue
	err := c.Get(endpoint, &issue)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func GetRepoLabels(org string, repo Repo, c api.RESTClient) (labels Labels) {

	// TODO: remove hardcoded limit, add paging
	// construct REST endpoint
	endpoint := fmt.Sprintf("repos/%s/%s/labels?per_page=100", org, repo)

	// intermediary response struct
	response := new([]Label)
	err := c.Get(endpoint, &response)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: this feels weird
	labels = make(Labels)

	// extract the label names into map
	for _, label := range *response {
		labels[label.Name] = label
	}
	log.Printf("found %d labels in %s/%s", len(labels), org, repo)

	return
}

func CopyIssue(org string, target Repo, sourceIssue Issue, c api.RESTClient) (response Issue) {

	// construct REST endpoint
	endpoint := fmt.Sprintf("repos/%s/%s/issues", org, target)

	// extra body to indicate bot
	extra := fmt.Sprintf("\n\n---\n\ncopied from %s", sourceIssue.URL)

	// construct the target issue
	targetIssue := Issue{
		Title: sourceIssue.Title,
		Body:  fmt.Sprintf("%s%s", sourceIssue.Body, extra),
	}

	// marshal the issue struct into JSON
	json, err := json.Marshal(targetIssue)
	if err != nil {
		log.Fatal(err)
	}

	// create the issue
	err = c.Post(endpoint, bytes.NewBuffer(json), &response)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func SyncLabels(dryRun bool, org string, source, target Repo, sourceLabels, targetLabels Labels, c api.RESTClient) {

	// for each label in the source repo
	for _, sourceLabel := range sourceLabels {

		// check if the label exists in the target repo
		targetLabel, exists := targetLabels[sourceLabel.Name]
		log.Printf("\t\tlabel %s exists in target repo: %t", sourceLabel.Name, exists)
		if exists {
			// check if the label matchest in the target repo
			if !labelsMatch(sourceLabel, targetLabel) {
				// update the label in the target repo
				log.Printf("\t\tupdating label %s in %s/%s", sourceLabel.Name, org, target)
				if !dryRun {
					updateLabel(org, target, sourceLabel, c)
				}
			} else {
				log.Printf("\t\tlabel %s in %s/%s is up to date", sourceLabel.Name, org, target)
			}
		} else {
			// create the label in the target repo
			log.Printf("\t\tcreating label %s in %s/%s", sourceLabel.Name, org, target)
			if !dryRun {
				createLabel(org, target, sourceLabel, c)
			}
		}
	}
}

func ReposToRepos(repos []string) (r Repos) {
	for _, repo := range repos {
		r = append(r, Repo{Name: repo})
	}
	return
}

// private helpers
func getTeamRepos(org, team string, c api.RESTClient) (repos Repos) {

	// construct REST endpoint
	endpoint := fmt.Sprintf("orgs/%s/teams/%s/repos", org, team)

	// intermediary response struct
	response := new([]struct{ Name string })
	err := c.Get(endpoint, &response)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func labelsMatch(source, target Label) bool {

	// check for matching color and description
	if source.Color == target.Color && source.Description == target.Description {
		return true
	}
	return false
}

func createLabel(org string, repo Repo, sourceLabel Label, c api.RESTClient) {

	// construct REST endpoint
	endpoint := fmt.Sprintf("repos/%s/%s/labels", org, repo)

	// construct the target label
	targetLabel := Label{
		Name:        sourceLabel.Name,
		Color:       sourceLabel.Color,
		Description: sourceLabel.Description,
	}

	// marshal the issue struct into JSON
	json, err := json.Marshal(targetLabel)
	if err != nil {
		log.Fatal(err)
	}

	// create the issue
	response := new(Label)
	err = c.Post(endpoint, bytes.NewBuffer(json), &response)
	if err != nil {
		log.Fatal(err)
	}
}

func updateLabel(org string, repo Repo, sourceLabel Label, c api.RESTClient) {

	// construct REST endpoint
	endpoint := fmt.Sprintf("repos/%s/%s/labels/%s", org, repo, sourceLabel.Name)

	// construct the target label
	targetLabel := Label{
		Name:        sourceLabel.Name,
		Color:       sourceLabel.Color,
		Description: sourceLabel.Description,
	}

	// marshal the issue struct into JSON
	json, err := json.Marshal(targetLabel)
	if err != nil {
		log.Fatal(err)
	}

	// update the issue
	response := new(Label)
	err = c.Patch(endpoint, bytes.NewBuffer(json), &response)
	if err != nil {
		log.Fatal(err)
	}
}

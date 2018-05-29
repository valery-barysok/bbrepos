package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Link struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type Value struct {
	Key   string            `json:"key"`
	Slug  string            `json:"slug"`
	Links map[string][]Link `json:"links"`
}

type Items struct {
	Values []Value `json:"values"`
}

const getProjects = "https://git.junolab.net/rest/api/1.0/projects?limit=1000"
const getRepositories = "https://git.junolab.net/rest/api/1.0/projects/%s/repos?limit=1000"

func get(url string, token string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return http.DefaultClient.Do(req)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERRROR: bitbucket token missed")
		fmt.Printf("Usage: %v <token>\n", os.Args[0])
		os.Exit(1)
	}

	token := os.Args[1]
	res, err := get(getProjects, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	var projects Items

	json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		return
	}

	var prKeys []string
	for _, val := range projects.Values {
		prKeys = append(prKeys, val.Key)
		printRepos(getRepos(val.Key, token))
	}
}

func printRepos(repos []string) {
	for _, repo := range repos {
		fmt.Println(repo)
	}
}

func getRepos(prKey string, token string) []string {
	res, err := get(fmt.Sprintf(getRepositories, prKey), token)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var repositories Items

	json.NewDecoder(res.Body).Decode(&repositories)
	if err != nil {
		return nil
	}

	var repos []string
	for _, val := range repositories.Values {
		for _, v := range val.Links["clone"] {
			if v.Name == "ssh" {
				repos = append(repos, fmt.Sprintf("%s,/%s,%s", v.Href, prKey, val.Slug))
				break
			}
		}
	}
	return repos
}

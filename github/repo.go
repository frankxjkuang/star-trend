package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Repository basic info
// source：https://api.github.com/repos/{Repository.FullName}
type Repository struct {
	FullName         string    `json:"full_name"`
	StargazersCount  int       `json:"stargazers_count"` // star count
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Forks            int       `json:"forks"`
	Watchers         int       `json:"watchers"`
	SubscribersCount int       `json:"subscribers_count"`
	CreateAtFormat   string    `json:"create_at_format"`
}

// GetRepoDetail get repo basic info
func (gh *Github) GetRepoDetail(fullName string) (*Repository, error) {
	repo := new(Repository)
	if fullName == "" {
		return repo, nil
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s", fullName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return repo, err
	}
	req.Header.Set("Accept", "application/vnd.Github.v3+json")
	if gh.HasToken() {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", gh.GetToken()))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return repo, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return repo, errHandle(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(repo)
	if err != nil {
		return repo, err
	}

	repo.CreateAtFormat = repo.CreatedAt.Format("2006-01-02 15:04:05")
	return repo, nil
}

func (repo *Repository) totalPage(pageSize int) int {
	return repo.StargazersCount/pageSize + 1
}

func (repo *Repository) lastPage(pageSize int) int {
	return repo.totalPage(pageSize)
}

func errHandle(resp *http.Response) error {
	var result struct {
		Message          string `json:"message"`
		DocumentationUrl string `json:"documentation_url"`
	}
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}
	return errors.New(result.Message)
}

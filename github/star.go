package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	maxPage = 400
)

var (
	errPageSizeLTEZero = errors.New("page size low than equal 0")
	errTooManyStars    = errors.New("repo has too many stargazers, Github won't allow us to list all stars")
)

// Stargazers single msg info
type Stargazers struct {
	StarredAt time.Time `json:"starred_at"`
}

// GetStargazers get star info
func (gh *Github) GetStargazers(repo *Repository) (stargazers []Stargazers, err error) {
	if gh.GetPerPage() <= 0 {
		return stargazers, errPageSizeLTEZero
	}
	if repo.totalPage(gh.GetPerPage()) > maxPage {
		return stargazers, errTooManyStars
	}
	var g errgroup.Group
	g.SetLimit(4)
	var lock sync.Mutex
	for page := 1; page <= repo.totalPage(gh.GetPerPage()); page++ {
		page := page
		g.Go(
			func() error {
				result, err := gh.getStargazersPage(repo, page)
				if err != nil {
					return err
				}
				if len(result) == 0 {
					return nil
				}
				lock.Lock()
				defer lock.Unlock()
				stargazers = append(stargazers, result...)
				return nil
			},
		)
	}
	err = g.Wait()
	sort.Slice(stargazers, func(i, j int) bool {
		return stargazers[i].StarredAt.Before(stargazers[j].StarredAt)
	})
	return
}

func (gh *Github) getStargazersPage(repo *Repository, page int) (stargazers []Stargazers, err error) {
	stargazers = make([]Stargazers, 0)
	url := fmt.Sprintf("https://api.github.com/repos/%s/stargazers?page=%d&per_page=%d", repo.FullName, page, gh.GetPerPage())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return stargazers, err
	}
	req.Header.Set("Accept", "application/vnd.Github.v3.star+json")
	if gh.HasToken() {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", gh.GetToken()))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return stargazers, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return stargazers, errHandle(resp)
	}

	return stargazers, json.NewDecoder(resp.Body).Decode(&stargazers)
}

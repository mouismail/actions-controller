package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

func getRepos(client *github.Client, owner string, limiter *rate.Limiter, wg *sync.WaitGroup, reposChan chan<- []*github.Repository) {
	defer wg.Done()

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		// Wait for the rate limiter
		if err := limiter.Wait(context.Background()); err != nil {
			log.Printf("Rate limiter error: %v", err)
			return
		}

		repos, resp, err := client.Repositories.List(context.Background(), owner, opt)
		if err != nil {
			log.Printf("Error fetching repositories: %v", err)
			return
		}

		reposChan <- repos

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}

func checkAppInstallation(client *github.Client, owner, repo, appSlug string, limiter *rate.Limiter, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()

	// Wait for the rate limiter
	if err := limiter.Wait(context.Background()); err != nil {
		log.Printf("Rate limiter error: %v", err)
		return
	}

	installations, _, err := client.Apps.ListInstallations(context.Background(), nil)
	if err != nil {
		log.Printf("Error fetching app installations: %v", err)
		return
	}
// TODO: o(n'2) to be refactored to be o(log n) max
	for _, installation := range installations {
		if installation.AppSlug != nil && *installation.AppSlug == appSlug {
			// Check if the app is installed on the specific repository
			repos, _, err := client.Apps.ListRepos(context.Background(), *installation.ID, nil)
			if err != nil {
				log.Printf("Error fetching repositories for installation %d: %v", *installation.ID, err)
				return
			}

			for _, repo := range repos.Repositories {
				if repo.GetFullName() == fmt.Sprintf("%s/%s", owner, repo) {
					resultChan <- fmt.Sprintf("App %s is installed on repo %s/%s", appSlug, owner, repo)
					return
				}
			}
		}
	}

	resultChan <- fmt.Sprintf("App %s is not installed on repo %s/%s", appSlug, owner, repo)
}

func main() {
	owner := "octocat"         // Replace with the desired GitHub username
	appSlug := "your-app-slug" // Replace with the desired GitHub app slug

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Create a rate limiter that allows 1 request per second with a burst of 5 requests
	limiter := rate.NewLimiter(1, 5)

	var wg sync.WaitGroup
	reposChan := make(chan []*github.Repository, 10)
	resultChan := make(chan string, 10)

	// Start a Go routine to fetch repositories
	wg.Add(1)
	go getRepos(client, owner, limiter, &wg, reposChan)

	// Close the reposChan when all repositories are fetched
	go func() {
		wg.Wait()
		close(reposChan)
	}()

	// Start Go routines to check app installations for each repository
	for repos := range reposChan {
		for _, repo := range repos {
			wg.Add(1)
			go checkAppInstallation(client, owner, repo.GetName(), appSlug, limiter, &wg, resultChan)
		}
	}

	// Close the resultChan when all checks are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Print the results
	for result := range resultChan {
		fmt.Println(result)
	}
}

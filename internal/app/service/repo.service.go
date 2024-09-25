package service

import (
	"fmt"
	"sync"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
)

func SyncOrgRepos(orgName string, sem chan struct{}, repoCount chan int) error {

	repos, err := repository.FetchRepoOfOrgFromGitea(1, orgName)
	if err != nil {
		return err
	}

	err = repository.SyncReposWithDB(repos)
	if err != nil {
		return err
	}
	repoCount <- len(repos)
	fmt.Printf("Synced %d repositories of %s\n", len(repos), orgName)
	return nil
}
func SyncAllRepos() (int, error) {
	var orgs []*model.Org
	orgs, err := repository.GetAllOrgs()
	if err != nil {
		return 0, err
	}

	err = repository.ClearRepos()
	if err != nil {
		return 0, err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, goRoutines)
	errorsChan := make(chan error, len(orgs))
	repoCount := make(chan int, len(orgs))
	totalRepos := 0

	for _, org := range orgs {
		wg.Add(1)

		go func(org *model.Org) {
			defer wg.Done()
			sem <- struct{}{}

			if err := SyncOrgRepos(org.Username, sem, repoCount); err != nil {
				errorsChan <- err
			}

			<-sem
		}(org)
	}

	wg.Wait()

	close(repoCount)
	close(errorsChan)

	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	for count := range repoCount {
		fmt.Println("Testing Channels: ", count)
		totalRepos += count
	}

	fmt.Println("Total Repos: ", totalRepos)

	return totalRepos, err
}

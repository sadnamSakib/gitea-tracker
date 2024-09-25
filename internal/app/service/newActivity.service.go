package service

import (
	"fmt"
	"sync"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
)

func SyncNewUserActivities(username string, sem chan struct{}) error {

	user, err := repository.GetUser(username)

	if err != nil {
		return err
	}
	lastUpdateTime := user.Last_updated
	activities, err := repository.FetchNewUserActivityFromGitea(1, username, lastUpdateTime)
	if err != nil {
		return err
	}
	err = repository.SyncNewActivitiesWithDB(username, activities)
	if err != nil {
		return err
	}

	return nil
}

func SyncDailyActivities() (int, error) {

	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return 0, err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, goRoutines)
	errorsChan := make(chan error, len(users))
	usersSynced := 0
	mu := sync.Mutex{}

	for _, user := range users {
		wg.Add(1)

		go func(user model.User) {
			defer wg.Done()
			sem <- struct{}{}

			if err := SyncNewUserActivities(user.Username, sem); err != nil {
				errorsChan <- err
			} else {

				mu.Lock()
				usersSynced++
				mu.Unlock()
			}

			<-sem
		}(user)
	}

	wg.Wait()
	close(errorsChan)
	close(sem)

	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}

	fmt.Println("Synced ", usersSynced, " users")
	return usersSynced, err
}

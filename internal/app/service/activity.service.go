package service

import (
	"fmt"
	"sync"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
)

func SyncUserActivities(userName string) error {

	activities, err := repository.FetchUserActivityFromGitea(1, userName)
	if err != nil {
		return err
	}
	fmt.Println("Activities fetched for", userName, ":", len(activities))
	err = repository.SyncActivitiesWithDB(userName, activities)
	if err != nil {
		return err
	}
	fmt.Printf("User %s synced with %d activities\n", userName, len(activities))
	return nil
}

func SyncAllActivities() (int, error) {

	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return 0, err
	}

	err = repository.ClearActivities()
	if err != nil {
		return 0, err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, goRoutines)
	errorsChan := make(chan error, len(users))
	userSynced := 0
	mu := sync.Mutex{}

	for _, user := range users {

		wg.Add(1)
		go func(user model.User) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if err := SyncUserActivities(user.Username); err != nil {
				errorsChan <- err
			} else {

				mu.Lock()
				userSynced++
				mu.Unlock()
			}
		}(user)
	}

	wg.Wait()
	close(errorsChan)
	close(sem)

	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}

	fmt.Println("Synced", userSynced, "users")
	return userSynced, err
}

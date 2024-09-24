package controller

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func SyncAllUsers() (int, error) {
	users, err := repository.FetchUsersFromGitea(1)
	if err != nil {
		return 0, err
	}

	err = repository.SyncUsersWithDB(users)
	if err != nil {
		return 0, err
	}

	return len(users), nil
}

func SyncUsers(c echo.Context) error {
	usersSynced, err := SyncAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Users": usersSynced})
}

func SyncUserHeatmap(username string, wg *sync.WaitGroup, sem chan struct{}) error {
	defer wg.Done()
	defer func() { <-sem }()
	heatmap, err := repository.FetchUserHeatmapActivityFromGitea(username)
	if err != nil {
		fmt.Println(username, err.Error())
		return err
	}
	err = repository.SyncHeatMaps(heatmap)
	if err != nil {
		fmt.Println(username, err.Error())
		return err
	}
	fmt.Printf("User %s's Heatmap synced \n", username)
	return nil
}

func SyncAllHeatmaps() (int, error) {
	users, err := repository.GetAllUsers("", "")
	if err != nil {
		return 0, err
	}
	err = repository.ClearHeatmaps()
	if err != nil {
		return 0, err
	}
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	errorsChan := make(chan error, len(users))
	userSynced := 0
	for _, user := range users {
		wg.Add(1)
		sem <- struct{}{}

		if err := SyncUserHeatmap(user.Username, &wg, sem); err != nil {
			errorsChan <- err
		} else {
			userSynced++
		}

	}
	wg.Wait()
	close(errorsChan)
	close(sem)
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}
	fmt.Printf("Synced %d users\n", userSynced)
	return len(users), err
}

func SyncHeatMaps(c echo.Context) error {
	users, err := SyncAllHeatmaps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"User's Heatmaps": users})
}

func SyncUserActivities(userName string) error {

	activities, err := repository.FetchUserActivityFromGitea(1, userName)
	if err != nil {
		return err
	}
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

	// Clear activities before starting the sync
	err = repository.ClearActivities()
	if err != nil {
		return 0, err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)             // Semaphore to limit concurrency
	errorsChan := make(chan error, len(users)) // Buffered error channel
	userSynced := 0
	mu := sync.Mutex{} // Mutex to safely update shared counter

	for _, user := range users {

		wg.Add(1)
		go func(user model.User) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// Sync user activities and capture errors
			if err := SyncUserActivities(user.Username); err != nil {
				errorsChan <- err
			} else {
				// Safely increment userSynced
				mu.Lock()
				userSynced++
				mu.Unlock()
			}
		}(user)
	}

	wg.Wait()         // Wait for all goroutines to finish
	close(errorsChan) // Close the error channel after all processing is complete
	close(sem)        // Close the semaphore channel

	// Handle errors after syncing
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}

	fmt.Println("Synced", userSynced, "users")
	return userSynced, err
}

func SyncActivities(c echo.Context) error {
	users, err := SyncAllActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Users Activities": users})
}

func SyncAllOrganizations() (int, error) {
	orgs, err := repository.FetchOrgsFromGitea(1)
	if err != nil {
		return 0, err
	}
	err = repository.ClearOrgs()
	if err != nil {
		return 0, err
	}
	err = repository.SyncOrgsWithDB(orgs)
	if err != nil {
		fmt.Println("Error here")
		return 0, err
	}
	fmt.Println(orgs, len(orgs))
	return len(orgs), nil
}

func SyncOrganizations(c echo.Context) error {
	orgsSynced, err := SyncAllOrganizations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Organizations": orgsSynced})
}

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
	sem := make(chan struct{}, 10) // Semaphore to limit concurrency
	errorsChan := make(chan error, len(orgs))
	repoCount := make(chan int, len(orgs)) // Buffered channel to avoid blocking on count
	totalRepos := 0

	for _, org := range orgs {
		wg.Add(1)
		// Goroutine to sync repos for each org concurrently
		go func(org *model.Org) {
			defer wg.Done()   // Mark the goroutine as done when finished
			sem <- struct{}{} // Acquire a spot in the semaphore

			// SyncOrgRepos handles repo syncing, and any error is returned to errorsChan
			if err := SyncOrgRepos(org.Username, sem, repoCount); err != nil {
				errorsChan <- err
			}

			<-sem // Release the semaphore spot
		}(org)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close channels after all processing is complete
	close(repoCount) // Close repoCount after all repos are synced
	close(errorsChan)

	// Check for any errors that occurred during syncing
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

func SyncRepos(c echo.Context) error {
	orgsSynced, err := SyncAllRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"Repos For Organization": orgsSynced})
}

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
	sem := make(chan struct{}, 10)             // Semaphore to limit concurrency
	errorsChan := make(chan error, len(users)) // Buffered to avoid blocking
	usersSynced := 0
	mu := sync.Mutex{} // Mutex to safely update shared counter

	for _, user := range users {
		wg.Add(1)

		// Run each sync operation in a separate goroutine
		go func(user model.User) {
			defer wg.Done()   // Ensure we call Done even if there's an error
			sem <- struct{}{} // Acquire a spot in the semaphore

			// SyncNewUserActivities and handle errors
			if err := SyncNewUserActivities(user.Username, sem); err != nil {
				errorsChan <- err
			} else {
				// Safely increment the usersSynced counter
				mu.Lock()
				usersSynced++
				mu.Unlock()
			}

			<-sem // Release the semaphore spot
		}(user)
	}

	wg.Wait()         // Wait for all goroutines to finish
	close(errorsChan) // Close the error channel after all work is done
	close(sem)        // Close the semaphore when done

	// Process any errors that were reported
	for e := range errorsChan {
		err = e
		fmt.Println(err)
	}

	fmt.Println("Synced ", usersSynced, " users")
	return usersSynced, err
}

func SyncNewActivity(c echo.Context) error {
	users, err := SyncDailyActivities()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"User's New Activities": users})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func DailySync(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	defer func() {
		err = repository.UpdateSyncStatus(true)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Error updating sync status: "+err.Error()))
		}
	}()

	// Send initial status to the client
	ws.WriteMessage(websocket.TextMessage, []byte("Sync process started"))
	startTime := time.Now()
	err = repository.UpdateSyncStatus(false)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error updating sync status: "+err.Error()))

	}

	orgs, err := SyncAllOrganizations()

	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing organizations: "+err.Error()))

	}
	fmt.Println("Organizations Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Organizations Synced: %d", orgs)))

	users, err := SyncAllUsers()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing users: "+err.Error()))

	}
	fmt.Println("Users Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Users Synced: %d", users)))

	repos, err := SyncAllRepos()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing repos: "+err.Error()))

	}
	fmt.Println("Repos Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Repositories Synced: %d", repos)))

	usersNewActivities, err := SyncDailyActivities()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing daily activities: "+err.Error()))

	}
	elapsedTime := int64(time.Since(startTime).Seconds())
	fmt.Println("Daily Activities Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("New Activities Synced For Users: %d", usersNewActivities)))
	err = repository.SyncSystemSummary(orgs, repos, users, time.Now().Local(), elapsedTime)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing system summary: "+err.Error()))
	}
	fmt.Println("System Summary Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Sync completed in %d seconds", elapsedTime)))
	return nil

}

package controller

import (
	"fmt"
	"net/http"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/service"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

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

	orgs, err := service.SyncAllOrganizations()

	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing organizations: "+err.Error()))

	}
	fmt.Println("Organizations Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Organizations Synced: %d", orgs)))

	users, err := service.SyncAllUsers()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing users: "+err.Error()))

	}
	fmt.Println("Users Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Users Synced: %d", users)))

	repos, err := service.SyncAllRepos()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing repos: "+err.Error()))

	}
	fmt.Println("Repos Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Repositories Synced: %d", repos)))

	usersNewActivities, err := service.SyncDailyActivities()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing daily activities: "+err.Error()))

	}
	fmt.Println("Daily Activities Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("New Activities Synced For Users: %d", usersNewActivities)))
	repoActivitiesSynced, err := service.SyncAllNewRepoActivities()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing repo activities: "+err.Error()))
	}
	fmt.Println("Repo Activities Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Repo Activities Synced: %d", repoActivitiesSynced)))
	elapsedTime := int64(time.Since(startTime).Seconds())
	err = repository.SyncSystemSummary(orgs, repos, users, time.Now().Local(), elapsedTime)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing system summary: "+err.Error()))
	}
	fmt.Println("System Summary Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Sync completed in %d seconds", elapsedTime)))
	return nil

}

func TotalSync(c echo.Context) error {
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

	orgs, err := service.SyncAllOrganizations()

	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing organizations: "+err.Error()))

	}
	fmt.Println("Organizations Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Organizations Synced: %d", orgs)))

	users, err := service.SyncAllUsers()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing users: "+err.Error()))

	}
	fmt.Println("Users Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Users Synced: %d", users)))

	repos, err := service.SyncAllRepos()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing repos: "+err.Error()))

	}
	fmt.Println("Repos Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Repositories Synced: %d", repos)))

	usersNewActivities, err := service.SyncAllActivities()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing daily activities: "+err.Error()))

	}

	fmt.Println("Daily Activities Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("New Activities Synced For Users: %d", usersNewActivities)))
	repoActivitiesSynced, err := service.SyncAllRepoActivities()
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing repo activities: "+err.Error()))
	}
	fmt.Println("Repo Activities Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Repo Activities Synced: %d", repoActivitiesSynced)))
	elapsedTime := int64(time.Since(startTime).Seconds())
	err = repository.SyncSystemSummary(orgs, repos, users, time.Now().Local(), elapsedTime)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error syncing system summary: "+err.Error()))
	}
	fmt.Println("System Summary Synchronised")
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Sync completed in %d seconds", elapsedTime)))
	return nil

}

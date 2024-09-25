package service

import (
	"fmt"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"github.com/robfig/cron/v3"
)

func SyncDailyData() {
	startTime := time.Now()
	err := repository.UpdateSyncStatus(false)
	if err != nil {
		_ = fmt.Errorf("error in updating sync status: %v", err)
		return
	}

	orgs, err := controller.SyncAllOrganizations()
	if err != nil {
		_ = fmt.Errorf("error in syncing organizations: %v", err)
		return
	}
	fmt.Println("Organizations Synchronised")

	users, err := controller.SyncAllUsers()
	if err != nil {
		_ = fmt.Errorf("error in syncing Users: %v", err)
		return
	}
	fmt.Println("Users Synchronised")

	repos, err := controller.SyncAllRepos()
	if err != nil {
		_ = fmt.Errorf("error in syncing repos: %v", err)
		return
	}
	fmt.Println("Repos Synchronised")

	_, err = controller.SyncDailyActivities()
	if err != nil {
		_ = fmt.Errorf("error in syncing repos: %v", err)
		return
	}
	elapsedTime := int64(time.Since(startTime).Seconds())
	err = repository.SyncSystemSummary(orgs, repos, users, time.Now().Local(), elapsedTime)
	if err != nil {
		_ = fmt.Errorf("error in syncing system summary: %v", err)
		return
	}
	fmt.Println("Daily Activities Synchronised")

}

func InitCronScheduler() *cron.Cron {
	c := cron.New(cron.WithLocation(time.Local))
	c.AddFunc("5 0 * * *", SyncDailyData)
	c.Start()
	fmt.Println("Cron Scheduler Initialized")
	return c
}

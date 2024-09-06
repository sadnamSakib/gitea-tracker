package service

import (
	"fmt"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
	"github.com/robfig/cron/v3"
)

func SyncDailyData() {
	_, err := controller.SyncAllOrganizations()
	if err != nil {
		_ = fmt.Errorf("error in syncing organizations: %v", err)
		return
	}
	fmt.Println("Organizations Synchronised")

	_, err = controller.SyncAllUsers()
	if err != nil {
		_ = fmt.Errorf("error in syncing Users: %v", err)
		return
	}
	fmt.Println("Users Synchronised")

	_, err = controller.SyncAllRepos()
	if err != nil {
		_ = fmt.Errorf("error in syncing repos: %v", err)
		return
	}
	fmt.Println("Repos Synchronised")
	_, err = controller.SyncAllHeatmaps()
	if err != nil {
		_ = fmt.Errorf("error in syncing repos: %v", err)
		return
	}
	_, err = controller.SyncDailyActivities()
	if err != nil {
		_ = fmt.Errorf("error in syncing repos: %v", err)
		return
	}
	fmt.Println("Daily Activities Synchronised")

}

func InitCronScheduler() *cron.Cron {
	c := cron.New()
	c.AddFunc("30 23 * * *", SyncDailyData)
	c.Start()
	fmt.Println("Cron Scheduler Initialized")
	return c
}

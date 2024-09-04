package service

import (
	"fmt"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/app/controller"
)

func SyncDailyData() error {
	err := controller.SyncAllOrganizations()
	if err != nil {
		return err
	}
	fmt.Println("Organizations Synchronised")

	err = controller.SyncAllUsers()
	if err != nil {
		return err
	}
	fmt.Println("Users Synchronised")

	err = controller.SyncAllRepos()
	if err != nil {
		return err
	}
	fmt.Println("Repos Synchronised")
	err = controller.SyncAllHeatmaps()
	if err != nil {
		return err
	}
	err = controller.SyncDailyActivities()
	if err != nil {
		return err
	}
	fmt.Println("Daily Activities Synchronised")
	return nil

}

func Scheduler(hour, minute, second int, job func() error) {
	for {

		now := time.Now()

		nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())

		if nextRun.Before(now) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		sleepDuration := nextRun.Sub(now)
		fmt.Printf("Sleeping for %v until next run at %v\n", sleepDuration, nextRun.Format("15:04:05"))

		time.Sleep(sleepDuration)

		err := job()
		if err != nil {
			fmt.Println(err)
		}
	}
}

package service

import (
	"fmt"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"github.com/robfig/cron/v3"
)

func TotalSync() {
	startTime := time.Now()

	err := repository.UpdateSyncStatus(false)
	if err != nil {
		fmt.Printf("Error updating sync status: %v\n", err)
		return
	}

	orgs, err := SyncAllOrganizations()
	if err != nil {
		fmt.Printf("Error syncing organizations: %v\n", err)
		return
	}
	fmt.Println("Organizations Synchronised")

	users, err := SyncAllUsers()
	if err != nil {
		fmt.Printf("Error syncing users: %v\n", err)
		return
	}
	fmt.Println("Users Synchronised")

	repos, err := SyncAllRepos()
	if err != nil {
		fmt.Printf("Error syncing repos: %v\n", err)
		return
	}
	fmt.Println("Repositories Synchronised")

	_, err = SyncAllActivities()
	if err != nil {
		fmt.Printf("Error syncing daily activities: %v\n", err)
		return
	}
	fmt.Println("Daily Activities Synchronised")
	_, err = SyncAllRepoActivities()
	if err != nil {
		fmt.Printf("Error syncing repo activities: %v\n", err)
		return
	}
	fmt.Println("Repo Activities Synchronised")

	elapsedTime := int64(time.Since(startTime).Seconds())

	err = repository.SyncSystemSummary(orgs, repos, users, time.Now().Local(), elapsedTime)
	if err != nil {
		fmt.Printf("Error syncing system summary: %v\n", err)
		return
	}
	fmt.Println("System Summary Synchronised")

	defer func() {
		err = repository.UpdateSyncStatus(true)
		if err != nil {
			fmt.Printf("Error updating sync status: %v\n", err)
		}
	}()

	fmt.Printf("Sync completed in %d seconds\n", elapsedTime)
}

func SyncDailyData() {
	startTime := time.Now()

	err := repository.UpdateSyncStatus(false)
	if err != nil {
		fmt.Printf("Error updating sync status: %v\n", err)
		return
	}

	orgs, err := SyncAllOrganizations()
	if err != nil {
		fmt.Printf("Error syncing organizations: %v\n", err)
		return
	}
	fmt.Println("Organizations Synchronised")

	users, err := SyncAllUsers()
	if err != nil {
		fmt.Printf("Error syncing users: %v\n", err)
		return
	}
	fmt.Println("Users Synchronised")

	repos, err := SyncAllRepos()
	if err != nil {
		fmt.Printf("Error syncing repos: %v\n", err)
		return
	}
	fmt.Println("Repositories Synchronised")

	_, err = SyncDailyActivities()
	if err != nil {
		fmt.Printf("Error syncing daily activities: %v\n", err)
		return
	}
	fmt.Println("Daily Activities Synchronised")
	_, err = SyncAllNewRepoActivities()
	if err != nil {
		fmt.Printf("Error syncing repo activities: %v\n", err)
		return
	}
	fmt.Println("Repo Activities Synchronised")

	elapsedTime := int64(time.Since(startTime).Seconds())

	err = repository.SyncSystemSummary(orgs, repos, users, time.Now().Local(), elapsedTime)
	if err != nil {
		fmt.Printf("Error syncing system summary: %v\n", err)
		return
	}
	fmt.Println("System Summary Synchronised")

	defer func() {
		err = repository.UpdateSyncStatus(true)
		if err != nil {
			fmt.Printf("Error updating sync status: %v\n", err)
		}
	}()

	fmt.Printf("Sync completed in %d seconds\n", elapsedTime)
}

func InitCronScheduler() *cron.Cron {
	dhakaTimezone, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		fmt.Println("Error loading Dhaka timezone, using fallback:", err)
		dhakaTimezone = time.FixedZone("Asia/Dhaka", 6*60*60)
	}
	c := cron.New(cron.WithLocation(dhakaTimezone))
	c.AddFunc("30 3 * * *", SyncDailyData)
	c.Start()
	fmt.Println("Cron Scheduler Initialized")
	return c
}

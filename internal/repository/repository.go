package repository

import (
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
)

const (
	userCollection      = "users"
	activitesCollection = "activities"
	repoCollection      = "repos"
	systemCollection    = "system"
	orgCollection       = "orgs"
)

var (
	client = &http.Client{}
)

func AggregateCommits(lastMonday, LastMonth, LastYear string, activities []model.Activity) (int, int, int, int) {
	var weekly int
	var monthly int
	var yearly int
	var allTime int
	for _, activity := range activities {
		date := activity.Date.Format("2006-01-02")
		if date >= lastMonday {
			weekly++
		}
		if date >= LastMonth {
			monthly++
		}
		if date >= LastYear {
			yearly++
		}
		allTime++

	}
	return weekly, monthly, yearly, allTime
}

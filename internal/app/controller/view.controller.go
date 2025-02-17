package controller

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/templates/components"
	"github.com/labstack/echo/v4"
)

func RenderOrganizations(c echo.Context) error {

	orgs, err := repository.GetAllOrgs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)

	ctx := context.Background()
	if err := components.Organizations(orgs).Render(ctx, c.Response().Writer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return nil
}
func RenderHome(c echo.Context) error {
	systemSummary, err := repository.GetSystemSummary()
	var orgs int
	var repos int
	var users int
	var last_updated time.Time
	var is_synced bool
	if err != nil {
		orgs = 0
		repos = 0
		users = 0
		last_updated = time.Now()
		is_synced = true
	} else {
		orgs = systemSummary.Total_orgs
		repos = systemSummary.Total_repos
		users = systemSummary.Total_users
		last_updated = systemSummary.Last_synced
		is_synced = systemSummary.Is_synced
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)

	ctx := context.Background()
	if err := components.Home(orgs, repos, users, last_updated, is_synced).Render(ctx, c.Response().Writer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return nil
}
func RenderRepos(c echo.Context) error {

	org := c.Param("org")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	repos, err := repository.GetAllReposFromOrg(org, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)

	ctx := context.Background()
	if err := components.Repos(repos, org).Render(ctx, c.Response().Writer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func RenderUsers(c echo.Context) error {

	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	users, err := repository.GetAllUsers(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)

	ctx := context.Background()
	if err := components.Users(users).Render(ctx, c.Response().Writer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func RenderUser(c echo.Context) error {
	user, err := repository.GetUser(c.Param("user"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var startDate string
	var endDate string
	viewBy := c.QueryParam("viewBy")
	if viewBy == "" {
		viewBy = "week"
	}
	now := time.Now()

	if viewBy == "week" {
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		lastMonday := now.AddDate(0, 0, -weekday+1)
		startDate = lastMonday.Format("2006-01-02")
		endDate = now.Format("2006-01-02")

	} else if viewBy == "month" {

		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	} else if viewBy == "year" {

		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	} else {
		endDate = now.Format("2006-01-02")
	}

	activities, err := repository.GetUserActivityByDateRange(c.Param("user"), startDate, endDate, "")
	if startDate == "" {
		minDate := endDate
		for _, activity := range activities {
			if activity.Date.Format("2006-01-02") < minDate {
				minDate = activity.Date.Format("2006-01-02")
			}
		}
		startDate = minDate
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var wg sync.WaitGroup
	var dailyCommitCountList, repoCommitCountList []int
	var dateList, repoNameList []string
	wg.Add(2)
	go func() {
		dailyCommitCountList, dateList = GetCommitCountWithDateList(activities, &wg, startDate, endDate, viewBy)
	}()
	go func() {
		repoCommitCountList, repoNameList = repoCommitCountWithNameList(activities, &wg)
	}()
	wg.Wait()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	ctx := context.Background()

	if err := components.User(user, dailyCommitCountList, repoCommitCountList, repoNameList, dateList, viewBy).Render(ctx, c.Response().Writer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func RenderRepo(c echo.Context) error {
	org := c.Param("org")
	repo := c.Param("repo")
	repoUsers, err := repository.SearchUsersOfRepo(org, repo, "", "", "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	repoObj, err := repository.GetRepo(org, repo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var startDate string
	var endDate string
	viewBy := c.QueryParam("viewBy")
	if viewBy == "" {
		viewBy = "week"
	}
	now := time.Now()

	if viewBy == "week" {
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		lastMonday := now.AddDate(0, 0, -weekday+1)
		startDate = lastMonday.Format("2006-01-02")
		endDate = now.Format("2006-01-02")

	} else if viewBy == "month" {

		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	} else if viewBy == "year" {

		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	} else {

		startDate = repoObj.Created.Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	userActivities := [][]model.Activity{}
	for _, user := range repoUsers {
		activities, err := repository.GetUserActivityByDateRange(user.Username, startDate, endDate, repo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		userActivities = append(userActivities, activities)
	}

	var wg sync.WaitGroup
	dailyCommitCountListForAllUsers := make([][]int, 0)
	dateListForAllUsers := make([]string, 0)

	for _, activities := range userActivities {
		wg.Add(1)
		commitCounts, dates := GetCommitCountWithDateList(activities, &wg, startDate, endDate, viewBy)
		wg.Wait()
		dailyCommitCountListForAllUsers = append(dailyCommitCountListForAllUsers, commitCounts)
		dateListForAllUsers = dates
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	ctx := context.Background()

	if err := components.Repo(org, repoObj, repoUsers, dailyCommitCountListForAllUsers, dateListForAllUsers, viewBy).Render(ctx, c.Response().Writer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func GetCommitCountWithDateList(activities []model.Activity, wg *sync.WaitGroup, startDate, endDate, viewBy string) ([]int, []string) {
	defer wg.Done()
	dates := make([]string, 0)
	var dateParser string
	dayInc := 0
	MonthInc := 0
	if viewBy == "week" {
		dateParser = "Monday"
		dayInc = 1
	} else if viewBy == "month" {
		dateParser = "2"
		dayInc = 1
	} else if viewBy == "year" {
		dateParser = "January"
		MonthInc = 1
	} else {
		dateParser = "January 2006"
		MonthInc = 1
	}

	var startDateObj, endDateObj time.Time
	startDateObj, _ = time.Parse("2006-01-02", startDate)
	endDateObj, _ = time.Parse("2006-01-02", endDate)

	endDateObj = endDateObj.AddDate(0, MonthInc, dayInc)

	for d := startDateObj; d.Before(endDateObj); d = d.AddDate(0, MonthInc, dayInc) {
		dates = append(dates, d.Format(dateParser))

	}

	mp := make(map[string]int)
	for _, activity := range activities {
		date := activity.Date.Format(dateParser)
		mp[date] += 1

	}

	commitCounts := make([]int, 0, len(mp))

	for _, date := range dates {
		commitCounts = append(commitCounts, mp[date])
	}

	return commitCounts, dates

}

func repoCommitCountWithNameList(activities []model.Activity, wg *sync.WaitGroup) ([]int, []string) {
	defer wg.Done()
	mp := make(map[string]int)
	for _, activity := range activities {
		repo := activity.Repo.Name
		mp[repo] += 1
	}
	commitCounts := make([]int, 0, len(mp))
	repos := make([]string, 0, len(mp))
	for repo, _ := range mp {
		repos = append(repos, repo)
	}
	sort.Strings(repos)
	for _, repo := range repos {
		commitCounts = append(commitCounts, mp[repo])
	}
	fmt.Println(len(commitCounts))
	return commitCounts, repos

}

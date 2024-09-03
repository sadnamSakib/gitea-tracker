package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
)

func GetAllOrganizationsFromGitea(page int, orgs *[]*model.Org) error {
	// Define the Gitea API URL with access token
	currentPageOrgs := []model.Org{}
	limit := 50
	url := fmt.Sprintf("https://gitea.vivasoftltd.com/api/v1/orgs?page=%d&limit=%d&access_token=%s", page, limit, config.AppConfig.GITEA.API_KEY)

	// Create a new HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&currentPageOrgs); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(currentPageOrgs) == 0 {
		return nil
	}
	for _, org := range currentPageOrgs {
		*orgs = append(*orgs, &org)
	}

	fmt.Println(len(*orgs))

	return nil
}

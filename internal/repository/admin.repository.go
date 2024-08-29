package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/config"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
)

const orgCollection = "orgs"

func GetAllOrganizationsFromGitea(orgs *[]*model.Org) error {
	// Define the Gitea API URL with access token
	url := "https://gitea.vivasoftltd.com/api/v1/admin/orgs?access_token=" + config.AppConfig.GITEA.API_KEY

	// Create a new HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Make the HTTP GET request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	fmt.Println("Request Made")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response into the orgs slice (using a pointer to the slice)
	if err := json.Unmarshal(body, orgs); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}

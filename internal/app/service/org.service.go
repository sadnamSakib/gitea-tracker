package service

import (
	"fmt"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
)

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

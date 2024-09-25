package service

import "gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"

func SyncAllUsers() (int, error) {
	users, err := repository.FetchUsersFromGitea(1)
	if err != nil {
		return 0, err
	}

	err = repository.SyncUsersWithDB(users)
	if err != nil {
		return 0, err
	}

	return len(users), nil
}

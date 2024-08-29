package controller

import (
	"fmt"
	"net/http"

	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/internal/repository"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"

	"github.com/labstack/echo/v4"
)

func GetAllUsers(c echo.Context) error {
	var users []*model.User
	var orgName string = c.Param("org")
	fmt.Println("Reached Controller")
	err := repository.GetAllUserFromGitea(orgName, &users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

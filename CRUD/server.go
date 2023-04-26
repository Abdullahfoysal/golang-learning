package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	users = map[int]*user{}
	seq   = 1
	lock  = sync.Mutex{}
)

func getAllUsers(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, users)
}
func createUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}
func getUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}
func updateUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}
func deleteUser(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()

	//middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Route => handler
	e.GET("/users", getAllUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1324"))
}

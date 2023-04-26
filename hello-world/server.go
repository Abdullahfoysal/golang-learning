package main

import(
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main(){
	e := echo.New()
	
	//Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Route => handler
	e.GET("/",func (c echo.Context) error  {
		return c.String(http.StatusOK,"Heelo, World!!!")
	})

	//Start server
	e.Logger.Fatal(e.Start(":1323"))

}
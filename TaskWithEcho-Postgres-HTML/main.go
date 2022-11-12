package main

import (
	Handler "myapp/internal/handlers"
	Logic "myapp/internal/logic"

	"github.com/labstack/echo"
)

func main() {
	Logic.InitTemplate()
	router := echo.New()
	router.Renderer = Logic.T
	router.Use(Handler.ConnectDB)
	router.GET("/", Handler.Autorization)
	router.GET("/mainForm", Handler.GetPersons)
	router.GET("/Add", Handler.Add)
	router.GET("/Remove", Handler.Remove)
	router.GET("/Edit", Handler.Edit)
	router.GET("/Form_handler_Autorization", Handler.Form_handler_Autorization)
	router.GET("/Form_handler_GetById", Handler.Form_handler_GetById)
	router.POST("/Form_handler_PostPerson", Handler.Form_handler_PostPerson)
	router.GET("/Form_handler_DeleteById", Handler.Form_handler_DeleteById)
	router.GET("/Form_handler_UpdatePersonById", Handler.Form_handler_UpdatePersonById)
	router.Logger.Fatal(router.Start(":8080"))
}

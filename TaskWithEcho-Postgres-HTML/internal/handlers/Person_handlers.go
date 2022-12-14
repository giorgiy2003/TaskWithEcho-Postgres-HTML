package Handler

import (
	"log"
	Logic "myapp/internal/logic"
	Model "myapp/internal/model"
	Repository "myapp/internal/repository"
	"net/http"

	"github.com/labstack/echo"
)

func Form_handler_PostPerson(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	var newPerson Model.Person
	newPerson.Email = c.FormValue("email")
	newPerson.Phone = c.FormValue("phone")
	newPerson.FirstName = c.FormValue("firstName")
	newPerson.LastName = c.FormValue("lastName")
	err := Logic.Create(newPerson)
	if err != nil {
		log.Println(err)
		return c.Render(400, "ErrorPage", map[string]interface{}{
			"Error": err,
		})
	}
	log.Println("Добавлена новая запись")
	return c.Render(200, "returnPage", nil)
}

func GetPersons(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	Persons, err := Logic.Read()
	if err != nil {
		log.Println(err)
		return c.Render(400, "ErrorPage", map[string]interface{}{
			"Error": err,
		})
	}
	if len(Persons) == 0 {
		return c.Render(200, "InfoPage", map[string]interface{}{
			"Info": "Нет информации",
		})
	}
	c.Render(200, "Title", map[string]interface{}{"Title": "Список сотрудников"}) //Вывод заголовка
	for _, Person := range Persons {
		c.Render(200, "mainForm", map[string]interface{}{
			"Id":        Person.Id,
			"Email":     Person.Email,
			"Phone":     Person.Phone,
			"FirstName": Person.FirstName,
			"LastName":  Person.LastName,
		})
	}
	return nil
}

func Form_handler_GetById(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	id := c.FormValue("id")
	Persons, err := Logic.ReadOne(id)
	if err != nil {
		log.Println(err)
		return c.Render(400, "ErrorPage", map[string]interface{}{
			"Error": err,
		})
	}
	if len(Persons) == 0 {
		return c.Render(200, "InfoPage", map[string]interface{}{
			"Info": "Нет информации",
		})
	}
	c.Render(200, "Title", map[string]interface{}{"Title": "Список сотрудников"}) //Вывод заголовка
	for _, Person := range Persons {
		c.Render(200, "mainForm", map[string]interface{}{
			"Id":        Person.Id,
			"Email":     Person.Email,
			"Phone":     Person.Phone,
			"FirstName": Person.FirstName,
			"LastName":  Person.LastName,
		})
	}
	return nil
}

func Form_handler_DeleteById(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	id := c.FormValue("id")
	err := Logic.Delete(id)
	if err != nil {
		log.Println(err)
		return c.Render(400, "ErrorPage", map[string]interface{}{
			"Error": err,
		})
	}
	log.Printf("Запись с id = %s  успешно удалена", id)
	return c.Render(200, "returnPage", nil)
}

func Form_handler_UpdatePersonById(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	var newPerson Model.Person
	id := c.FormValue("id")
	newPerson.Email = c.FormValue("email")
	newPerson.Phone = c.FormValue("phone")
	newPerson.FirstName = c.FormValue("firstName")
	newPerson.LastName = c.FormValue("lastName")
	err := Logic.Update(newPerson, id)
	if err != nil {
		log.Println(err)
		return c.Render(400, "ErrorPage", map[string]interface{}{
			"Error": err,
		})
	}
	log.Printf("Запись с id = %s  успешно обновлена", id)
	return c.Render(200, "returnPage", nil)
}

func ConnectDB(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := Repository.OpenTable(); err != nil {
			log.Printf("Не удалось подключиться к базе данных: %v", err)
			return c.Render(500, "Connection_failed", map[string]interface{}{
				"Error": err,
			})
		}
		return next(c)
	}
}

func Form_handler_Autorization(c echo.Context) error {
	login := c.FormValue("Email")
	password := c.FormValue("Password1")
	err := Logic.Autorization(login, password)
	if err != nil {
		return c.Render(200, "Authorization", map[string]interface{}{"err": err}) //Вывод ошибки
	}
	log.Println("Авторизация прошла успешно")
	return c.Redirect(http.StatusSeeOther, "/mainForm")
}

func Autorization(c echo.Context) error {
	return c.Render(200, "Authorization", nil)
}
func Add(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	return c.Render(200, "CreatePerson", nil)
}
func Remove(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	return c.Render(200, "DeleteById", nil)
}
func Edit(c echo.Context) error {
	if err := Logic.Proverka(); err != nil {
		return c.Render(400, "Connection_failed", map[string]interface{}{
			"Error": err,
		})
	}
	return c.Render(200, "EditPerson", nil)
}

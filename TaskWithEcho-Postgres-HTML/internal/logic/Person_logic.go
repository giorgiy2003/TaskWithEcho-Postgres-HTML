package Logic

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	Model "myapp/internal/model"
	Repository "myapp/internal/repository"
	"os"
	"strconv"
	"strings"

	"html/template"
	"io"

	"github.com/labstack/echo"
)

var Auth = false

func Create(p Model.Person) error {
	p.Email = strings.TrimSpace(p.Email)
	p.Phone = strings.TrimSpace(p.Phone)
	p.FirstName = strings.TrimSpace(p.FirstName)
	p.LastName = strings.TrimSpace(p.LastName)
	if p.Email == "" || p.Phone == "" || p.FirstName == "" || p.LastName == "" {
		return errors.New("невозможно добавить запись, не все поля заполнены!")
	}
	if _, err := Repository.Connection.Exec(`INSERT INTO "person" ("person_email", "person_phone", "person_firstName", "person_lastName") VALUES ($1, $2,$3,$4)`, p.Email, p.Phone, p.FirstName, p.LastName); err != nil {
		return err
	}
	return nil
}

func ReadOne(id string) ([]Model.Person, error) {
	person_id, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("неверно введён параметр id!")
	}
	row, err := Repository.Connection.Query(`SELECT * FROM "person" WHERE "person_id" = $1`, person_id)
	if err != nil {
		return nil, err
	}
	var personInfo = []Model.Person{}
	for row.Next() {
		var p Model.Person
		err := row.Scan(&p.Id, &p.Email, &p.Phone, &p.FirstName, &p.LastName)
		if err != nil {
			return nil, err
		}
		personInfo = append(personInfo, p)
	}
	return personInfo, nil
}

func Read() ([]Model.Person, error) {
	row, err := Repository.Connection.Query(`SELECT * FROM "person" ORDER BY "person_id"`)
	if err != nil {
		return nil, err
	}
	var personInfo = []Model.Person{}
	for row.Next() {
		var p Model.Person
		err := row.Scan(&p.Id, &p.Email, &p.Phone, &p.FirstName, &p.LastName)
		if err != nil {
			return nil, err
		}
		personInfo = append(personInfo, p)
	}
	return personInfo, nil
}

func Update(p Model.Person, id string) error {
	if err := dataExist(id); err != nil {
		return err
	}
	p.Email = strings.TrimSpace(p.Email)
	p.Phone = strings.TrimSpace(p.Phone)
	p.FirstName = strings.TrimSpace(p.FirstName)
	p.LastName = strings.TrimSpace(p.LastName)
	if p.Email == "" || p.Phone == "" || p.FirstName == "" || p.LastName == "" {
		return errors.New("невозможно редактировать запись, не все поля заполнены!")
	}
	if _, err := Repository.Connection.Exec(`UPDATE "person" SET "person_email" = $1,"person_phone" = $2,"person_firstName" = $3,"person_lastName" = $4  WHERE "person_id" = $5`, p.Email, p.Phone, p.FirstName, p.LastName, id); err != nil {
		return err
	}
	return nil
}

func Delete(id string) error {
	if err := dataExist(id); err != nil {
		return err
	}
	if _, err := Repository.Connection.Exec(`DELETE FROM "person" WHERE "person_id" = $1`, id); err != nil {
		return err
	}
	return nil
}

func dataExist(id string) error {
	persons, err := ReadOne(id)
	if err != nil {
		return err
	}
	if len(persons) == 0 {
		return fmt.Errorf("записи с id = %s не существует", id)
	}
	return nil
}

//Авторизация
func Autorization(login, password string) error {
	login = strings.TrimSpace(login)
	Password_file, err := os.Open("./internal/userPassword.txt")
	if err != nil {
		log.Printf("Невозможно открыть файл userPassword: %v \n", err)
		return errors.New("Сайт временно недоступен")
	}
	defer Password_file.Close()

	Login_file, err := os.Open("./internal/userLogin.txt")
	if err != nil {
		log.Printf("Невозможно открыть файл userLogin: %v \n", err)
		return errors.New("Сайт временно недоступен")
	}
	defer Login_file.Close()

	Pass, err := io.ReadAll(Password_file)
	if err != nil {
		return err
	}
	Log, err := io.ReadAll(Login_file)
	if err != nil {
		return err
	}

	//Hashed value of password
	password = MD5_Encode(password) //Кодируем данные полученные от пользователя

	//Hashed value of login
	login = MD5_Encode(login)

	//Сравниваем с данными полученными из файлов .txt
	if login != string(Log) || password != string(Pass) {
		return errors.New("Введён неверный логин или пароль!")
	}
	Auth = true
	return nil
}

//Кодирование данных
func MD5_Encode(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func Proverka() error {
	if !Auth {
		return errors.New("Для начала работы необходимо авторизоваться!")
	}
	return nil
}

var T *Template

type Template struct {
	templates *template.Template
}

func (T *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return T.templates.ExecuteTemplate(w, name, data)
}

func InitTemplate() {
	T = &Template{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}
}

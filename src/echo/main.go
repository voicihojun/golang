package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"database/sql"
	_"github.com/lib/pq"
)

func main() {
	// fmt.Println("hello Golang")

	db, err := sql.Open("postgres", "user=postgres password=password dbname=goapi sslmode=disable")

	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB Connected...")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	type Employee struct {
		Userid		string `json:"userid"`
		Createdat	string `json:"createdat"`
		Updatedat	string `json:"updatedat"`
		Username	string `json:"username"`
	}

	e.GET("/employee/:userid", func(c echo.Context) error {
		userid := c.Param("userid")
		sqlStatement := "SELECT userid, createdat, updatedat, username FROM employee WHERE userid = $1"
		rows, err := db.Query(sqlStatement, userid)
		if err != nil {
			return err
		}
		defer rows.Close()
		employee := []Employee{}

		for rows.Next() {
			e := Employee{}
			err2 := rows.Scan(&e.Userid, &e.Createdat, &e.Updatedat, &e.Username)
			if err2 != nil {
				return err2
			}
			employee = append(employee, e)
		}
		return c.JSON(http.StatusCreated, employee[0])
	})

	e.PUT("/employee", func(c echo.Context) error {
		u := new(Employee)
		
		if err := c.Bind(u); err != nil {
			return err
		}
		fmt.Println(u)
		sqlStatement := "UPDATE employee SET updatedat=now(), username=$2 WHERE userid=$1"
		res, err := db.Query(sqlStatement, u.Userid, u.Username)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, u.Userid)
	})

	e.POST("/employee", func(c echo.Context) error {
		u := new(Employee)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "INSERT INTO employee (userid, createdat, updatedat, username) VALUES ($1, now(), now(), $2)"
		res, err := db.Query(sqlStatement, u.Userid, u.Username)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("hello")
			fmt.Println(res)
			
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, "ok")
	})
	

	e.Logger.Fatal(e.Start(":1800"))


}
package main

import (
	"fmt"
	"lesson4/DAO"
	"lesson4/controllers"
	"lesson4/database"
	"lesson4/services"
)

func main() {
	db := database.InitDB()
	//test := testFunction.NewTestUserDAO(DAO.NewUserDAO(db))
	//test.Test()
	fmt.Println(DAO.NewSelectionDAO(db).ListSelections())
	admCon := services.NewAdministerService(controllers.NewAdministerUser(db), controllers.NewStudentUser(db))
	admCon.Serve()

}

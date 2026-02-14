package main

import (
	"collegeWaleServer/internal/database"
	"collegeWaleServer/internal/model"
	"log"
)

var modelsToMigrate = []any{
	&model.College{},
	&model.Attendance{},
	&model.Student{},
	&model.Subject{},
	&model.User{},
}

func main() {
	dbService := database.New()
	err := dbService.DB.AutoMigrate(modelsToMigrate...)

	if err != nil {
		log.Fatalf("failed to migrate db %v", err)
	}

	log.Println("Database migration completed successfully")

}

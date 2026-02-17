package main

import (
	"collegeWaleServer/internal/db"
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
	dbService := db.New()
	err := dbService.GetDatabase().AutoMigrate(modelsToMigrate...)

	if err != nil {
		log.Fatalf("failed to migrate db %v", err)
	}

	log.Println("Database migration completed successfully")

}

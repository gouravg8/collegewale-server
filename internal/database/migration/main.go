package main

import (
	"collegeWaleServer/internal/database"
	"log"
)

var modelsToMigrate = []any{
	// &models.College{},
	// &models.Attendance{},
	// &models.Student{},
	// &models.Subject{},
	// &models.User{},
}

func main() {
	dbService := database.New()
	err := dbService.DB.AutoMigrate(modelsToMigrate...)

	if err != nil {
		log.Fatalf("failed to migrate db %v", err)
	}

	log.Println("Database migration completed successfully")

}

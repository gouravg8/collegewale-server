package main

import (
	"collegeWaleServer/internal/database"
	"collegeWaleServer/internal/models"
	"log"
)

var modelsToMigrate = []any{
	&models.College{},
}

func main() {
	dbService := database.New()
	err := dbService.DB.AutoMigrate(modelsToMigrate...)

	if err != nil {
		log.Fatalf("failed to migrate db %v", err)
	}

	log.Println("Database migration completed successfully")

}

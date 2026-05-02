package main

import (
	"collegeWaleServer/internal/db"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/utils"
	"encoding/base64"
	"log"
)

var modelsToMigrate = []any{
	&model.College{},
	&model.Attendance{},
	&model.Student{},
	&model.Subject{},
	&model.User{},
	&model.Attendance{},
	&model.Courses{},
	&model.Role{},
	&model.Application{},
}

func main() {
	dbService := db.New()
	gormDB := dbService.GetDatabase()

	if err := gormDB.AutoMigrate(modelsToMigrate...); err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}
	log.Println("Database migration completed successfully")

	seedRoles(dbService)
	seedSuperAdmin(dbService)
}

func seedRoles(dbService *db.Service) {
	gormDB := dbService.GetDatabase()
	allRoles := []roles.Roles{roles.Admin, roles.Student, roles.College}
	for _, r := range allRoles {
		role := model.Role{Name: r}
		result := gormDB.Where("name = ?", r).FirstOrCreate(&role)
		if result.Error != nil {
			log.Fatalf("failed to seed role %s: %v", r, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("created role: %s", r)
		}
	}
	log.Println("Roles seeded successfully")
}

func seedSuperAdmin(dbService *db.Service) {
	gormDB := dbService.GetDatabase()

	// Check if super_admin already exists
	var existing model.User
	if err := gormDB.Where("username = ?", "super_admin").First(&existing).Error; err == nil {
		log.Println("super_admin already exists, skipping seed")
		return
	}

	// Hash "Admin@123" — base64-encode first to match the app's HashPassword convention
	encoded := base64.StdEncoding.EncodeToString([]byte("Admin@123"))
	hash, err := utils.HashPassword(encoded)
	if err != nil {
		log.Fatalf("failed to hash super admin password: %v", err)
	}

	var adminRole model.Role
	if err := gormDB.Where("name = ?", roles.Admin).First(&adminRole).Error; err != nil {
		log.Fatalf("admin role not found — run seedRoles first: %v", err)
	}

	superAdmin := model.User{
		Username:     "super_admin",
		Email:        "super_admin@collegewale.com",
		PasswordHash: hash,
	}

	if err := gormDB.Create(&superAdmin).Error; err != nil {
		log.Fatalf("failed to create super_admin: %v", err)
	}

	// Explicitly insert into user_roles join table
	if err := gormDB.Model(&superAdmin).Association("Roles").Append(&adminRole); err != nil {
		log.Fatalf("failed to assign admin role to super_admin: %v", err)
	}
	log.Println("super_admin created successfully")
	log.Println("  username : super_admin")
	log.Println("  password : Admin@123")
}


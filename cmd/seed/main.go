package main

import (
	"collegeWaleServer/internal/db"
	"collegeWaleServer/internal/enums/roles"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/utils"
	"encoding/base64"
	"log"

	"gorm.io/gorm"
)

func encodedPassword(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func ensureRole(gdb *gorm.DB, name roles.Roles) model.Role {
	var role model.Role
	if err := gdb.Where("name = ?", name).FirstOrCreate(&role, model.Role{Name: name}).Error; err != nil {
		log.Fatalf("failed to ensure role %q: %v", name, err)
	}
	return role
}

func ensureUser(gdb *gorm.DB, username, email, rawPassword string, role model.Role, collegeID *uint) *model.User {
	var existing model.User
	err := gdb.Where("username = ?", username).First(&existing).Error
	if err == nil {
		log.Printf("seed user %q already exists, skipping", username)
		return &existing
	}

	passwordHash, err := utils.HashPassword(encodedPassword(rawPassword))
	if err != nil {
		log.Fatalf("failed to hash password for %q: %v", username, err)
	}

	user := model.User{
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		Roles:        []model.Role{role},
		CollegeID:    collegeID,
	}
	if err := gdb.Create(&user).Error; err != nil {
		log.Fatalf("failed to create seed user %q: %v", username, err)
	}
	log.Printf("seed user created: username=%s password=%s role=%s", username, rawPassword, role.Name)
	return &user
}

func ensureStudentProfile(gdb *gorm.DB, user *model.User, rollNumber string, courseID uint) {
	var existing model.Student
	err := gdb.Where("user_id = ?", user.ID).First(&existing).Error
	if err == nil {
		log.Printf("student profile for %q already exists, skipping", user.Username)
		return
	}

	student := model.Student{
		FirstName:  "Rohan",
		LastName:   "Verma",
		Email:      user.Email,
		Phone:      "9876500000",
		RollNumber: rollNumber,
		CourseID:   courseID,
		Year:       3,
		Gender:     "male",
		Semester:   "5",
		UserID:     user.ID,
	}
	if err := gdb.Create(&student).Error; err != nil {
		log.Fatalf("failed to create student profile for %q: %v", user.Username, err)
	}
	log.Printf("student profile created for %q (roll=%s)", user.Username, rollNumber)
}

func ensureCourse(gdb *gorm.DB, name string) {
	var existing model.Courses
	err := gdb.Where("name = ?", name).First(&existing).Error
	if err == nil {
		return
	}
	if err := gdb.Create(&model.Courses{Name: name}).Error; err != nil {
		log.Fatalf("failed to create course %q: %v", name, err)
	}
	log.Printf("seed course created: %s", name)
}

func main() {
	dbService := db.New()
	gdb := dbService.GetDatabase()

	adminRole := ensureRole(gdb, roles.Admin)
	collegeAdminRole := ensureRole(gdb, roles.CollegeAdmin)
	staffRole := ensureRole(gdb, roles.Staff)
	studentRole := ensureRole(gdb, roles.Student)

	for _, course := range []string{"COMPUTER SCIENCE", "ELECTRONICS", "MECHANICAL", "CIVIL", "MATHEMATICS"} {
		ensureCourse(gdb, course)
	}

	// Platform super-admin: creates colleges and college-admin accounts.
	ensureUser(gdb, "superadmin", "superadmin@brightfield.edu", "admin123", adminRole, nil)

	// Demo college + its college admin, staff, and student for local end-to-end testing.
	var college model.College
	err := gdb.Where("code = ?", "BFC001").First(&college).Error
	if err != nil {
		college = model.College{
			Name:  "Brightfield College",
			Code:  "BFC001",
			Phone: "9999999999",
			Email: "college@brightfield.edu",
			Seats: 500,
		}
		if err := gdb.Create(&college).Error; err != nil {
			log.Fatalf("failed to create demo college: %v", err)
		}
		log.Printf("demo college created: code=%s", college.Code)
	} else {
		log.Printf("demo college %q already exists, skipping", college.Code)
	}

	ensureUser(gdb, "collegeadmin", "collegeadmin@brightfield.edu", "admin123", collegeAdminRole, &college.ID)
	ensureUser(gdb, "staffuser", "staff@brightfield.edu", "staff123", staffRole, &college.ID)
	studentUser := ensureUser(gdb, "studentuser", "student@brightfield.edu", "student123", studentRole, &college.ID)

	var csCourse model.Courses
	if err := gdb.Where("name = ?", "COMPUTER SCIENCE").First(&csCourse).Error; err == nil {
		ensureStudentProfile(gdb, studentUser, "CS21B045", csCourse.ID)
	}

	log.Println("seeding complete")
}

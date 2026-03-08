package views

import "collegeWaleServer/internal/model"

type CoursesFilter struct {
	Name         string `json:"name"`
	WithSubjects bool   `json:"with_subjects"`
}

type CoursesResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Subjects    []string `json:"subjects,omitempty"`
}

func NewCoursesResponse(c model.Courses) CoursesResponse {
	response := CoursesResponse{
		Name:        c.Name,
		Description: c.Description,
	}
	response.Subjects = make([]string, 0)
	for _, s := range c.Subjects {
		response.Subjects = append(response.Subjects, s.Name)
	}
	return response
}

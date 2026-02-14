package service

import (
	"collegeWaleServer/errz"
	"collegeWaleServer/internal/model"
	"collegeWaleServer/internal/views"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type RegistryService struct {
	db *gorm.DB
}

func NewRegistryService(db *gorm.DB) *RegistryService {
	return &RegistryService{db}
}

func (s *RegistryService) RegisterCollege(req views.College) error {
	// --- Input Validation ---
	var existingCount int64
	if err := s.db.Model(&model.College{}).Where("code = ?", req.Code).Count(&existingCount).Error; err != nil {
		return err
	} else if existingCount > 0 {
		return errz.NewBadRequest("college already exists")
	}

	clg := model.College{
		Name:       req.Name,
		Code:       req.Code,
		Phone:      req.Phone,
		Email:      req.Email,
		CourseType: req.CourseType,
		Seats:      req.Seats,
		Logo:       req.Logo,
	}

	if err := s.db.Model(&model.College{}).Create(&clg).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errz.NewBadRequest(pgErr.Detail)
		}
		return err
	}

	return nil
}

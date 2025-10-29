package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type UsGeographyRepo struct {
	*gormx.BaseRepo[*model.UsGeography]
	m *model.UsGeography
}

func NewUsGeographyRepo(db *gorm.DB) *UsGeographyRepo {
	return &UsGeographyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.UsGeography](db),
		m:        new(model.UsGeography),
	}
}

// GetStateUsps get us_geography unique state_usps
func (repo *UsGeographyRepo) GetStateUsps(ctx context.Context) (result []string, err error) {
	err = repo.Client(ctx).
		Select("state_usps").
		Distinct("state_usps").
		Order("state_usps").
		Find(&result).Error
	// slice.Unique(result)
	// slice.Sort(result)
	return result, errors.WithStack(err)
}

// GetCitiesByStateUsps get city name by state usps
func (repo *UsGeographyRepo) GetCitiesByStateUsps(ctx context.Context, stateUsps string) (result []string, err error) {
	err = repo.Client(ctx).
		Select("city_name").
		Where("state_usps = ?", stateUsps).
		Distinct("city_name").
		Order("city_name").
		Find(&result).Error
	// slice.Unique(result)
	// slice.Sort(result)
	return result, errors.WithStack(err)
}

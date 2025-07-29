package repo

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/app/borrower/rpc/internal/repo/model"
	"zeebee-backend/common/gormx"
)

type PropertyRepo struct {
	*gormx.BaseRepo[*model.Property]
	m *model.Property
}

func NewPropertyRepo(db *gorm.DB) *PropertyRepo {
	return &PropertyRepo{
		BaseRepo: gormx.NewBaseRepo[*model.Property](db),
		m:        new(model.Property),
	}
}

// GetOrCreateWithTx query or create property
func (repo *PropertyRepo) GetOrCreateWithTx(ctx context.Context, tx *gorm.DB, data *model.Property) (result *model.Property, err error) {
	// if not existed,create one
	err = tx.WithContext(ctx).
		Where("borrower_id = ? and property_type = ? and unit = ?", data.BorrowerId, data.PropertyType, data.Unit).
		Where("address->>'state' = ? and address->>'city' = ? and address->>'zip_code' = ? and address->>'address1' = ? and address->>'address2' = ?", data.Address.State, data.Address.City, data.Address.ZipCode, data.Address.Address1, data.Address.Address2).
		Find(&result).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		err = tx.WithContext(ctx).Create(data).Error
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = data
	}
	return
}

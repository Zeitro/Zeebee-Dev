package repo

import (
	"context"
	"fmt"

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

func (repo *PropertyRepo) GetOneWithTx(ctx context.Context, tx *gorm.DB, id int64, preloads []string) (result *model.Property, err error) {
	query := tx.WithContext(ctx).Where(repo.m.FieldId(), id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err = query.First(&result).Error
	return
}

// GetOrCreateWithTx query or create property
func (repo *PropertyRepo) GetOrCreateWithTx(ctx context.Context, tx *gorm.DB, data *model.Property) (result *model.Property, err error) {
	// if not existed,create one
	err = tx.WithContext(ctx).
		Where(repo.m.FieldBorrowerId(), data.BorrowerId).
		Where(repo.m.FieldPropertyType(), data.PropertyType).
		Where(repo.m.FieldUnit(), data.Unit).
		Where(fmt.Sprintf("%s->>'state'", repo.m.FieldAddress()), data.Address.State).
		Where(fmt.Sprintf("%s->>'city'", repo.m.FieldAddress()), data.Address.City).
		Where(fmt.Sprintf("%s->>'zip_code'", repo.m.FieldAddress()), data.Address.ZipCode).
		Where(fmt.Sprintf("%s->>'address1'", repo.m.FieldAddress()), data.Address.Address1).
		Where(fmt.Sprintf("%s->>'address2'", repo.m.FieldAddress()), data.Address.Address2).
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

func (repo *PropertyRepo) GetOneWithTxByProperty(ctx context.Context, tx *gorm.DB, data *model.Property) (result *model.Property, err error) {
	err = tx.WithContext(ctx).
		Where(repo.m.FieldBorrowerId(), data.BorrowerId).
		Where(repo.m.FieldPropertyType(), data.PropertyType).
		Where(repo.m.FieldUnit(), data.Unit).
		Where(fmt.Sprintf("%s->>'state'", repo.m.FieldAddress()), data.Address.State).
		Where(fmt.Sprintf("%s->>'city'", repo.m.FieldAddress()), data.Address.City).
		Where(fmt.Sprintf("%s->>'zip_code'", repo.m.FieldAddress()), data.Address.ZipCode).
		Where(fmt.Sprintf("%s->>'address1'", repo.m.FieldAddress()), data.Address.Address1).
		Where(fmt.Sprintf("%s->>'address2'", repo.m.FieldAddress()), data.Address.Address2).
		Find(&result).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result == nil || result.Id == 0 {
		return nil, nil
	}
	return
}

func (repo *PropertyRepo) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, id int64, data *model.Property) (result *model.Property, err error) {
	err = tx.WithContext(ctx).Select("*").
		Omit(repo.m.FieldId(), repo.m.FieldCreatedAt()).
		Where(repo.m.FieldId(), id).
		Updates(data).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

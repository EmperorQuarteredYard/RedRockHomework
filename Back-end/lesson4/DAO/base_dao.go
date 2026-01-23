package DAO

import (
	"gorm.io/gorm"
)

type BaseDAO struct {
	db *gorm.DB
}

func NewBaseDAO(db *gorm.DB) *BaseDAO {
	return &BaseDAO{db: db}
}

func (d *BaseDAO) Create(model interface{}) error {
	return d.db.Create(model).Error
}

func (d *BaseDAO) GetByID(model interface{}, id int64) error {
	return d.db.First(model, id).Error
}

func (d *BaseDAO) Update(model interface{}) error {
	return d.db.Save(model).Error
}

func (d *BaseDAO) Delete(model interface{}) error {
	return d.db.Delete(model).Error
}

func (d *BaseDAO) List(model interface{}, list *[]interface{}) error {
	err := d.db.Find(model, list).Error
	return err
}

func (d *BaseDAO) Exist(model interface{}, id int64) (bool, error) {
	var count int64
	err := d.db.Model(model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

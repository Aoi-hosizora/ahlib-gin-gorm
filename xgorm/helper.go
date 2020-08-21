package xgorm

import (
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/jinzhu/gorm"
)

type Helper struct {
	db *gorm.DB
}

func WithDB(db *gorm.DB) *Helper {
	return &Helper{db: db}
}

func (h *Helper) Pagination(limit int32, page int32) *gorm.DB {
	return h.db.Limit(limit).Offset((page - 1) * limit)
}

func (h *Helper) Count(model interface{}, where interface{}) (int, error) {
	cnt := 0
	rdb := h.db.Model(model).Where(where).Count(&cnt)
	return cnt, rdb.Error
}

func (h *Helper) Exist(model interface{}, where interface{}) (bool, error) {
	cnt, err := h.Count(model, where)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (h *Helper) Create(model interface{}, object interface{}) (xstatus.DbStatus, error) {
	rdb := h.db.Model(model).Create(object)
	return CreateDB(rdb)
}

func (h *Helper) Update(model interface{}, where interface{}, object interface{}) (xstatus.DbStatus, error) {
	if where == nil {
		where = object
	}
	rdb := h.db.Model(model).Where(where).Update(object)
	return UpdateDB(rdb)
}

func (h *Helper) Delete(model interface{}, where interface{}, object interface{}) (xstatus.DbStatus, error) {
	if where == nil {
		where = object
	}
	rdb := h.db.Model(model).Where(where).Delete(object)
	return DeleteDB(rdb)
}

func CreateDB(rdb *gorm.DB) (xstatus.DbStatus, error) {
	if IsMySqlDuplicateEntryError(rdb.Error) {
		return xstatus.DbExisted, nil
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return xstatus.DbFailed, rdb.Error
	}

	return xstatus.DbSuccess, nil
}

func UpdateDB(rdb *gorm.DB) (xstatus.DbStatus, error) {
	if IsMySqlDuplicateEntryError(rdb.Error) {
		return xstatus.DbExisted, nil
	} else if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if rdb.RowsAffected == 0 {
		return xstatus.DbNotFound, nil
	}

	return xstatus.DbSuccess, nil
}

func DeleteDB(rdb *gorm.DB) (xstatus.DbStatus, error) {
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if rdb.RowsAffected == 0 {
		return xstatus.DbNotFound, nil
	}

	return xstatus.DbSuccess, nil
}

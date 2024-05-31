package dbutil

import (
	"gorm.io/gorm"
)

func PanicError(db *gorm.DB) *gorm.DB {
	if err := db.Error; err != nil {
		panic(err)
	}
	return db
}

func MultiCondition(db *gorm.DB, conditions []any) *gorm.DB {
	for _, condition := range conditions {
		if condition != nil {
			cond, ok := condition.([]any)
			if ok {
				if len(cond) > 1 {
					db = db.Where(cond[0], cond[1:]...)
				}
			} else {
				db = db.Where(condition)
			}
		}
	}
	return db
}

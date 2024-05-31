package model

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"nacos/database"
	"nacos/database/dbutil"
	"nacos/util"
)

type Page struct {
	OrderBy  string
	PageNum  int `form:"pageNo" json:"pageNo"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type PageResult struct {
	PageNum int `json:"pageNumber"`
	Pages   int `json:"pagesAvailable"`
	Total   int `json:"totalCount"`
	List    any `json:"pageItems"`
}

func Paginate[M any](condition any, page *Page) PageResult {
	return PaginateResult[M, M]([]any{condition}, page)
}

func PaginateResult[M, R any](condition any, page *Page) PageResult {
	var model M
	conditions, ok := condition.([]any)
	if !ok {
		conditions = append(conditions, condition)
	}
	countDB := dbutil.MultiCondition(db.GORM.Model(model), conditions)
	queryDB := dbutil.MultiCondition(db.GORM.Model(model), conditions)
	return paginate[R](countDB, queryDB, page)
}

func PaginateSQL[R any](sql string, args []any, page *Page) PageResult {
	countDB := db.GORM.Raw(fmt.Sprintf("select count(1) from (%s) temp_table", sql), args...)
	queryDB := db.GORM.Raw(sql, args...)
	return paginate[R](countDB, queryDB, page)
}

func paginate[R any](countDB *gorm.DB, queryDB *gorm.DB, page *Page) PageResult {
	rewritePage(page)
	pageResult := PageResult{
		PageNum: page.PageNum,
		List:    []struct{}{},
	}
	var count int64
	countDB.Count(&count)
	pageResult.Total = int(count)
	pages := int(math.Ceil(float64(count) / float64(page.PageSize)))
	pageResult.Pages = pages
	if count == 0 || pageResult.PageNum > pages {
		return pageResult
	}
	var data []R
	offset := (page.PageNum - 1) * page.PageSize
	queryDB.Order(page.OrderBy).Offset(offset).Limit(page.PageSize).Find(&data)
	pageResult.List = data
	return pageResult
}

func rewritePage(page *Page) {
	pageNum := page.PageNum
	pageSize := page.PageSize
	pageNum = util.ConditionalExpression(pageNum <= 0, 1, pageNum)
	pageSize = util.ConditionalExpression(pageSize <= 0, 10, util.ConditionalExpression(pageSize > 100, 100, pageSize))
	page.PageNum = pageNum
	page.PageSize = pageSize
}

package paginator

import (
	"gorm.io/gorm"
)

// Page 标准分页结构体，接收最原始的DO
// 使用示例：
// page := paginator.Page[User]{PageNum: 1, PageSize: 10}
// err := page.SelectPages(db.Model(&User{}))
type Page[T any] struct {
	PageNum  int64 `json:"pageNum"`  // 当前页码
	PageSize int64 `json:"pageSize"` // 每页大小
	Total    int64 `json:"total"`    // 总记录数
	Pages    int64 `json:"pages"`    // 总页数
	Data     []T   `json:"data"`     // 查询结果数据
}

// SelectPages 各种查询条件先在query设置好后再放进来
// 使用示例：
// query := db.Model(&User{}).Where("age > ?", 18)
// page := paginator.Page[User]{PageNum: 1, PageSize: 10}
// err := page.SelectPages(query)
func (page *Page[T]) SelectPages(query *gorm.DB) (e error) {
	var model T
	// 获取总记录数
	query.Model(&model).Count(&page.Total)
	if page.Total == 0 {
		page.Data = []T{} // 没有数据，返回空列表
		return
	}
	// 查询当前页数据并分页
	e = query.Model(&model).Scopes(paginate(page)).Find(&page.Data).Error
	return
}

// Paginate 分页逻辑
func paginate[T any](page *Page[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 使用 PageNum 代替 CurrentPage
		if page.PageNum <= 0 {
			page.PageNum = 1 // 如果页码小于等于0，设为1
		}
		// 检查 PageSize 的合法性
		switch {
		case page.PageSize > 10000:
			page.PageSize = 10000 // 限制一下分页大小
		case page.PageSize <= 0:
			page.PageSize = 10 // 默认每页10条
		}
		// 计算总页数
		page.Pages = page.Total / page.PageSize
		if page.Total%page.PageSize != 0 {
			page.Pages++
		}
		// 页码修正，避免超过最大页数
		p := page.PageNum
		if page.PageNum > page.Pages {
			p = page.Pages
		}
		// 计算偏移量
		size := page.PageSize
		offset := int((p - 1) * size)
		return db.Offset(offset).Limit(int(size))
	}
}

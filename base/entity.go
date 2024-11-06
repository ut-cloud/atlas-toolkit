package base

import (
	"context"
	"github.com/go-kratos/kratos/v2/metadata"
	"gorm.io/gorm"
	"time"
)

type Entity struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreateBy  string         `gorm:"type:varchar(100);not null;comment:创建人" json:"create_by"`
	UpdateBy  string         `gorm:"type:varchar(100);not null;comment:更新人" json:"update_by"`
}

func (m *Entity) insertEntity(ctx context.Context) {
	var uuid string
	if md, ok := metadata.FromServerContext(ctx); ok {
		uuid = md.Get("identity")
	}
	m.CreateBy = uuid
	m.CreatedAt = time.Now()
}

func (m *Entity) updateEntity(ctx context.Context) {
	var uuid string
	if md, ok := metadata.FromServerContext(ctx); ok {
		uuid = md.Get("identity")
	}
	m.UpdateBy = uuid
	m.UpdatedAt = time.Now()
}

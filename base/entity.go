package base

import (
	"context"
	"github.com/ut-cloud/atlas-toolkit/utils"
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

func (m *Entity) InsertEntity(ctx context.Context) {
	if m.CreateBy == "" {
		m.CreateBy = utils.GetLoginUserId(ctx)
	}
}

func (m *Entity) UpdateEntity(ctx context.Context) {
	if m.UpdateBy == "" {
		m.UpdateBy = utils.GetLoginUserId(ctx)
	}
}

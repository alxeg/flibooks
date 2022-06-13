package orm

import (
	"github.com/alxeg/flibooks/pkg/utils"
	"gorm.io/gorm"
)

type Author struct {
	ID   uint   `gorm:"primary_key"`
	Name string `json:"name" gorm:"not null;unique_index"`
}

func (auth *Author) AfterFind(db *gorm.DB) error {
	auth.Name = utils.UpperInitialAll(auth.Name)
	return nil
}

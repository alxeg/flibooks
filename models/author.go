package models

import (
    "github.com/alxeg/flibooks/utils"
)

type Author struct {
    ID   int
    Name string `json:"name" sql:"not null;unique_index"`
}

func (author *Author) AfterFind() {
    author.Name = utils.UpperInitialAll(author.Name)
}

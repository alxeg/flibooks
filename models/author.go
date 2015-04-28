package models

import (
    _ "github.com/alxeg/flibooks/utils"
)

type Author struct {
    ID   int
    Name string `json:"name" sql:"not null;unique_index"`
}

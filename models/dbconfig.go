package models

type DBConfig struct {
    DBType   string `json:"db-type"`
    DBParams string `json:"db-params"`
}

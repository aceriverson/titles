package core

import (
	"titles.run/site_api/services/interfaces"
)

type Core struct {
	DB interfaces.DBService
}

func NewSiteCore(db interfaces.DBService) *Core {
	core := &Core{
		DB: db,
	}

	return core
}

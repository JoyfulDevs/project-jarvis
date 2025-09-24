package app

import (
	"github.com/jh1104/publicapi"
	"github.com/jh1104/publicapi/forecast"
	"github.com/jh1104/publicapi/specialday"
)

type DataPortalService struct {
}

func NewDataPortalService(authKey string) *DataPortalService {
	c := publicapi.NewClient(authKey)
	specialday.SetDefaultClient(c)
	forecast.SetDefaultClient(c)
	return &DataPortalService{}
}

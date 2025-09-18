package app

import (
	"github.com/jh1104/publicapi"
	"github.com/jh1104/publicapi/specialday"
)

type DataPortalService struct {
}

func NewDataPortalService(authKey string) *DataPortalService {
	specialday.SetDefaultClient(publicapi.NewClient(authKey))
	return &DataPortalService{}
}

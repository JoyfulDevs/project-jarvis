package app

import "github.com/devafterdark/project-jarvis/pkg/dataportal"

type DataPortalService struct {
	client *dataportal.Client
}

func NewDataPortalService(authKey string) *DataPortalService {
	return &DataPortalService{
		client: dataportal.NewClient(authKey),
	}
}

package app

import "github.com/genians/endpoint-lab-slack-bot/pkg/dataportal"

type DataPortalService struct {
	client *dataportal.Client
}

func NewDataPortalService(authKey string) *DataPortalService {
	return &DataPortalService{
		client: dataportal.NewClient(authKey),
	}
}

package service

import "github.com/eirsyl/flexit/app"

type webhookService struct {

}

func NewWebhookService() *webhookService {
	return &webhookService{}
}

func (ws *webhookService) Run(cnf app.Config) error {
	return nil
}

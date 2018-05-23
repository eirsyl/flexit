package debugserver

import (
	"errors"
	"github.com/eirsyl/flexit/log"
	"net/http"
)

type HealthCheck interface {
	GetName() string
	Healty() error
}

type Checker interface {
	CheckHealth() error
	Handler() http.Handler
}

type basicChecker struct {
	logger log.Logger
	checks []HealthCheck
}

var (
	HealthChecksFailed = errors.New("Health checks failed")
)

func NewChecker(logger log.Logger, checks ...HealthCheck) Checker {
	return &basicChecker{logger, checks}
}

func (bc *basicChecker) CheckHealth() error {
	var failed = false
	for _, checker := range bc.checks {
		if err := checker.Healty(); err != nil {
			bc.logger.WithField("check", checker.GetName()).Error("health check failed", err)
			failed = true
		}
	}

	if failed {
		return HealthChecksFailed
	}
	return nil
}

func (bc *basicChecker) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := bc.CheckHealth(); err != nil {
			http.Error(w, "Failed", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK"))
	})
}

package app

type App interface {
	GetName() string
	Run(srv Service) error
}

type Config interface {

}

type Service interface {
	Run(cnf Config) error
}
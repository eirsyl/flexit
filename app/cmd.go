package app

type CmdApp struct {
	name string
}

func New(name string) App {
	return &CmdApp{
		name: name,
	}
}

func (ca *CmdApp) GetName() string {
	return ca.name
}

func (ca *CmdApp) Run(srv Service) error {
	return nil
}

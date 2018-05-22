package app

import (
	"strings"
)

type App interface {
	GetName() string
	GetShortName() string
	GetDescription() string
}

type app struct {
	name        string
	description string
}

func NewApp(name, description string) App {
	return &app{
		name:        name,
		description: description,
	}
}

func (a *app) GetName() string {
	return a.name
}

func (a *app) GetShortName() string {
	fullName := a.name
	s := strings.Split(fullName, ".")
	if len(s) == 0 {
		return ""
	}
	return s[len(s)-1]
}

func (a *app) GetDescription() string {
	return a.description
}

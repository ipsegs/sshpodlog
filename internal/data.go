package data

import (
	"log"
)

type ClientConfig struct {
	Server        string
	Port          int
	Username      string
	PrivateKey    string
	KctlCtxSwitch string
}

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Config   ClientConfig
}

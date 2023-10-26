package data

import (
	"log"
)

type ClientConfig struct {
	Server        string
	Username      string
	PrivateKey    string
	KctlCtxSwitch string
	Port          int
}

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Config   ClientConfig
}

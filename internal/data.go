package data

import (
	"log"
)
type Config struct {
	Server        string
	Port          int
	Username      string
	KctlCtxSwitch string
	PrivateKey    string
}

type Application struct {
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	Config    Config
	Namespace string
}
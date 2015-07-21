package ricecore

import (
	"os/user"
)

var rdbDir string
var homeDir string

//Initializes ricecore, setting global variables, etc.
func InitCore() {
	usr, _ := user.Current()
	homeDir = usr.HomeDir
	rdbDir = homeDir + "/.rdb/"
}

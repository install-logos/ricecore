package ricecore

import (
	"encoding/json"
	"os"
	"os/user"
)

var rdbDir string

//A rice object represents a rice object
//The name
type Rice struct {
    Name    string `json:"name"`
	Program string `json:"prog"`
	Root    string `json:"root"`
	Files   map[string]string `json:"files"`
}

//A package object represents a package object which will be used for querys
//and uploads. Essentially, packages are stored online, rices are stored locally
type Package struct {
	Name       string
	Desc       string
	Author     string
	Screenshot string
	Program    string
	URL        string
}

//Initializes ricecore, setting global variables, etc.
func InitCore() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	rdbDir = dir + "/.rdb/"
}

//Creates a rice and stores the information as a json
func CreateRice(name string, prog string, root string, files map[string]string) (rice *Rice, err error) {
	rice = &Rice{Name: name, Program: prog, Root: root, Files: files}

	jsonData, err := json.MarshalIndent(rice, "", "  ")
	if err != nil {
		return nil, err
	}

	riceDir := rdbDir + prog + "/" + name + "/"
	if !exists(riceDir) {
		os.MkdirAll(riceDir, 0764)
	}

	jsonFile, err := os.Create(riceDir + "rice.json")

	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	return rice, nil
}

//Loads a rice.json file from an existing rice and returns the struct
func GetRice(name string, prog string) (rice *Rice, err error) {
	riceDir := rdbDir + prog + "/" + name + "/"

    jsonFile, err := os.Open(riceDir + "rice.json")
    if err != nil {
        return nil, err
    }

    rice = new(Rice)

    jsonParser := json.NewDecoder(jsonFile)
    if err = jsonParser.Decode(&rice); err != nil {
        return nil, err
    }

    return rice, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

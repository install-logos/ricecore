package ricecore

import (
    "errors"
	"encoding/json"
	"os"
	"os/user"
)

var rdbDir string
var homeDir string

//A rice struct represents a rice with its files and other basic info
type Rice struct {
    Name    string `json:"name"`
	Program string `json:"prog"`
	Root    string `json:"root"`
	Files   []*RiceFile `json:"files"`
}

type RiceFile struct {
    Location string `json:"location"`
    File     string `json:"file"`
}

//Initializes ricecore, setting global variables, etc.
func InitCore() {
	usr, _ := user.Current()
	homeDir = usr.HomeDir
	rdbDir = homeDir + "/.rdb/"
}

//Creates a rice and stores the information as a json
func CreateRice(name string, prog string, root string, files []*RiceFile) (rice *Rice, err error) {
	rice = &Rice{Name: name, Program: prog, Root: root, Files: files}

	jsonData, err := json.MarshalIndent(rice, "", "  ")
	if err != nil {
		return nil, err
	}

	riceDir := rdbDir + prog + "/" + name + "/"
	if !exists(riceDir) {
		os.MkdirAll(riceDir, 0664)
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

//Initializes a created local rice, extracting the files from the directory to 
//the rdb dir and symlinking them back
func (rice Rice) InitLocalRice() (err error){
	riceDir := rdbDir + rice.Program + "/" + rice.Name + "/"
	progDir := expandDir(rice.Root)

    if err:= os.Chdir(riceDir); err != nil {
        return errors.New("Error, this rice was not created properly")
    }

    for _, rf := range rice.Files {
        if !exists(rf.Location){
            os.MkdirAll(rf.Location, 0755)
        }

        if err = os.Rename(progDir + rf.Location + rf.File, riceDir + rf.Location + rf.File); err != nil {
            return errors.New("Error, this rice was not initialized properly: File: " + rf.Location + rf.File + " was not properly moved. Additional info: " + err.Error())
        }
    }
    return nil
}

//Installs a Rice by symlinking the files into the specified dirs
func (rice Rice) InstallRice() (err error) {
	riceDir := rdbDir + rice.Program + "/" + rice.Name + "/"
	progDir := expandDir(rice.Root)

    if err:= os.Chdir(riceDir); err != nil {
        return errors.New("Error, this rice does not exist")
    }

    for _, rf := range rice.Files {
        if !exists(progDir + rf.Location){
            os.MkdirAll(progDir + rf.Location, 0755)
        }

        if err = os.Symlink(riceDir + rf.Location + rf.File, progDir + rf.Location + rf.File); err != nil {
            return errors.New("Error, this rice was not symlinked properly: File: " + rf.Location + rf.File + " was not properly symlinked. Additional info: " + err.Error())
        }
    }

    return nil
}

package ricecore

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

//A rice struct represents a rice with its files and other basic info
type Rice struct {
	Name    string      `json:"name"`
	Program string      `json:"prog"`
	Root    string      `json:"root"`
	Files   []*RiceFile `json:"files"`
}

//A ricefile struct is a file within a rice. It contains the file name and
//its respective location
type RiceFile struct {
	Location string `json:"location"`
	File     string `json:"file"`
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
		os.MkdirAll(riceDir, 0775)
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
	defer jsonFile.Close()

	rice = new(Rice)

	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(&rice); err != nil {
		return nil, err
	}

	return rice, nil
}

//Returns the currently active rice for a program
func GetActiveRice(prog string) (rice *Rice, err error) {
	progDir := rdbDir + prog + "/"

	s, err := ioutil.ReadFile(progDir + ".active")
	if err != nil {
		return nil, errors.New("Error, the .active file does not exist for this program.")
	}

	rname := string(s)
	return GetRice(rname, prog)
}

//Deactivates the currently active rice for a program
func DeactivateCurrentRice(prog string) (err error) {
	crice, err := GetActiveRice(prog)
	if err != nil {
		return err
	}

	if err = crice.deactivate(); err != nil {
		return err
	}

	return nil
}

//Initializes a created local rice, extracting the files from the directory to
//the rdb dir and symlinking them back
func (rice Rice) LocalInit() (err error) {
	riceDir := rdbDir + rice.Program + "/" + rice.Name + "/"
	progDir := expandDir(rice.Root)

	for _, rf := range rice.Files {
		if !exists(riceDir + rf.Location) {
			os.MkdirAll(riceDir+rf.Location, 0755)
		}

		if err = os.Rename(progDir+rf.Location+rf.File, riceDir+rf.Location+rf.File); err != nil {
			return errors.New("Error, this rice was not initialized properly: File: " + rf.Location + rf.File + " was not properly moved. Additional info: " + err.Error())
		}
	}
	return nil
}

//Activates a Rice by symlinking the files into the specified dirs
func (rice Rice) Activate() (err error) {
	riceDir := rdbDir + rice.Program + "/" + rice.Name + "/"
	progDir := expandDir(rice.Root)

	for _, rf := range rice.Files {
		if !exists(progDir + rf.Location) {
			os.MkdirAll(progDir+rf.Location, 0755)
		}

		if err = os.Symlink(riceDir+rf.Location+rf.File, progDir+rf.Location+rf.File); err != nil {
			return errors.New("Error, this rice was not symlinked properly: File: " + rf.Location + rf.File + " was not properly symlinked. Additional info: " + err.Error())
		}
	}
	activeRice := []byte(rice.Name)
	if err = ioutil.WriteFile(rdbDir+rice.Program+"/.active", activeRice, 0755); err != nil {
		return errors.New("Error, the .active file could not be created properly. Additional info: " + err.Error())
	}

	return nil
}

//Deactivates a rice by deleting all specified symlinks for that rice.
//Not publicly available because it would be another dumb thing to
//screw up. DeactivateCurrentRice should be used publicly
func (rice Rice) deactivate() (err error) {
	activeFile := rdbDir + rice.Program + "/.active"
	progDir := expandDir(rice.Root)

	for _, rf := range rice.Files {
		if !exists(progDir + rf.Location + rf.File) {
			return errors.New("Error, a specified file does not seem to exist")
		}

		if err = os.Remove(progDir + rf.Location + rf.File); err != nil {
			return errors.New("Error, this rice was not deactivated properly: File: " + rf.Location + rf.File + " was not properly removed. Additional info: " + err.Error())
		}
	}

	if exists(activeFile) {
		if err = os.Remove(activeFile); err != nil {
			return errors.New("Error, could not remove .active file. Additional info: " + err.Error())
		}
	}

	return nil
}

//Swaps in a rice by deactivating the currently active rice
//and activating the given rice
func (rice Rice) Swap() (err error) {
	if err := DeactivateCurrentRice(rice.Program); err != nil {
		return err
	}

	if err = rice.Activate(); err != nil {
		return err
	}

	return nil
}

//Uninstalls a rice for a program, deleting the symlinks
//and moving the files back to their original locations
func (rice Rice) Uninstall() (err error) {
	//Deactivates any active rice
	DeactivateCurrentRice(rice.Program)

	riceDir := rdbDir + rice.Program + "/" + rice.Name + "/"
	progDir := expandDir(rice.Root)

	for _, rf := range rice.Files {
		if !exists(progDir + rf.Location) {
			os.MkdirAll(progDir+rf.Location, 0755)
		}

		//Do better safety here, possibly use copys(these are more painful)
		if err = os.Rename(riceDir+rf.Location+rf.File, progDir+rf.Location+rf.File); err != nil {
			return errors.New("Error, this rice was not uninstalled properly: File: " + rf.Location + rf.File + " was not properly moved. Additional info: " + err.Error())
		}
	}

	if err = os.RemoveAll(riceDir); err != nil {
		return errors.New("Error, the remaining files could not be removed. Additional info: " + err.Error())
	}

	return nil
}

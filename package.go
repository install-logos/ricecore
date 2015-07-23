package ricecore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//A package object represents a package object which will be used for querys
//and uploads.
type Package struct {
	Name       string `json:"name"`
	Desc       string `json:"description"`
	Author     string `json:"author"`
	Screenshot string `json:"cover"`
	Program    string `json:"program"`
	URL        string `json:"upstream"`
}

//The result of a query - used for unmarshalling results
type QueryResult []*Package

//Querys the backend with the given keyword
//and returns an array of packages as a result
func QueryPackages(keyword string) (qRes *QueryResult, err error) {
	resp, err := http.Get("http://rice.kagayaite.com/query?q=" + keyword)
	if err != nil {
		return nil, errors.New("Error, could not query the RiceBE Database. Additional info: " + err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	qRes = new(QueryResult)

	if err = json.Unmarshal(data, &qRes); err != nil {
		return nil, errors.New("Could not convert response to a JSON. Additional info: " + err.Error())
	}

	return qRes, nil
}

//Uploads a package to the ricebe server
func (pack Package) Upload() (err error) {
	resp, err := http.Post("http://rice.kagayaite.com/upload?upstream="+pack.URL, "text", nil)
	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(response)
	return nil
}

//Loads a package.json file for an existing rice and returns the struct
func GetPackage(rice *Rice) (pack *Package, err error) {
	riceDir := rdbDir + rice.Program + "/" + rice.Name + "/"

	jsonFile, err := os.Open(riceDir + "package.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	pack = new(Package)

	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(&pack); err != nil {
		return nil, err
	}

	return pack, nil
}

//Creates a package from a rice and other data and stores the information as a json
func CreatePackage(rice *Rice, desc string, author string, screenshot string, url string) (pack *Package, err error) {
	pack = &Package{Name: rice.Name, Desc: desc, Author: author, Screenshot: screenshot, Program: rice.Program, URL: url}

	jsonData, err := json.MarshalIndent(pack, "", "  ")
	if err != nil {
		return nil, err
	}

	riceDir := rdbDir + rice.Program + "/" + rice.Name + "/"
	if !exists(riceDir) {
		return nil, errors.New("Error, this rice does not appear to be installed locally. Please install it and try again.")
	}

	jsonFile, err := os.Create(riceDir + "package.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	return pack, nil
}

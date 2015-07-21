package ricecore

import (
    "net/http"
    "errors"
    "io/ioutil"
    "encoding/json"
)

//A package object represents a package object which will be used for querys
//and uploads. Essentially, packages are stored online, rices are stored locally
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

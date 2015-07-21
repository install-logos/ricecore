package ricecore

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

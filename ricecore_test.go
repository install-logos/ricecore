package ricecore

import (
    "os"
    "io/ioutil"
    "testing"
)

func TestCreate(t *testing.T) {
    InitCore()

    //Setup
    os.Mkdir(homeDir + "/test/", 0755)
	testFileCont := []byte("le test")
    if err := ioutil.WriteFile(homeDir +"/test/test1", testFileCont, 0755); err != nil {
        t.Error(err.Error())
    }

    var files []*RiceFile
    files = make([]*RiceFile, 1)
    f1 := &RiceFile{Location:"./", File:"test1"}
    files[0] = f1
    _, err := CreateRice("test", "test-prog", "~/test/", files)
    if err != nil {
        t.Error(err.Error())
    }
    r, err := GetRice("test", "test-prog")
    if err != nil {
        t.Error(err.Error())
    }
    if err = r.LocalInit(); err != nil {
        t.Error(err.Error())
    }
}

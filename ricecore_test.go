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

func TestActivate(t *testing.T) {
    r, err := GetRice("test", "test-prog")
    if err != nil {
        t.Error(err.Error())
    }

    if err = r.Activate(); err != nil {
        t.Error(err.Error())
    }
}

func TestGetActiveRice(t *testing.T) {
    _, err := GetActiveRice("test-prog")
    if err != nil {
        t.Error(err.Error())
    }
}

func TestDeactivate(t *testing.T) {
    if err := DeactivateCurrentRice("test-prog"); err != nil {
        t.Error(err.Error())
    }
}

func TestUninstall(t *testing.T) {
    r, err := GetRice("test", "test-prog")
    if err != nil {
        t.Error(err.Error())
    }

    if err = r.Uninstall(); err != nil {
        t.Error(err.Error())
    }
}

func TestQuery(t *testing.T) {
    _, err := QueryPackages("logos")
    if err != nil {
        t.Error(err.Error())
    }
    //TODO: Validate results
}

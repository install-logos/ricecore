package ricecore

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	InitCore()

	//Setup
	os.Mkdir(homeDir+"/test/", 0755)
	testFileCont := []byte("le test")
	if err := ioutil.WriteFile(homeDir+"/test/test1", testFileCont, 0755); err != nil {
		t.Error(err.Error())
	}

	var files []*RiceFile
	files = make([]*RiceFile, 1)
	f1 := &RiceFile{Location: "./", File: "test1"}
	files[0] = f1
	_, err := CreateRice("test", "test-prog", "~/test/", files)
	if err != nil {
		t.Error(err.Error())
	}
	r, err := GetRice("test", "test-prog")
	if err != nil {
		t.Error(err.Error())
	}
	if err = r.FirstInit(); err != nil {
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

func TestCreate2(t *testing.T) {
	os.Mkdir(homeDir+"/test2/", 0755)
	testFileCont := []byte("le second test")
	if err := ioutil.WriteFile(homeDir+"/test2/test2", testFileCont, 0755); err != nil {
		t.Error(err.Error())
	}

	var files []*RiceFile
	files = make([]*RiceFile, 1)
	f1 := &RiceFile{Location: "./", File: "test2"}
	files[0] = f1
	_, err := CreateRice("test2", "test-prog", "~/test/", files)
	if err != nil {
		t.Error(err.Error())
	}
	r, err := GetRice("test2", "test-prog")
	if err != nil {
		t.Error(err.Error())
	}
	if err = r.LocalInit("~/test2/"); err != nil {
		t.Error(err.Error())
	}
}

func TestSwap(t *testing.T) {
	r, err := GetRice("test2", "test-prog")
	if err != nil {
		t.Error(err.Error())
	}

	r.Swap()
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

func TestDownload(t *testing.T) {
	res, err := QueryPackages("test")
	if err != nil {
		t.Error(err.Error())
	}
	//TODO: Validate results
	r := *res
	_, err = r[0].Download()
	if err != nil {
		t.Error(err.Error())
	}
}

package utils_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/xwi88/kit4go/utils"
)

var tmpDir string

func setup() {
	tmpDir = os.TempDir()
	tmpDir = fmt.Sprintf("%v/%v", tmpDir, "xwi88_file_util_test")
	if err := os.MkdirAll(tmpDir, os.FileMode(0777)); err != nil {
		log.Fatalf("Mkdir err:%v", err.Error())
	}
}

func teardown() {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatal(err)
	}
}

func TestCopyFile(t *testing.T) {
	setup()
	defer teardown()

	dstFileName := "1.txt.bak"
	srcFileName := "1.txt"

	srcFile := fmt.Sprintf("%v/%v", tmpDir, srcFileName)
	dstFile := fmt.Sprintf("%v/%v", tmpDir, dstFileName)

	f, err := os.Create(srcFile)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println(f.Name())
	defer f.Close()

	if w, err := utils.CopyFile(dstFile, srcFile); err != nil {
		log.Fatalf("CopyFile err:%v", err)
	} else {
		log.Printf("CopyFile %v", w)
	}
}

func TestGetFileInfo(t *testing.T) {
	setup()
	defer teardown()

	srcFileName := "1.txt"
	srcFile := fmt.Sprintf("%v/%v", tmpDir, srcFileName)
	f, err := os.Create(srcFile)
	if err != nil {
		t.Errorf("%v", err)
	}
	log.Println(f.Name())
	defer f.Close()

	fi := utils.GetFileInfo(srcFile)
	t.Logf("GetFileInfo:%v", fi)
}

func TestIsDir(t *testing.T) {
	setup()
	defer teardown()

	if isDir := utils.IsDir(tmpDir); isDir {
		t.Logf("%v is dir:%v", tmpDir, isDir)
	} else {
		t.Errorf("%v is not dir:%v", tmpDir, isDir)
	}

	newTmpDir := tmpDir + "2012"
	if isDir := utils.IsDir(newTmpDir); !isDir {
		t.Logf("%v is dir:%v", newTmpDir, isDir)
	} else {
		t.Errorf("%v is dir:%v", newTmpDir, isDir)
	}
}

func TestIsExist(t *testing.T) {
	setup()
	defer teardown()

	exist := utils.IsExist(tmpDir)
	if exist {
		t.Logf("%v exist %v", tmpDir, exist)
	} else {
		t.Errorf("%v exist %v", tmpDir, exist)
	}
}

func TestIsFile(t *testing.T) {
	setup()
	defer teardown()

	isFile := utils.IsFile(tmpDir)
	if isFile {
		t.Errorf("%v is file %v", tmpDir, isFile)
	} else {
		t.Logf("%v is file %v", tmpDir, isFile)
	}
}

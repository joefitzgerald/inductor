package cpy

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type fileCopier struct {
	srcDir string
	outDir string
}

// New create a new cpy instance
func New() Copier {
	return &fileCopier{}
}

func (cp *fileCopier) Copy(srcDir, outDir string) error {
	if err := cp.initCopyDirs(srcDir, outDir); err != nil {
		return err
	}
	return filepath.Walk(srcDir, cp.walkFile)
}

func (cp *fileCopier) initCopyDirs(srcDir, outDir string) error {
	sfi, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	if !sfi.IsDir() {
		return errors.New("Expected Copy source to be a directory")
	}

	tfi, err := os.Stat(outDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(outDir, 0777)
		} else {
			return err
		}
	} else {
		if !tfi.IsDir() {
			return errors.New("Expected Copy destination to be a directory")
		}
	}
	cp.srcDir = srcDir
	cp.outDir = outDir
	return nil
}

func (cp *fileCopier) walkFile(sf string, sfi os.FileInfo, err error) error {
	if sfi.IsDir() {
		if strings.HasPrefix(sfi.Name(), ".") {
			// don't copy hidden dirs
			return filepath.SkipDir
		}
		return nil
	}

	// don't copy template files
	if strings.HasPrefix(sfi.Name(), ".") || filepath.Ext(sf) == ".template" || filepath.Ext(sf) == ".partial" {
		return nil
	}

	// we have a file, calculate its relative destination location
	rel := strings.TrimPrefix(sf, cp.srcDir)
	df := filepath.Join(cp.outDir, rel)
	dir := filepath.Dir(df)

	if err = cp.mkdir(dir); err != nil {
		return err
	}
	if err = cp.copyFile(sf, df); err != nil {
		return err
	}
	return nil
}

func (cp *fileCopier) mkdir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func (cp *fileCopier) copyFile(source, target string) (err error) {
	//fmt.Println(fmt.Sprintf("%s => %s", source, target))
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := sf.Close(); cerr != nil {
			err = cerr
		}
	}()
	tf, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := tf.Close(); cerr != nil {
			err = cerr
		}
	}()
	_, err = io.Copy(tf, sf)
	if err != nil {
		return err
	}
	err = tf.Sync()
	return err
}

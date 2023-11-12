package uzip

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Compress(dst io.Writer, root string) error {
	//prepare zip
	zipWriter := zip.NewWriter(dst)
	defer zipWriter.Close()

	addToZip := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		//directory
		if info.IsDir() {
			_, err = zipWriter.Create(relativePath + "/")
			return err
		}

		//file
		fileWriter, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(fileWriter, file)
		return err
	}

	//compress files into the zip buffer
	return filepath.WalkDir(root, addToZip)
}

func Extract(desPath string, zipHandle *zip.Reader) error {
	for _, fileHandle := range zipHandle.File {
		path := filepath.Join(desPath, fileHandle.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
			return err
		}

		if !fileHandle.FileInfo().IsDir() {
			srcFile, err := fileHandle.Open()
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

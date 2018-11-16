package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"archive/zip"
	"log"
)

func UnzipFiles(src string, dest string) []string{
	files, err := Unzip(src, dest)
	if err != nil {
		log.Fatal(err)
	}
	
	//fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))
	return files
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {
	
	var filenames []string
	
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()
	
	for _, f := range r.File {
		
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()
		
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		
		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}
		
		filenames = append(filenames, fpath)
		
		if f.FileInfo().IsDir() {
			
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			
		} else {
			
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}
			
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}
			
			_, err = io.Copy(outFile, rc)
			
			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			
			if err != nil {
				return filenames, err
			}
			
		}
	}
	return filenames, nil
}

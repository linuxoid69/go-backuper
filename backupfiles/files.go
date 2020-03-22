package backupfiles

import (
	"archive/zip"
	"github.com/linuxoid69/go-backuper/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ZipFiles - create zip archive
// files - files for adding to an archive
// zipFileName - a name of zip file
func ZipFiles(files []string, zipFileName string) error {

	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

// AddFileToZip - adding files to zip archive
func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// Find - find files in directory return files and directories
// dirname - path type string.
// outputType - (files, dirs) if arg is empty will be output both file and dir.
func Find(dirname string, outputType string) (files []string, err error) {

	err = filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})
	if err != nil {
		return []string{}, err
	}
	filesOutput := []string{}

	switch output := outputType; output {
	case "dirs":
		for _, i := range files {
			if isDirectory(i) {
				filesOutput = append(filesOutput, i)
			}
		}
		return filesOutput, nil
	case "files":
		for _, i := range files {
			if !isDirectory(i) {
				filesOutput = append(filesOutput, i)
			}
		}
		return filesOutput, nil
	default:
		return files, nil
	}

}

// IsDirectory - detect is directory or not
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	return info.IsDir()
}

// CreateArchive - create archive
func CreateArchive(files []string, nameArchive string) (err error) {
	totalFiles := []string{}
	for _, i := range files {
		out, err := Find(i, "files")

		if err != nil {
			log.Printf("Error %v", err)
		}

		for _, f := range out {
			totalFiles = append(totalFiles, f)
		}
	}
	err = ZipFiles(totalFiles, nameArchive)
	if err != nil {
		return err
	}
	return nil

}

// CreateArchives - create archives
func CreateArchives(cfg *config.Config) {
	t := strings.Replace(time.Now().Format("02-01-2006"), "-", "_", 2)
	for _, p := range cfg.BackupSourceFilesPath["projects"][0] {
		pathZipFile := cfg.BackupStoragePath + "/" + p["project_name"] + "_" + t + "_" + cfg.NameZipFile

		err := CreateArchive(strings.Split(p["dirs"], ","), pathZipFile)

		if err != nil {
			log.Printf("Archive %+v didn't successful created\n", pathZipFile)
		} else {
			log.Printf("Archive %+v was successful created\n", pathZipFile)

		}
	}
}

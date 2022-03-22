package backupfiles

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/linuxoid69/go-backuper/config"
	"github.com/linuxoid69/go-backuper/internal/helpers"
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

	if _, err = io.Copy(writer, fileToZip); err != nil {
		return err
	}

	return nil
}

// CreateArchive - create archive.
func CreateArchive(files []string, nameArchive string) (err error) {
	allFiles := []string{}
	rDays := 0

	for _, i := range files {
		out, err := helpers.Find(i, "files", rDays)

		if err != nil {
			logrus.Errorf("Error %w", err)
		}

		for _, f := range out {
			allFiles = append(allFiles, f)
		}
	}

	if err = ZipFiles(allFiles, nameArchive); err != nil {
		return err
	}

	return nil
}

// CreateArchives - create archives
func CreateArchives(cfg *config.Config) (backupFiles []string, err error) {
	dateString := strings.Replace(time.Now().Format("02-01-2006"), "-", "_", 2)
	for _, p := range cfg.BackupSourceFilesPath["projects"][0] {
		pathZipFile := fmt.Sprintf("%s/%s_%s_%s", cfg.BackupStoragePath,
			p["project_name"],
			dateString,
			cfg.NameZipFile)

		if err := CreateArchive(strings.Split(p["dirs"], ","), pathZipFile); err != nil {
			logrus.Errorf("Archive %+v didn't successful created\n", pathZipFile)

			return backupFiles, err
		} else {
			logrus.Infof("Archive %+v was successful created\n", pathZipFile)
		}

		backupFiles = append(backupFiles, pathZipFile)
	}

	return backupFiles, nil
}

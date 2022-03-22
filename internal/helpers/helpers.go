package helpers

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Find - find files in directory return files and directories
/*
   dirname - path type string.
   outputType - type dir or file.
   rDays - retention days.
   outputType - (files, dirs) if arg is empty will be output both file and dir.
*/
func Find(dirname string, outputType string, rDays int) (result []string, err error) {

	err = filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {

			if rDays > 0 {

				if err != nil {
					return err
				}

				if dirname != path {

					oldTime := info.ModTime().Unix()
					nowTime := time.Now().Unix()
					retentionDays := rDays * 86400

					if int(nowTime-oldTime) >= retentionDays {
						result = append(result, path)
					}
				}
			} else {
				if err != nil {
					return err
				}
				result = append(result, path)
			}

			return nil
		})

	if err != nil {
		return []string{}, err
	}

	filesOutput := []string{}

	switch output := outputType; output {
	case "dirs":
		for _, i := range result {
			if isDirectory(i) {
				filesOutput = append(filesOutput, i)
			}
		}
		return filesOutput, nil
	case "files":
		for _, i := range result {
			if !isDirectory(i) {
				filesOutput = append(filesOutput, i)
			}
		}
		return filesOutput, nil
	default:
		return result, nil
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

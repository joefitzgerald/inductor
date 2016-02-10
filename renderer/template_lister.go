package renderer

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FilterFilesToRender filters the list of files that should be rendered
// according to inductor rules for template precedence.
func FilterFilesToRender(files []string, osName string) []string {
	fileMap := make(map[string]string)
	for _, f := range files {
		// get the filename without the OS name (if any)
		filename := filepath.Base(f)
		key := strings.Replace(filename, "-"+osName, "", 1)
		fmt.Println(key)

		// add an entry, but ensure we keep the OS specific entries
		// if it produced the same key but is a longer path, it must be OS specific
		curFile, _ := fileMap[key]
		if len(f) > len(curFile) {
			curFile = f
		}
		fileMap[key] = curFile
	}

	//var distinctFiles []string
	distinctFiles := make([]string, 0, len(fileMap))
	for _, v := range fileMap {
		distinctFiles = append(distinctFiles, v)
	}
	return distinctFiles
}

// ListFiles returns all files in the current working directory that start with
// the specific base file name.
func ListFiles(baseDir string, filenameBase string) ([]string, error) {
	return filepath.Glob(filepath.Join(baseDir, filenameBase) + "*")
}

package renderer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FilterFilesToRender filters the list of files that should be rendered
// according to inductor rules for template precedence.
func FilterFilesToRender(files []string, osName string) []string {
	fileMap := make(map[string]string)
	for _, file := range files {
		// get the filename without the directory
		filename := filepath.Base(file)

		// if this file is OS specific, parse out that OS name
		re := regexp.MustCompile("^\\w+-(\\w+)\\.")
		matches := re.FindStringSubmatch(filename)
		if matches != nil && len(matches) > 1 {
			fileOsName := matches[1]

			// skip any files which have an OS in their name that isn't ours
			if !strings.EqualFold(fileOsName, osName) {
				continue
			}
		}

		// file key shouldn't contain the OS
		key := strings.Replace(filename, "-"+osName, "", 1)

		// add an entry, but ensure we keep the OS specific entries
		// if it produced the same key but is a longer path, it must be OS specific
		curFile, _ := fileMap[key]
		if len(file) > len(curFile) {
			curFile = file
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
func ListFiles(baseDir string, filenameBase string) []string {
	var files []string
	entries, _ := filepath.Glob(filepath.Join(baseDir, filenameBase) + "*")
	for _, e := range entries {
		if isFile(e) {
			files = append(files, e)
		}
	}
	return files
}

func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

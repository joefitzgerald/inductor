package cpy

// Copier will recursively copy all the file from the source directory
// to the given output directory
type Copier interface {
	Copy(srcDir, outDir string) error
}

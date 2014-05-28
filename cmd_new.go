package main

import (
	"fmt"
	"go/build"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var newCmd = &Command{
	UsageLine: "new [app_path]",
	Short:     "create a new armadillo application",
	Long: `'new' creates a skeleton Armadillo app in the src of folder of your GOPATH

		It uses the last portion of the path as the new app's name

		For example:

			armadillo new myAppName	#creates a new skeleton project in ~GOPATH/src/myAppName with the name 'myAppName'

			armadillo new repp/projects/ApplicationName	#creates a new skeleton project in ~GOPATH/src/repp/projects/ApplicationName with the name 'ApplicationName'`,
}

var (
	srcRoot      string
	appPath      string
	appName      string
	skeletonPath string
	templateData map[string]interface{}
)

func init() {
	newCmd.Run = createNewApp
}

func createNewApp(args []string) {
	if len(args) == 0 {
		errorf("No app path provided.\nRun 'armadillo help new' for usage.\n")
	}
	if len(args) > 1 {
		errorf("Too many arguments provided.\nRun 'armadillo help new' for usage.\n")
	}

	setSrcRoot()
	setAppPaths(args)
	createSkeletonApp()
}

func setSrcRoot() {
	// lookup go path
	gopath := build.Default.GOPATH
	if gopath == "" {
		errorf("Abort: GOPATH environment variable is not set. " +
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
	}

	// set go src path
	srcRoot = filepath.Join(filepath.SplitList(gopath)[0], "src")
}

func setAppPaths(args []string) {
	var err error
	importPath := args[0]
	if filepath.IsAbs(importPath) {
		errorf("Abort: '%s' looks like a directory.  Please provide a Go import path instead.", importPath)
	}

	_, err = build.Import(importPath, "", build.FindOnly)
	if err == nil {
		errorf("Abort: Import path %s already exists.\n", importPath)
	}

	armadilloPkg, err2 := build.Import(IMPORT_PATH, "", build.FindOnly)
	if err2 != nil {
		errorf("Abort: Could not find Armadillo source code: %s\n", err)
	}

	appPath = filepath.Join(srcRoot, filepath.FromSlash(importPath))
	appName = filepath.Base(appPath)
	skeletonPath = filepath.Join(armadilloPkg.Dir, "skeleton")
	templateData = map[string]interface{}{
		"AppName": appName,
	}

}

func createSkeletonApp() {
	var err error
	err = os.MkdirAll(appPath, 0777)
	if err != nil {
		errorf(err.Error())
		return
	}
	copyDir(skeletonPath)
}

func copyDir(srcDir string) error {
	var fullSrcDir string
	// Handle symlinked directories.
	f, err := os.Lstat(srcDir)
	if err == nil && f.Mode()&os.ModeSymlink == os.ModeSymlink {
		fullSrcDir, err = os.Readlink(srcDir)
		if err != nil {
			panic(err)
		}
	} else {
		fullSrcDir = srcDir
	}

	return filepath.Walk(fullSrcDir, walkDir)
}

func walkDir(srcPath string, info os.FileInfo, err error) error {
	// Get the relative path from the source base, and the corresponding path in
	// the dest directory.
	relSrcPath := strings.TrimLeft(srcPath[len(skeletonPath):], string(os.PathSeparator))
	destPath := path.Join(appPath, relSrcPath)

	// Skip dot files and dot directories.
	if strings.HasPrefix(relSrcPath, ".") {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// Create a subdirectory if necessary.
	if info.IsDir() {
		err := os.MkdirAll(path.Join(appPath, relSrcPath), 0777)
		if !os.IsExist(err) {
			panicOnError(err, "Failed to create directory")
		}
		return nil
	}

	// If this file ends in ".template", render it as a template.
	if strings.HasSuffix(relSrcPath, ".template") {
		copyTemplateFile(destPath[:len(destPath)-len(".template")], srcPath, templateData)
		return nil
	}

	// Copy files over
	copyFile(destPath, srcPath)
	return nil
}

func copyFile(destFilename, srcFilename string) {
	destFile, err := os.Create(destFilename)
	panicOnError(err, "Failed to create file "+destFilename)

	srcFile, err := os.Open(srcFilename)
	panicOnError(err, "Failed to open file "+srcFilename)

	_, err = io.Copy(destFile, srcFile)
	panicOnError(err,
		fmt.Sprintf("Failed to copy data from %s to %s", srcFile.Name(), destFile.Name()))

	err = destFile.Close()
	panicOnError(err, "Failed to close file "+destFile.Name())

	err = srcFile.Close()
	panicOnError(err, "Failed to close file "+srcFile.Name())
}

func copyTemplateFile(destFilename, srcFilename string, data map[string]interface{}) {
	tmpl, err := template.ParseFiles(srcFilename)
	panicOnError(err, "Failed to parse template "+srcFilename)

	f, err := os.Create(destFilename)
	panicOnError(err, "Failed to create "+destFilename)

	err = tmpl.Execute(f, data)
	panicOnError(err, "Failed to render template "+srcFilename)

	err = f.Close()
	panicOnError(err, "Failed to close "+f.Name())
}

package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// ------------------------------------------------------------------
// kalpavrikshaConfig : config structure for the tool
type kalpavrikshaConfig struct {
	StorageConfig // Storage config

	// Base config for what to do
	NumberOfDirs  int64 // Number of directories to be created
	DirDepth      int64 // Number of sub-directories to be created inside each directory
	NumberOfFiles int64 // Number of files to be created
	FileSize      int64 // Size of each file to be created
	Parallelism   int   // Number of threads to run in parallel

	InputTypeStr string     // Type of input in string : Zero / Rand / File
	InputType    SourceType // Type of input : Zero / Rand / File

	SourceFilePath string // In case of input is coming from a file, path to that file
	Tier           string // blob tier to set on upload

	Delete  bool // Delete the previously generated data on given path
	SetTier bool // Change Tier of previously generated data on given path

	CreateStub bool // Create directory stub files on the given path
	DeleteStub bool // Delete directory stub files on the given path
}

type Kalpavriksha struct {
	// Storage accoutn related config
	storage Storage // Storage client

	// Worker related internal objects
	jobs      chan workItem  // Channel holding jobs to be performed
	results   chan workItem  // Channel holding jobs which are done
	wgWorkers sync.WaitGroup // Wait group for all workers

	// Data source provider
	dataSrc dataSource // Source of data for input
}

// global variable holding all of the config
var config kalpavrikshaConfig
var kalpavriksha Kalpavriksha

// Methods to operate on config
func sanitizeConfig() error {
	err := config.InputType.Parse(config.InputTypeStr)
	if err != nil {
		return err
	}

	if config.InputType == ESourceType.FILE() {
		if _, err := os.Stat(config.SourceFilePath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("source file %s does not exists", config.SourceFilePath)
		}
	}

	config.FileSize = config.FileSize * 1024 * 1024
	readStorageParams()
	return nil
}

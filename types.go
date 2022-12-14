package main

import (
	"reflect"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/JeffreyRichter/enum/enum"
)

// ------------------------------------------------------------------
// kalpavrikshaConfig : config structure for the tool
type kalpavrikshaConfig struct {
	NumberOfDirs  int64 // Number of directories to be created
	NumberOfFiles int64 // Number of files to be created
	FileSize      int64 // Size of each file to be created
	Parallelism   int   // Number of threads to run in parallel

	InputTypeStr string     // Type of input in string : Zero / Rand / File
	InputType    SourceType // Type of input : Zero / Rand / File

	SourceFilePath  string // In case of input is coming from a file, path to that file
	DestinationPath string // Provide destination path (post container)

	StorageAccountName      string // Name of the destination storage account
	StorageAccountKey       string // Key of the destination storage account
	StorageEndPoint         string // Type of storage account blob/dfs
	StorageAccountContainer string // Destination container in the storage account

	StorageClient *container.Client // Client to hold storage connection

	jobs      chan workItem  // Channel holding jobs to be performed
	results   chan workItem  // Channel holding jobs which are done
	wgWorkers sync.WaitGroup // Wait group for all workers

	src dataSource // Source of data for input
}

// global variable holding all of the config
var config kalpavrikshaConfig

// ------------------------------------------------------------------
// Input Source Type
type SourceType int

var ESourceType = SourceType(0).INVALID_SRC()

func (SourceType) INVALID_SRC() SourceType {
	return SourceType(0)
}

func (SourceType) ZERO() SourceType {
	return SourceType(1)
}

func (SourceType) RANDOM() SourceType {
	return SourceType(2)
}

func (SourceType) FILE() SourceType {
	return SourceType(3)
}

func (f SourceType) String() string {
	return enum.StringInt(f, reflect.TypeOf(f))
}

func (a *SourceType) Parse(s string) error {
	enumVal, err := enum.ParseInt(reflect.TypeOf(a), s, true, false)
	if enumVal != nil {
		*a = enumVal.(SourceType)
	}

	return err
}

// ------------------------------------------------------------------
// Job status Type
type JobStatusType int

var EJobStatusType = JobStatusType(0).INVALID_JOBTYPE()

func (JobStatusType) INVALID_JOBTYPE() JobStatusType {
	return JobStatusType(0)
}

func (JobStatusType) WAIT() JobStatusType {
	return JobStatusType(1)
}

func (JobStatusType) INPROGRESS() JobStatusType {
	return JobStatusType(2)
}

func (JobStatusType) SUCCESS() JobStatusType {
	return JobStatusType(3)
}

func (JobStatusType) FAILED() JobStatusType {
	return JobStatusType(4)
}

func (f JobStatusType) String() string {
	return enum.StringInt(f, reflect.TypeOf(f))
}

func (a *JobStatusType) Parse(s string) error {
	enumVal, err := enum.ParseInt(reflect.TypeOf(a), s, true, false)
	if enumVal != nil {
		*a = enumVal.(JobStatusType)
	}

	return err
}

// ------------------------------------------------------------------
// Object type
type ObjectType int

var EObjectType = ObjectType(0).INVALID_OBJECTTYPE()

func (ObjectType) INVALID_OBJECTTYPE() ObjectType {
	return ObjectType(0)
}

func (ObjectType) FILE() ObjectType {
	return ObjectType(1)
}

func (ObjectType) DIR() ObjectType {
	return ObjectType(2)
}

func (f ObjectType) String() string {
	return enum.StringInt(f, reflect.TypeOf(f))
}

func (a *ObjectType) Parse(s string) error {
	enumVal, err := enum.ParseInt(reflect.TypeOf(a), s, true, false)
	if enumVal != nil {
		*a = enumVal.(ObjectType)
	}

	return err
}

// ------------------------------------------------------------------
// Azure storage related env variables
const (
	EnvAzStorageAccount          = "AZURE_STORAGE_ACCOUNT"
	EnvAzStorageAccessKey        = "AZURE_STORAGE_ACCESS_KEY"
	EnvAzStorageAccountContainer = "AZURE_STORAGE_ACCOUNT_CONTAINER"
)

// ------------------------------------------------------------------

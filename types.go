package main

import (
	"reflect"

	"github.com/JeffreyRichter/enum/enum"
)

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
	EnvAzStorageSAS              = "AZURE_STORAGE_SAS_TOKEN"
	EnvAzStorageAccountContainer = "AZURE_STORAGE_ACCOUNT_CONTAINER"
)

// ------------------------------------------------------------------

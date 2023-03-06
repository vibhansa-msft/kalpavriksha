package main

import (
	"fmt"
	golog "log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/JeffreyRichter/enum/enum"
)

const (
	responseStatusString = "RESPONSE Status:"
)

// ------------------------------------------------------------------
// Storage type
type StorageType int

var EStorageType = StorageType(0).INVALID_STORAGE()

func (StorageType) INVALID_STORAGE() StorageType {
	return StorageType(0)
}

func (StorageType) BLOB() StorageType {
	return StorageType(1)
}

func (StorageType) FILE() StorageType {
	return StorageType(2)
}

func (StorageType) DATALAKE() StorageType {
	return StorageType(3)
}

func (f StorageType) String() string {
	return enum.StringInt(f, reflect.TypeOf(f))
}

func (a *StorageType) Parse(s string) error {
	enumVal, err := enum.ParseInt(reflect.TypeOf(a), s, true, false)
	if enumVal != nil {
		*a = enumVal.(StorageType)
	}

	return err
}

// ------------------------------------------------------------------

type StorageConfig struct {
	StorageAccountName      string // Name of the destination storage account
	StorageAccountKey       string // Key of the destination storage account
	StorageAccountSAS       string // SAS Key of the destination storage account
	StorageEndPoint         string // Type of storage account blob/dfs
	StorageAccountContainer string // Destination container in the storage account

	AccountType     StorageType // Type of storage account Blob / File / Datalake
	DestinationPath string      // Provide destination path (post container)

	UpdateMD5 bool            // Set MD5SUM on upload
	BlobTier  blob.AccessTier // Set tier value on upload
}

type UploadOptions struct {
	Tier   *blob.AccessTier
	MD5Sum []byte
}

type DeleteOptions struct {
}

type Storage interface {
	Init() error
	TestConnection() error
	UploadData(name string, data []byte, o *UploadOptions) error
	Delete(name string, o *DeleteOptions) error
	SetTier(name string, tier blob.AccessTier) error
	CreateStub(name string) error
	ListBlobs(name string) *runtime.Pager[container.ListBlobsHierarchyResponse]
	GetProperties(name string) (blob.GetPropertiesResponse, error)
}

func setupLogging() {
	log.SetEvents(log.EventRequest, log.EventResponse)
	log.SetListener(func(cls log.Event, msg string) {
		switch cls {
		case log.EventRequest:
		case log.EventRetryPolicy:
		case log.EventLRO:
			// We do not want to log the request
			break

		case log.EventResponse:
			index := strings.Index(msg, responseStatusString)
			if index > 0 {
				index += len(responseStatusString) + 1
				respCode := msg[index : index+3]
				respCodeVal, _ := strconv.Atoi(respCode)
				if respCodeVal >= 400 {
					golog.Printf("Request failed with status code %v\n", respCode)
				}
			}
			break
		default:
			break
		}
	})
}

func readStorageParams() {
	config.BlobTier = blob.AccessTier(config.Tier)
	config.AccountType = EStorageType.BLOB()
	config.StorageAccountName = os.Getenv(EnvAzStorageAccount)
	config.StorageAccountKey = os.Getenv(EnvAzStorageAccessKey)
	config.StorageAccountSAS = os.Getenv(EnvAzStorageSAS)
	config.StorageAccountContainer = os.Getenv(EnvAzStorageAccountContainer)
}

func createStorage(t StorageType, c StorageConfig) (Storage, error) {
	// Setup logging for storage SDK
	setupLogging()

	var stobj Storage

	if t == EStorageType.BLOB() {
		stobj = &BlobStorage{StorageConfig: c}
	} else {
		return nil, fmt.Errorf("invalid storage type")
	}

	err := stobj.Init()
	if err != nil {
		return nil, err
	}

	err = stobj.TestConnection()
	if err != nil {
		return nil, err
	}

	return stobj, nil
}

package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

const (
	responseStatusString = "RESPONSE Status:"
)

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
					fmt.Println("Request failed with status code ", respCode)
				}
			}
		default:
			break
		}
	})
}

func setupStorageConnection() error {
	// Setup logging for storage SDK
	setupLogging()

	// Create credential object using storage account name and key
	cred, err := azblob.NewSharedKeyCredential(config.StorageAccountName, config.StorageAccountKey)
	if err != nil {
		return err
	}

	containerURL := fmt.Sprintf("https://%s.%s.core.windows.net/%s", config.StorageAccountName, config.StorageEndPoint, config.StorageAccountContainer)
	containerClient, err := container.NewClientWithSharedKeyCredential(containerURL, cred, nil)
	if err != nil {
		return err
	}

	config.StorageClient = containerClient
	return testConnectionByList()
}

func testConnectionByList() error {
	// Try to list the container and see if auth gets validated or not
	maxResults := int32(2)
	pager := config.StorageClient.NewListBlobsHierarchyPager("/", &container.ListBlobsHierarchyOptions{
		Include:    container.ListBlobsInclude{Metadata: true},
		MaxResults: &maxResults,
	})

	if pager == nil {
		return fmt.Errorf("failed to authenticate to storage account")
	}

	for pager.More() {
		_, err := pager.NextPage(context.TODO())
		if err != nil {
			return err
		}

		// We are able to get some blobs so it means connection is successful
		return nil
	}

	return nil
}

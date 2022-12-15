package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type BlobStorage struct {
	StorageConfig
	StorageClient *container.Client // Client to hold storage connection
}

func (bs *BlobStorage) Init() error {
	// Create credential object using storage account name and key
	cred, err := azblob.NewSharedKeyCredential(bs.StorageAccountName, bs.StorageAccountKey)
	if err != nil {
		return err
	}

	containerURL := fmt.Sprintf("https://%s.%s.core.windows.net/%s", bs.StorageAccountName, bs.StorageEndPoint, bs.StorageAccountContainer)
	bs.StorageClient, err = container.NewClientWithSharedKeyCredential(containerURL, cred, nil)
	if err != nil {
		return err
	}

	return nil
}

func (bs *BlobStorage) TestConnection() error {
	// Try to list the container and see if auth gets validated or not
	maxResults := int32(2)
	pager := bs.StorageClient.NewListBlobsHierarchyPager("/", &container.ListBlobsHierarchyOptions{
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

func (bs *BlobStorage) UploadData(name string, data []byte, o *UploadOptions) error {
	blockBlobClient := bs.StorageClient.NewBlockBlobClient(filepath.Join(bs.DestinationPath, name))

	opts := &azblob.UploadBufferOptions{}
	if o != nil {
		if o.MD5Sum != nil {
			opts.TransactionalContentMD5 = o.MD5Sum
		}

		if o.Tier != nil {
			opts.AccessTier = o.Tier
		}
	}

	_, err := blockBlobClient.UploadBuffer(context.TODO(), data, opts)

	return err
}

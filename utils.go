package main

import (
	"errors"
	"fmt"
	"os"
)

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

	readStorageParams()
	return nil
}

func readStorageParams() {
	config.StorageAccountName = os.Getenv(EnvAzStorageAccount)
	config.StorageAccountKey = os.Getenv(EnvAzStorageAccessKey)
	config.StorageAccountContainer = os.Getenv(EnvAzStorageAccountContainer)
}

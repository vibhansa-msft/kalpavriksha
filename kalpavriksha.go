package main

import (
	"flag"
	"fmt"
	//"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

func main() {
	fmt.Println("Kalpavrikhsa starting")

	// Parse the user config
	flag.Parse()

	// Sanitize the config
	if err := sanitizeConfig(); err != nil {
		// In case of any config failure, terminate
		fmt.Println("invalid config.", err.Error())
		flag.Usage()
		return
	}

	err := setupStorageConnection()
	if err != nil {
		fmt.Println("failed to connect to storaged.", err.Error())
		return
	}

	config.src, err = createDataSource()
	if err != nil {
		fmt.Println("failed to create data source.", err.Error())
		return
	}

	startWorkers()
	fmt.Println("Kalpavrikhsa completed")
}

func init() {
	config = kalpavrikshaConfig{}

	flag.Int64Var(&config.NumberOfDirs, "dirs", 1, "Number of directories to be created")
	flag.Int64Var(&config.NumberOfFiles, "files", 1, "Number of files to be created per directory")
	flag.Int64Var(&config.FileSize, "size", 1, "Size of each file to be created")
	flag.IntVar(&config.Parallelism, "concurrency", 64, "Number of threads to run in parllel")

	flag.StringVar(&config.InputTypeStr, "type", "random", "Type of source ZERO / RANDOM / FILE")

	flag.StringVar(&config.SourceFilePath, "src-file", "", "Source file to be used for data")
	flag.StringVar(&config.DestinationPath, "dst-path", "", "Destination path after the container where files will be created")

	flag.StringVar(&config.StorageEndPoint, "acct-type", "blob", "Stroage account type")
}

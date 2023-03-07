package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	file, err := os.OpenFile("kalpavriksha.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Kalpavrikhsa starting")

	// Parse the user config
	flag.Parse()

	// Sanitize the config
	if err = sanitizeConfig(); err != nil {
		// In case of any config failure, terminate
		fmt.Println("invalid config.", err.Error())
		flag.Usage()
		return
	}

	kalpavriksha.storage, err = createStorage(EStorageType.BLOB(), config.StorageConfig)
	if err != nil {
		fmt.Println("failed to connect to storage.", err.Error())
		return
	}

	kalpavriksha.dataSrc, err = createDataSource(config.InputType)
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
	flag.Int64Var(&config.DirDepth, "depth", 0, "Number of sub-directory inside which file will be created")
	flag.Int64Var(&config.NumberOfFiles, "files", 1, "Number of files to be created per directory")
	flag.Int64Var(&config.FileSize, "size", 1, "Size of each file to be created")
	flag.IntVar(&config.Parallelism, "concurrency", 64, "Number of threads to run in parllel")

	flag.StringVar(&config.InputTypeStr, "type", "random", "Type of source ZERO / RANDOM / FILE")

	flag.StringVar(&config.SourceFilePath, "src-file", "", "Source file to be used for data")
	flag.StringVar(&config.DestinationPath, "dst-path", "", "Destination path after the container where files will be created")

	flag.StringVar(&config.StorageEndPoint, "acct-type", "blob", "Stroage account type")

	flag.BoolVar(&config.UpdateMD5, "md5", false, "Set MD5 Sum on upload")
	flag.StringVar(&config.Tier, "tier", "none", "Tier to be set for each file")

	flag.BoolVar(&config.Delete, "delete", false, "Delete the data set instead of generation")
	flag.BoolVar(&config.SetTier, "set-tier", false, "Change the tier of previously generated dataset")

	flag.BoolVar(&config.CreateStub, "create-stub", false, "Create directory stub on the given path")
	flag.BoolVar(&config.DeleteStub, "delete-stub", false, "Delete directory stub on the given path")
}

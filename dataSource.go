package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// -------------------------------------------------------------------
type dataSource interface {
	Init(i interface{}) error
	GetData(size uint64) ([]byte, error)
}

//-------------------------------------------------------------------

type zeroDataSource struct {
}

func (zds *zeroDataSource) Init(i interface{}) error {
	return nil
}

func (zds *zeroDataSource) GetData(size uint64) ([]byte, error) {
	data := make([]byte, size)
	return data, nil
}

// -------------------------------------------------------------------
type randomDataSource struct {
}

func (rds *randomDataSource) Init(i interface{}) error {
	rand.Seed(time.Now().UnixNano())
	return nil
}

func (rds *randomDataSource) GetData(size uint64) ([]byte, error) {
	data := make([]byte, size)
	n, err := rand.Read(data)

	if err != nil {
		return nil, err
	}

	if uint64(n) != size {
		return nil, fmt.Errorf("failed to generate data of size %v (generated only %v)", size, n)
	}

	return data, nil
}

// -------------------------------------------------------------------
type fileDataSource struct {
	filePath string
	data     []byte
}

func (fds *fileDataSource) Init(i interface{}) error {
	fds.filePath = i.(string)
	f, err := os.Open(fds.filePath)
	if err != nil {
		return err
	}

	data := make([]byte, config.FileSize)
	_, err = f.Read(data)
	if err != nil {
		return err
	}

	f.Close()
	return nil
}

func (fds *fileDataSource) GetData(size uint64) ([]byte, error) {
	return fds.data, nil
}

//-------------------------------------------------------------------

func createDataSource() (dataSource, error) {
	t := config.InputType
	if t == ESourceType.ZERO() {
		return &zeroDataSource{}, nil
	} else if t == ESourceType.RANDOM() {
		return &randomDataSource{}, nil
	} else if t == ESourceType.FILE() {
		f := &fileDataSource{}
		err := f.Init(config.SourceFilePath)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	return nil, fmt.Errorf("invalid source type")
}

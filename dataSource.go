package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

//-------------------------------------------------------------------
type dataSource interface {
	Init(i interface{}) 
	GetData(size uint64) error
}
//-------------------------------------------------------------------

type zeroDataSource struct {

}

func (zds *zeroDataSource) Init(i interface{}) {
	
}

func (zds *zeroDataSource) GetData(size uint64) ([]byte, error) {
	data := make([]byte, size)
	return data, nil
}

//-------------------------------------------------------------------
type randomDataSource struct {
	
}

func (rds *randomDataSource) Init(i interface{}) {
	rand.Seed(time.Now().UnixNano())
}

func (rds *randomDataSource) GetData(size uint64) ([]byte, error) {
	data := make([]byte, size)
	n, err := rand.Read(data)
	
	if err  != nil {
		return nil, err
	}

	if uint64(n) != size {
		return nil, fmt.Errorf("failed to generate data of size %v (generated only %v)", size, n)
	}

	return data, nil
}

//-------------------------------------------------------------------
type fileDataSource struct {
	filePath 		string
}

func (fds *fileDataSource) Init(i interface{}) {
	fds.filePath = i.(string)
}

func (fds *fileDataSource) GetData(size uint64) ([]byte, error) {
	return os.ReadFile(fds.filePath)
}


//-------------------------------------------------------------------
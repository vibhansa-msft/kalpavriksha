package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// -------------------------------------------------------------------
type dataSource interface {
	Init(i interface{}) error
	GetData(size uint64) ([]byte, error)
	GetMd5Sum(data []byte) []byte
}

func getMD5Sum(data []byte) []byte {
	x := md5.Sum(data)
	return x[:]
}

//-------------------------------------------------------------------

type zeroDataSource struct {
	md5sum []byte
}

func (zds *zeroDataSource) Init(i interface{}) error {
	return nil
}

func (zds *zeroDataSource) GetData(size uint64) ([]byte, error) {
	data := make([]byte, size)
	return data, nil
}

func (zds *zeroDataSource) GetMd5Sum(data []byte) []byte {
	return zds.md5sum
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

func (rds *randomDataSource) GetMd5Sum(data []byte) []byte {
	return getMD5Sum(data)
}

// -------------------------------------------------------------------
type fileDataSourceConfig struct {
	filename string
	filesize int64
}

type fileDataSource struct {
	fileDataSourceConfig
	data   []byte
	md5sum []byte
}

func (fds *fileDataSource) Init(i interface{}) error {
	fds.fileDataSourceConfig = i.(fileDataSourceConfig)
	f, err := os.Open(fds.filename)
	if err != nil {
		return err
	}

	data := make([]byte, fds.filesize)
	_, err = f.Read(data)
	if err != nil {
		return err
	}

	f.Close()

	fds.md5sum = getMD5Sum(data)
	return nil
}

func (fds *fileDataSource) GetData(size uint64) ([]byte, error) {
	return fds.data, nil
}

func (fds *fileDataSource) GetMd5Sum(data []byte) []byte {
	return fds.md5sum
}

//-------------------------------------------------------------------

func createDataSource(t SourceType) (dataSource, error) {
	if t == ESourceType.ZERO() {
		return &zeroDataSource{}, nil
	} else if t == ESourceType.RANDOM() {
		return &randomDataSource{}, nil
	} else if t == ESourceType.FILE() {
		f := &fileDataSource{}
		err := f.Init(fileDataSourceConfig{
			filename: config.SourceFilePath,
			filesize: config.FileSize,
		})
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	return nil, fmt.Errorf("invalid source type")
}

package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

// -------------------------------------------------------------------
type dataSource interface {
	Init(i interface{}) error
	GetData() ([]byte, error)
	GetMd5Sum(data []byte) []byte
}

func getMD5Sum(data []byte) []byte {
	x := md5.Sum(data)
	return x[:]
}

//-------------------------------------------------------------------

type zeroDataSource struct {
	size   int64
	md5sum []byte
	data   []byte
}

func (zds *zeroDataSource) Init(i interface{}) error {
	size := i.(int64)
	zds.size = size
	zds.data = make([]byte, size)
	zds.md5sum = getMD5Sum(zds.data)
	return nil
}

func (zds *zeroDataSource) GetData() ([]byte, error) {
	return zds.data, nil
}

func (zds *zeroDataSource) GetMd5Sum(data []byte) []byte {
	return zds.md5sum
}

// -------------------------------------------------------------------
type randomDataSource struct {
	size int64
}

func (rds *randomDataSource) Init(i interface{}) error {
	rand.Seed(time.Now().UnixNano())
	size := i.(int64)
	rds.size = size
	return nil
}

func (rds *randomDataSource) GetData() ([]byte, error) {
	size := rds.size
	if size < 0 {
		// Negative value means generate a random number upto size and create file of that size
		size = rand.Int63n(int64(math.Abs(float64(rds.size))))
	}

	data := make([]byte, size)
	n, err := rand.Read(data)

	if err != nil {
		return nil, err
	}

	if int64(n) != size {
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

	fds.data = make([]byte, fds.filesize)
	n, err := f.Read(fds.data)
	if err != nil {
		return err
	}

	if int64(n) < fds.fileDataSourceConfig.filesize {
		log.Printf("Size : %d, FileSize : %d, remaining will be filled with 0s\n", fds.fileDataSourceConfig.filesize, n)
	}

	f.Close()

	fds.md5sum = getMD5Sum(fds.data)
	return nil
}

func (fds *fileDataSource) GetData() ([]byte, error) {
	return fds.data, nil
}

func (fds *fileDataSource) GetMd5Sum(data []byte) []byte {
	return fds.md5sum
}

//-------------------------------------------------------------------

func createDataSource(t SourceType) (dataSource, error) {
	if t == ESourceType.ZERO() {
		f := &zeroDataSource{}
		err := f.Init(config.FileSize)
		if err != nil {
			return nil, err
		}
		return f, nil

	} else if t == ESourceType.RANDOM() {
		f := &randomDataSource{}
		err := f.Init(config.FileSize)
		if err != nil {
			return nil, err
		}
		return f, nil

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

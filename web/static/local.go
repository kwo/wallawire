package static

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func NewLocalStore(dir string) (AssetStore, error) {

	info, errInfo := os.Stat(dir)
	if errInfo != nil {
		return nil, errInfo
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path does not exist: %s", dir)
	}

	indexPath := path.Join(dir, "index.html")
	index, errIndex := os.Stat(indexPath)
	if errIndex != nil {
		return nil, errIndex
	}

	if index.IsDir() {
		return nil, fmt.Errorf("index file does not exist: %s", indexPath)
	}

	return &LocalStore{
		dir: dir,
	}, nil

}

type LocalAsset struct {
	name string
	info os.FileInfo
	data []byte
}

func (z *LocalAsset) Name() string {
	return z.name
}

func (z *LocalAsset) ModTime() time.Time {
	return z.info.ModTime()
}

func (z *LocalAsset) Data() []byte {
	return z.data
}

type LocalStore struct {
	dir string
}

func (z *LocalStore) AssetFile(name string) Asset {

	filepath := path.Join(z.dir, name)
	info, errInfo := os.Stat(filepath)
	if errInfo == nil {
		data, errData := ioutil.ReadFile(filepath)
		if errData == nil {
			return &LocalAsset{
				name: info.Name(),
				info: info,
				data: data,
			}
		}
	}

	return nil
}

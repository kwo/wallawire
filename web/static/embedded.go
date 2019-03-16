package static

import (
	"time"

	"wallawire/web/assets"
)

func NewEmbeddedStore() (AssetStore, error) {
	return &EmbdeddedStore{}, nil
}

type EmbeddedAsset struct {
	file Asset
}

func (z *EmbeddedAsset) Name() string {
	return z.file.Name()
}

func (z *EmbeddedAsset) ModTime() time.Time {
	return z.file.ModTime()
}

func (z *EmbeddedAsset) Data() []byte {
	return z.file.Data()
}

type EmbdeddedStore struct{}

func (z *EmbdeddedStore) AssetFile(name string) Asset {
	file := assets.AssetFile(name)
	if file != nil {
		return &EmbeddedAsset{
			file: file,
		}
	}
	return nil
}

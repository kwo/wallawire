// Code generated by genesis.
// DO NOT EDIT.

package schema

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"sync"
	"time"
)

var assetMap = map[string]*File{
	"/1_init.sql": &File{
		name:    "/1_init.sql",
		hash:    "fcfd7757c98b10f404d5f20b6a95f2d17e29ee2c19e8cbfb36d4932d84577c25",
		modTime: time.Unix(1552334582, 162511848),
		payload: `
H4sIAAAAAAAC/5yT3Y7TMBCF7/0Uc5mIrrQgxE0lpDSZbq1NncVx+nMVBWxopLSp4kJ5fGTnz5SWwlaqlHh8PnvOmTw8wJt9+a0p
TgqyIwk5BgJBBLMYgc6BJQJwQ1ORwnetGg0eASglOL8so1H/bLazLI7hhdNlwLfwjNsJAZClLj5XqtPNkiS+VJhd5oRDsVdtZRXw
cBFw78N7f+SGCwyfwYuRPYmFNxOcLr1e5vvwER59yBj9lKEBjrD/AjowQzkWWp/rRua7Qu/AIt6+e/Tddq9RfpM5uC+NKk69FUCZ
wCfkf3pxlHd2EX9K+rxWFNeDfRpGS0iQkhRjDAXEyRq5U5nzZNmGOnKu5d7UlXJy/5e4re+vsbvLzm3t1ijm5l72WvatnUl7ueEw
jnPkyEIchreUNgIjvafo2u4UP4qqlPnXpt73WYyrp3pMyKw6ZrRJ5KWc9Gf6bm+URbi56K2UPzOtGl5XyvwhYW67HWXSN21g7icc
1ecDiXjy0o6EYTue2emYtvXB10tPb9StH3/R6in5FQAA///Xz1UcTAQAAA==
`,
	},
	"/2_data.sql": &File{
		name:    "/2_data.sql",
		hash:    "6dec4ebda957b47cf3f750871a544b3688c6bbd7d1f53c5c0c58474eb1ad28b0",
		modTime: time.Unix(1552334582, 162692082),
		payload: `
H4sIAAAAAAAC/7RTXYsTMRR931+Rt7Q4gSQ3994MPhU7iwt1V/qhvpVkkrGFbad0Vvr3ZToqVURG0Lf7Aecczr1HKfHqsP98Di9Z
bE53m/erarkWD4/rJ3Fun3MnJvtUiGM45Kn4MFtsqpWYSI05ETAqzZ6UyyaoWEetDALUHrhEl2QhZEiH/VFOX4+EDbFsbKmNwpCD
chBJ+cZGxbEpOSXQXvse9kuXzz3qT7D98Bts2nchPudUXIc9yUBViPqcw8t1cUpDcQpdd2nPabsL3W5690ML6AQxmqBC0wTlfCLl
sawV1AaBmzKnGGUh7meLVVUImfKhvcoqhJznQys2Q1N9Wi9nb9YTmU9tvZOFOLaXyXT6h4W0DiwZ68CAtg4NOkfkkNizcUyWsgsO
mYkQwSBSxMAaNCETokMqQaMHBnClzUAEDGgoEjGwoZIyaHBUu5qS82A5sGVw/FtHt/21xORa9tb27Xafbo420qhxP/PLq/xbBaPe
qzfhNhPz9nK8m1eLal2J++XTuxtNH99Wy+q7HvHw+DfBGCnllnjIzUD6v/i+BgAA//8rMoMEDgQAAA==
`,
	},
}

var assetNames = []string{
	"/1_init.sql",
	"/2_data.sql",
}

// File represents a single embedded asset file.
type File struct {
	name    string
	hash    string
	modTime time.Time
	payload string
	data    []byte
	once    sync.Once
}

// Name returns the full path of the file.
func (f *File) Name() string { return f.name }

// Hash returns the SHA256 hash of the file's data.
func (f *File) Hash() string { return f.hash }

// ModTime returns the last modified date of the file when it was generated.
func (f *File) ModTime() time.Time { return f.modTime }

// Data returns the raw embedded data for the file.
func (f *File) Data() []byte {
	f.once.Do(func() {
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.payload))
		gr, errReader := gzip.NewReader(b64)
		if errReader != nil {
			return
		}
		data, err := ioutil.ReadAll(gr)
		if err != nil {
			return
		}
		f.data = data
	})
	return f.data
}
// Asset returns the raw data given an embedded filename.
// Returns nil if the asset cannot be found.
func Asset(name string) []byte {
	if f := AssetFile(name); f != nil {
		return f.Data()
	}
	return nil
}

// AssetFile returns the File object given an embedded filename.
// Returns nil if the asset cannot be found.
func AssetFile(name string) *File {
	if f := assetMap[name]; f != nil {
		return f
	}
	return nil
}

// AssetNames returns a sorted list of all embedded asset filenames.
func AssetNames() []string {
	return assetNames
}


package disc

import (
	"github.com/zoroqi/rubbish/disc/hash"
	"io/fs"
	"path/filepath"
)

type FileMeta struct {
	name       string
	path       string
	fastMd5sum string
	md5sum     string
	dir        bool
	size       int64
	err        error
}

var digest hash.Digest = hash.Crc32

func FastScan(path string) []FileMeta {
	meta := make([]FileMeta, 0, 1024)
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var fastDigest string
		if !info.IsDir() {
			fastDigest, err = hash.FileFastMd5sum(path, 1024)
		}

		meta = append(meta, FileMeta{
			path:       path,
			name:       info.Name(),
			fastMd5sum: fastDigest,
			dir:        info.IsDir(),
			size:       info.Size(),
			err:        err,
		})
		return err
	})
	return meta
}

func SlowScan(files []FileMeta) ([]FileMeta, error) {
	for i := 0; i < len(files); i++ {
		if !files[i].dir {
			dis, err := hash.FileMd5sum(files[i].path, digest)
			files[i].md5sum = dis
			files[i].err = err
		}
	}
	return files, nil
}

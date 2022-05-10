package hash

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func FileFastMd5sum(path string, digestSize int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bs := make([]byte, digestSize)
	n, err := f.Read(bs)
	if err != nil && err != io.EOF {
		return "", err
	}

	b := md5.Sum(bs[:n])
	return hex.EncodeToString(b[0:]), nil
}

type Digest = func(string) (string, error)

func Md5sum(path string) (string, error) {
	cmd := exec.Command("md5sum", path)
	outBuffer := &bytes.Buffer{}
	cmd.Stdout = outBuffer
	err := cmd.Run()
	cmd.Run()
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(outBuffer)
	if err != nil {
		return "", err
	}
	out := string(bs)
	split := strings.Split(out, " ")
	if len(split) == 2 {
		return strings.TrimSpace(split[0]), nil
	}
	return "", errors.New(fmt.Sprintf("md5sum error, %s", out))
}

func Crc32(path string) (string, error) {
	cmd := exec.Command("crc32", path)
	outBuffer := &bytes.Buffer{}
	cmd.Stdout = outBuffer
	err := cmd.Run()
	cmd.Run()
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(outBuffer)
	if err != nil {
		return "", err
	}
	out := string(bs)
	return strings.TrimSpace(out), nil
}

func FileMd5sum(path string, digest Digest) (string, error) {
	return digest(path)
}

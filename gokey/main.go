package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"flag"
	"fmt"
	"golang.design/x/clipboard"
	"golang.org/x/term"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

var (
	salt      = flag.String("s", "", "salt")
	length    = flag.Int("l", 16, "length")
	plaintext = flag.String("p", "", "plaintext")
	charset   = flag.Int("charset", 0, "charset, default base32")
	clip      = flag.Bool("c", false, "copy to clipboard")
)

func main() {
	flag.Parse()
	initCharsets(charsets)
	if *plaintext == "" {
		fmt.Println("Error: plaintext is required")
		return
	}
	if *salt == "" {
		var saltBytes []byte
		var saltBytesAgain []byte
		var err error
		for {
			for len(saltBytes) == 0 {
				fmt.Fprint(os.Stderr, "salt: ")
				saltBytes, err = term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Fprintln(os.Stderr, "")
			}

			fmt.Fprint(os.Stderr, "salt again: ")
			saltBytesAgain, err = term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Fprintln(os.Stderr, "")
			if bytes.Equal(saltBytes, saltBytesAgain) {
				break
			} else {
				fmt.Fprintln(os.Stderr, "salt do not match. Try again.")
				saltBytes = nil
				continue
			}
		}
		*salt = string(saltBytes)
	}
	// 重盐中解析长度和字符集
	if ss := strings.SplitN(*salt, "_", 3); len(ss) >= 3 {
		if ss[2] == "" {
			fmt.Println("Error: salt is required")
			return
		}
		l, err := strconv.Atoi(ss[0])
		if err == nil && (*length < 5 || *length > 64) {
			*length = l
		}
		c, err := strconv.Atoi(ss[1])
		if err == nil && charsets[c] != "" {
			*charset = c
		}
	}

	if *length < 5 || *length > 64 {
		fmt.Println("Error: length must be [5,64]")
		return
	}

	if c := charsets[*charset]; c == "" {
		fmt.Println("charsets don't exist, so use default base32")
	}

	key, err := generateKeyWithSha512(*plaintext, *salt, *length, *charset)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if *clip {
		clipboard.Init()
		clipboard.Write(clipboard.FmtText, []byte(key))
	} else {
		fmt.Println(key)
	}
}

var charsets = map[int]string{
	// num: 16/64(0.25); letter: 36/64(0.56); specific: 12/64(0.18)
	1: "0123456789abcdefghijklmabcdeNOPQRNOPQRSTUVWXYZ!@#!@#!@#!@#012345",
	// num: 16/64(0.25); letter: 36/64(0.56); specific: 12/64(0.18)
	2: "0123456789abcdefghijklmabcdeNOPQRNOPQRSTUVWXYZ!@#$%^!@#$%^012345",
	// num: 16/64(0.25); letter: 36/64(0.56); specific: 12/64(0.18)
	3: "0123456789abcdefghijklmabcdeNOPQRNOPQRSTUVWXYZ!@#$%^,./<>?012345",
	// num: 18/64(0.28); letter: 47/64(0.72);
	4: "0123456789abcdefghijklmabcdeNOPQRNOPQRSTUVWXYZijklmVWXYZ01234567",
}
var baseEncoding = map[int]encoding{}

type encoding interface {
	EncodeToString([]byte) string
}

func initCharsets(charsets map[int]string) {
	for i, v := range charsets {
		switch len(v) {
		case 64:
			baseEncoding[i] = base64.NewEncoding(v)
		case 32:
			baseEncoding[i] = base32.NewEncoding(v)
		default:
			baseEncoding[i] = base32.StdEncoding
		}
	}
}

func generateKeyWithSha512(plaintext, salt string, length int, charset int) (string, error) {
	// 将盐和明文组合
	combined := []byte(fmt.Sprintf("%s_%s", plaintext, salt))

	// 使用 SHA-512 算法生成哈希值
	slices.Reverse(combined)
	hash := sha512.Sum512(combined)

	// 使用 base32 编码生成密钥
	encoding, ok := baseEncoding[charset]
	if !ok {
		encoding = base32.StdEncoding
	}

	encoded := encoding.EncodeToString(hash[:])
	// 截取指定长度的密钥
	truncatedKey := encoded[:length]

	return truncatedKey, nil
}

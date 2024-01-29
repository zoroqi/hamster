package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
)

func main() {
	c := flag.Bool("c", false, "compress or uncompress")
	i := flag.String("i", "", "input file")
	o := flag.String("o", "", "output file")

	flag.Parse()
	if *i == "" {
		fmt.Println("input file is required")
		return
	}
	if *o == "" {
		fmt.Println("output file is required")
		return
	}
	input, err := os.Open(*i)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer input.Close()
	output, err := os.OpenFile(*o, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer output.Close()

	if *c {
		err = compress(input, output)
	} else {
		err = uncompress(input, output)
	}
	if err != nil {
		fmt.Println(err)
	}
}

const file_size_byte_len = 8

func uncompress(input io.Reader, output io.Writer) error {
	// Step 1: Decode the image
	image, _, err := image.Decode(input)
	if err != nil {
		return err
	}
	// Step 2: Get the file size
	filesize := make([]byte, file_size_byte_len)
	flag := true
	index := int64(0)
	l := int64(math.MaxInt64)

	// 简单优化以下, 1024个字节写一次
	const buffersize = 1024
	buffer := [buffersize]byte{}
	bufferIndex := 0
	// Step 3: Encode the image to the file
Outer:
	for y := 0; y < image.Bounds().Dy(); y++ {
		for x := 0; x < image.Bounds().Dx(); x++ {
			rgba := image.At(x, y)
			r, _, _, _ := rgba.RGBA()
			if flag {
				filesize[index] = byte(r)
				index++
				if index >= file_size_byte_len {
					l = unsize(filesize)
					index = 0
					flag = false
				}
			} else {
				buffer[bufferIndex] = byte(r)
				index++
				bufferIndex++
				if bufferIndex >= buffersize {
					if _, err := output.Write(buffer[:]); err != nil {
						return err
					}
					bufferIndex = 0
				}
				if index >= l {
					break Outer
				}
			}
		}
	}
	if bufferIndex > 0 {
		if _, err := output.Write(buffer[:bufferIndex]); err != nil {
			return err
		}
	}
	return nil
}

func compress(input io.Reader, output io.Writer) error {
	// Step 1: Read the file
	bs, err := io.ReadAll(input)
	if err != nil {
		return err
	}
	// Step 2: Create an image
	l := int64(len(bs))
	imgSize := int(math.Sqrt(float64(l+file_size_byte_len))) + 1
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	index := int64(0)

	filesize := size(l)
	flag := true
	// Step 3: Set the color of each pixel
Outer:
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if flag {
				img.Set(x, y, color.Gray{filesize[index]})
				index++
				if index >= file_size_byte_len {
					index = 0
					flag = false
				}
			} else {
				img.Set(x, y, color.Gray{bs[index]})
				index++
				if index >= l {
					break Outer
				}
			}
		}
	}

	// Step 4: Encode the image to the file
	if err := png.Encode(output, img); err != nil {
		return err
	}
	return nil
}

func size(l int64) []byte {
	bs := make([]byte, 8)
	bs[0] = byte(l >> 56)
	bs[1] = byte(l >> 48)
	bs[2] = byte(l >> 40)
	bs[3] = byte(l >> 32)
	bs[4] = byte(l >> 24)
	bs[5] = byte(l >> 16)
	bs[6] = byte(l >> 8)
	bs[7] = byte(l)
	return bs
}

func unsize(bs []byte) int64 {
	var l int64
	l |= int64(bs[0]) << 56
	l |= int64(bs[1]) << 48
	l |= int64(bs[2]) << 40
	l |= int64(bs[3]) << 32
	l |= int64(bs[4]) << 24
	l |= int64(bs[5]) << 16
	l |= int64(bs[6]) << 8
	l |= int64(bs[7])
	return l
}

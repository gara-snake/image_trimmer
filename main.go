package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// 4 : 3
const height = 960
const width = 1280

func main() {

	flag.Parse()

	files := flag.Args()

	for _, f := range files {

		file, e1 := os.Open(f)
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}()

		if e1 != nil {
			fmt.Println(e1)
			return
		}

		imgSrc, t, e2 := image.Decode(file)

		if e2 != nil {
			fmt.Println(e2)
			return
		}

		fmt.Println("処理中..." + file.Name())

		rctSrc := imgSrc.Bounds()

		rw := float64(width) / float64(rctSrc.Dx())
		rh := float64(height) / float64(rctSrc.Dy())

		var sizeRate float64

		if rw < rh {
			sizeRate = rh
		} else {
			sizeRate = rw
		}

		x := int(float64(rctSrc.Dx()) * sizeRate)
		y := int(float64(rctSrc.Dy()) * sizeRate)

		imgDst := image.NewRGBA(image.Rect(0, 0, x, y))
		draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), imgSrc, rctSrc, draw.Over, nil)

		imgWeb := image.NewRGBA(image.Rect(0, 0, width, height))

		ypoint := y - height

		draw.Draw(imgWeb, imgWeb.Bounds(), imgDst, image.Point{0, ypoint}, draw.Over)

		d, n := filepath.Split(file.Name())

		ext := filepath.Ext(n)

		n = strings.ReplaceAll(n, ext, "")

		n = "web_" + n + ext

		dst, e3 := os.Create(filepath.Join(d, n))
		if e3 != nil {
			fmt.Println(e3)
			return
		}
		defer dst.Close()

		switch t {
		case "jpeg":
			if err := jpeg.Encode(dst, imgWeb, &jpeg.Options{Quality: 100}); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "gif":
			if err := gif.Encode(dst, imgWeb, nil); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "png":
			if err := png.Encode(dst, imgWeb); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		default:
			fmt.Fprintln(os.Stderr, "format error")
		}

	}

}

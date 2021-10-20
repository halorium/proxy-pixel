package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {
	// get origin from cli flag argument
	origin := flag.String("origin", "", "The origin server")
	flag.Parse()

	fmt.Println("listening at: http://localhost:3000")
	fmt.Printf("origin: %s\n", *origin)

	handler := getHandler(*origin)

	http.ListenAndServe(":3000", handler)
}

func getHandler(origin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse the url
		url, _ := url.Parse(origin)

		// create the reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(url)

		// set timeout
		proxy.Transport = &http.Transport{
			ResponseHeaderTimeout: 5 * time.Second,
		}

		// Update the headers
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host

		proxy.ModifyResponse = convertImage
		proxy.ServeHTTP(w, r)
	}
}

func convertImage(res *http.Response) error {
	img, imgType, err := image.Decode(res.Body)
	if err != nil {
		return err
	}

	// colorPalette := img.ColorModel()

	// create grayscale image
	bounds := img.Bounds()
	grayScale := image.NewGray(
		image.Rectangle{
			image.Point{bounds.Min.X, bounds.Min.Y},
			image.Point{bounds.Max.X, bounds.Max.Y}})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldPixel := img.At(x, y)
			pixel := color.GrayModel.Convert(oldPixel)
			grayScale.Set(x, y, pixel)
		}
	}

	var b bytes.Buffer
	newImg := bufio.NewWriter(&b)

	if imgType == "png" {
		err = png.Encode(newImg, grayScale)
		if err != nil {
			return err
		}
	} else if imgType == "jpeg" {
		err = jpeg.Encode(newImg, grayScale, &jpeg.Options{})
		if err != nil {
			return err
		}
	}

	res.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))

	return nil
}

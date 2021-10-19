package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
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

	// create grayscale image
	// bounds := img.Bounds()
	// grayScale := image.NewGray(
	// 	image.Rectangle{
	// 		image.Point{bounds.Min.X, bounds.Min.Y},
	// 		image.Point{bounds.Max.X, bounds.Max.Y}})

	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
	// 		imageColor := img.At(x, y)
	// 		rr, gg, bb, _ := imageColor.RGBA()
	// 		r := math.Pow(float64(rr), 2.2)
	// 		g := math.Pow(float64(gg), 2.2)
	// 		b := math.Pow(float64(bb), 2.2)
	// 		m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
	// 		Y := uint16(m + 0.5)
	// 		grayColor := color.Gray{uint8(Y >> 8)}
	// 		grayScale.Set(x, y, grayColor)
	// 	}
	// }

	// fmt.Printf("grayScale: %#v\n", grayScale)

	// for x := 0; x < w; x++ {
	// 	for y := 0; y < h; y++ {
	// 		imageColor := img.At(x, y)
	// 		rr, gg, bb, _ := imageColor.RGBA()
	// 		r := math.Pow(float64(rr), 2.2)
	// 		g := math.Pow(float64(gg), 2.2)
	// 		b := math.Pow(float64(bb), 2.2)
	// 		m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
	// 		Y := uint16(m + 0.5)
	// 		grayColor := color.Gray{uint8(Y >> 8)}
	// 		grayScale.Set(x, y, grayColor)
	// 	}
	// }

	var b bytes.Buffer
	newImg := bufio.NewWriter(&b)

	if imgType == "png" {
		// err = png.Encode(newImg, grayScale)
		err = png.Encode(newImg, img)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("img type: %#v\n", imgType)
	}

	res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))

	return nil
}

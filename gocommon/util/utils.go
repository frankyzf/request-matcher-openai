package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"moul.io/http2curl"

	"github.com/nfnt/resize"
)

func PrintResponse(resp *http.Response) string {
	buf, _ := ioutil.ReadAll(resp.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

	buf2 := new(bytes.Buffer)
	buf2.ReadFrom(rdr1)
	s := buf2.String()
	size := len(s)
	if size > 1024 {
		size = 1024
	}
	fmt.Printf("response body:%v\n", s[:size])
	resp.Body = rdr2
	return s
}

func curlAppend(c http2curl.CurlCommand, newSlice ...string) http2curl.CurlCommand {
	d, _ := interface{}(c).([]string)
	d = append(d, newSlice...)
	k, _ := interface{}(d).(http2curl.CurlCommand)
	return k
}

type nopCloser struct {
	io.Reader
}

func bashEscape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}

func (nopCloser) Close() error { return nil }

func PrintCurlCommand(req *http.Request) (*http2curl.CurlCommand, error) {
	command := http2curl.CurlCommand{}

	command = curlAppend(command, "curl")

	command = curlAppend(command, "-X", bashEscape(req.Method))

	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = nopCloser{bytes.NewBuffer(body)}
		size := len(string(body))
		if size > 1024 {
			size = 1024
		}
		if len(string(body)) > 0 {
			bodyEscaped := bashEscape(string(body[0:size]))
			command = curlAppend(command, "-d", bodyEscaped)
		}
	}

	var keys []string

	for k := range req.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command = curlAppend(command, "-H", bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header[k], " "))))
	}

	command = curlAppend(command, bashEscape(req.URL.String()))

	return &command, nil
}

func ResizeJpeg(buf []byte, width, height uint) ([]byte, error) {
	// format, _, _ := image.DecodeConfig(bytes.NewReader(buf))
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return []byte{}, err
	}
	// img = resize.Resize(width, 0, img, resize.Lanczos3)
	img = resize.Thumbnail(width, height, img, resize.Lanczos3)
	out := bytes.NewBuffer(nil)
	err = jpeg.Encode(out, img, nil)
	data := out.Bytes()
	return data, err
}

func ResizeJpegRatio(buf []byte, ratio float64) ([]byte, error) {
	format, _, _ := image.DecodeConfig(bytes.NewReader(buf))
	height := uint(float64(format.Height) * ratio)
	width := uint(float64(format.Width) * ratio)
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return []byte{}, err
	}
	// img = resize.Resize(width, 0, img, resize.Lanczos3)
	img = resize.Thumbnail(width, height, img, resize.Lanczos3)
	out := bytes.NewBuffer(nil)
	err = jpeg.Encode(out, img, nil)
	data := out.Bytes()
	return data, err
}

func CompressImage(data []byte, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}
	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func ReadOrientation(mmbuf []byte) (int, error) {
	x, err := exif.Decode(bytes.NewReader(mmbuf))
	if err != nil {
		return 0, err
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return 0, err
	}
	orientVal, err := orientation.Int(0)
	if err != nil {
		return 0, err
	}
	return orientVal, nil
}

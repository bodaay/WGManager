package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

// Got it from: https://gist.github.com/xyproto/f4915d7e208771f3adc4

//GzipWrite Write gzipped data to a Writer
func GzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	defer gw.Close()
	gw.Write(data)
	return err
}

//GUnzipWrite Write gunzipped data to a Writer
func GUnzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}

/*
//Example:
	s := "some data"
	fmt.Println("original:\t", s)
	var buf bytes.Buffer
	err := gzipWrite(&buf, []byte(s))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("compressed:\t", buf.String())
	var buf2 bytes.Buffer
	err = gunzipWrite(&buf2, buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("decompressed:\t", buf2.String())
*/

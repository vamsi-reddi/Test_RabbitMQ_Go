package config

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func ReadConfigFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("error reading config file, ", err)
		return "", err
	}

	return string(data), err
}


func Decompress(compressedData string) (string, error) {
	
	data, err := base64.StdEncoding.DecodeString(compressedData)
	if err != nil {
		return "", err
	}

	rdata := bytes.NewReader(data)
	r, _ := zlib.NewReader(rdata)

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String(), nil
	

}
package demodulation

import (
	"bytes"
	"log"

	"github.com/ulikunitz/xz"
)

func decompress(data []byte) []byte {
	buf := bytes.NewReader(data)

	r, err := xz.NewReader(buf)

	if err != nil {
		log.Fatalf("NewReader error %s", err)
	}

	var reader bytes.Buffer

	reader.ReadFrom(r)

	return reader.Bytes()
}

package transmission

import (
	"bytes"
	"log"

	"github.com/ulikunitz/xz"
)

func compress(data []byte) []byte {
	var buf bytes.Buffer
	// compress text
	w, err := xz.NewWriter(&buf)
	if err != nil {
		log.Fatalf("xz.NewWriter error %s", err)
	}

	w.Write(data)

	if err := w.Close(); err != nil {
		log.Fatalf("w.Close error %s", err)
	}

	return buf.Bytes()
}

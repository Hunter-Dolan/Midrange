package transmission

import (
	"crypto/sha1"

	"github.com/Hunter-Dolan/midrange/fec"
	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/utils"
)

type packet struct {
	data *[]bool
	hash *[]bool

	hashBytes []byte

	hashConfidenceCollection *fec.DataConfidenceCollection
	dataConfidenceCollection *fec.DataConfidenceCollection
}

func (p *packet) setBytes(b []byte) {
	binaryData := utils.BinEncode(b)
	p.data = &binaryData
	p.calculateHash(b)
}

func (p *packet) calculateHash(b []byte) {
	sha1Data := sha1.Sum(b)

	p.hashBytes = sha1Data[:options.PacketChecksumLengthBytes]
	binaryHash := utils.BinEncode(p.hashBytes)
	p.hash = &binaryHash
}

func (p *packet) replicateHash(level int) {
	data := *p.data
	for i := 0; i < level; i++ {
		data = append(data, (*p.hash)...)
	}

	p.data = &data
}

func (p *packet) packetData() []bool {
	packetData := append(*p.data, *p.hash...)

	return packetData
}

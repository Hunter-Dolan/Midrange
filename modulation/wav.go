package modulation

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func (m *Modulator) BuildWav(filename string) {

	numberOfFrames := m.FrameCount()

	rate := m.Options.Kilobitrate * 1000
	bitDepth := 16
	duration := numberOfFrames * m.Options.FrameDuration

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	w := bufio.NewWriter(file)
	dataSize := rate * (duration / 1000) * (bitDepth / 8)
	buildWavHeader(w, bitDepth, rate, dataSize)

	for i := 0; i < numberOfFrames; i++ {
		wave := *m.NextWave()

		for _, amplitude := range wave {
			binary.Write(w, binary.LittleEndian, int16(amplitude))
		}
	}

	w.Flush()

}

const riff = "RIFF"
const wave = "WAVE"
const _fmt = "fmt "
const data = "data"

func buildWavHeader(buf io.Writer, bitDepth, sampleRate, dataSize int) {

	buf.Write([]byte(riff))
	binary.Write(buf, binary.LittleEndian, uint32(36+dataSize))
	buf.Write([]byte(wave))
	buf.Write([]byte(_fmt))
	binary.Write(buf, binary.LittleEndian, uint32(16))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))

	//ByteRate
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate*(bitDepth/8)))
	binary.Write(buf, binary.LittleEndian, uint16((bitDepth / 8)))
	binary.Write(buf, binary.LittleEndian, uint16(bitDepth))
	buf.Write([]byte(data))
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))
}

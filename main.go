package main

import (
	"fmt"

	"github.com/Hunter-Dolan/midrange/options"
	"github.com/Hunter-Dolan/midrange/transmission"
)

func main() {
	options := options.Options{}
	options.BaseFrequency = 1000
	options.FrameDuration = 500
	options.Kilobitrate = 96 * 2
	options.Bandwidth = 1000
	options.NoiseLevel = 20
	options.OMFSKConstant = 2.5
	options.NFFTPower = 17

	sentTransmission := transmission.NewTransmission()

	sentTransmission.SetOptions(options)

	data := `
		GNU GENERAL PUBLIC LICENSE
			 Version 3, 29 June 2007

Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
Everyone is permitted to copy and distribute verbatim copies
of this license document, but changing it is not allowed.

						Preamble

The GNU General Public License is a free, copyleft license for
software and other kinds of works.

The licenses for most software and other practical works are designed
to take away your freedom to share and change the works.  By contrast,
the GNU General Public License is intended to guarantee your freedom to
share and change all versions of a program--to make sure it remains free
software for all its users.  We, the Free Software Foundation, use the
GNU General Public License for most of our software; it applies also to
any other work released this way by its authors. You can apply it to
your programs, too.
	`

	sentTransmission.SetData(data)

	sentTransmission.BuildWav("tone.wav")

	signal := sentTransmission.Wave()

	recievedTransmission := transmission.NewTransmission()
	recievedTransmission.SetOptions(options)
	sentTransmission.SetWave(signal)

	fmt.Println(string(sentTransmission.Data()[:]))

	fmt.Println("Transfer successful: ", string(sentTransmission.Data()[:]) == data)

}

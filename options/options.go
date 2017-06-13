package options

type Options struct {
	BaseFrequency int
	FrameDuration int
	Kilobitrate   int
	Bandwidth     int

	OMFSKConstant float64

	NFFTPower float64

	// Debug
	NoiseLevel int
}

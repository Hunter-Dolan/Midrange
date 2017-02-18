package options

const HeaderDataLengthIntegerBitLength = 20
const HeaderTransmissionChecksumLength = 32
const HeaderChecksumRedundancy = 2

const PacketTotalLength = 256
const PacketChecksumLength = 32
const PacketChecksumLengthBytes = PacketChecksumLength / 8

const PacketDataLength = PacketTotalLength - PacketChecksumLength
const PacketDataLengthBytes = PacketDataLength / 8

const HeaderLengthBits = HeaderTransmissionChecksumLength + HeaderDataLengthIntegerBitLength + (HeaderChecksumRedundancy * PacketChecksumLength)
const HeaderLengthWithoutRedundancy = HeaderTransmissionChecksumLength + HeaderDataLengthIntegerBitLength

const MaxCorrectionRoundsHashValidated = 20
const MaxCorrectionRoundsHashNotValidated = 20

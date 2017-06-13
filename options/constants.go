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

const MaxCorrectionRoundsHashValidated = 12
const MaxCorrectionRoundsHashNotValidated = 12

const BitDepth = 32

const OffsetGroups = 1
const GuardDivisor = 1

const HeaderInterval = 2

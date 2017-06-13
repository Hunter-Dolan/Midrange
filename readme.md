# Midrange

Work in progress.

Disclaimer: Midrange still is under very active development. You probably shouldn't use it for anything yet. The protocol, packet structure, modulation techniques, and demodulation techniques will change.

Midrange is also very slow right now. The only focus is the speed, spectral efficiency, and reliability of transmitting large amounts of data in high noise environments.

Midrange focuses extensively on being the new digital standard for HF radio communications.

It is currently being Developed by Hunter Dolan, and software engineer and not a ham radio operator. 

Anyone who wants to provide assistance, comments, questions, concerns, etc can feel free to contact him directly at hunter@wearemocha.com.

# Goals

HF radio bands allow for intercontinental communication between radios at very low power through the use of skywave propagation. The unique properties of frequencies in this band have been exploited since the early days of radio communications.

Midrange's goal is to provide a new open standard for HF digital communications.

Midrange will be able to offer a high bitrate on a very low bandwidth channel. With adaptive symbol rates, advanced forward error correction, Midrange signals can still be decoded slightly below the noise floor.

The goals of midrange are as follows:

Open Header System:

- Provide an open header standard that will allow midrange to evolve in the future and maintain reverse compatibility. Headers will be modulated at a low bit rate in order to ensure decodablity under even under weak signal conditions. And will provide:

  - 4 bit Modulation Code: The type of modulation the data will be encoded with. Codes 0-4 will be designated for future midrange development while 5-14 will be designated for alternate (non midrange) modulation schemes. Finally code 15 will be reserved for future use (most likely to allow for extension bits to be added to the header for an even greater modulation set) 
 - 12 bit Modulation Configuration
 - 42 bit Call Sign
 - 8 bit length of transmission in quarter seconds (not applicable in beacon mode)
 - 16 bit checksum

- Enable automatic link establishment
- Enable frequency sharing (using TDM)

Midrange 1:

- Use the open header system
- Enable both beacon communications, 1 to 1 communications, and n to n communications
- Capable of transmitting up to 1 bit per hz
- Capable of being decoded at or even below the noise floor (at a lower bit rate)
- Highly configurable (configurable bandwidth, carrier count, frame rate, amplitude levels, etc)
- Able to use forward error correction to correct a large number of bits
- Can be modulated and demodulated with a SoundCard (no special hardware)
- Is able to mitigate multipath, high noise, doppler shift, and other problems commonly experienced on the HF band
- Able to be used on any mode of transport including walki talki, sound and microphone, fm radio, etc.

Midrange 2:

- Employs all features of Midrange 1
- Capable of being decoded well below the noise floor (at a very low bit rate)
- Capable of transmitting up to 2.5 bits for every 1hz of bandwidth

Midrange 3:

- Employs all of the features of Midrange 2
- Enables communications to be shared on the same band (through the use of FDM and TDM)
- Enables a real time communication registry that will allow frequencies to be scheduled for transmissions


Overall Vision:

Midrange will be used to connect the furthest reaches of the globe. It will be used as an alternative to costly satellite systems, expensive and monopolized high rate proprietary HF modes. It does not intend to, will it be able to, compete with high rate modes in the VHF and UHF spectrum.

Potential Applications include:

- HAM Communications (including digital audio)
- Maritime Communications
- Remote Research Communications
- Commercial Trunked Communication Services

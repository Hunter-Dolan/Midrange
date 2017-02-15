# Time Burst Modulation

A method of modulation that stores information in time between signals.

In it's simplest form, on / off keying may be used. More complex applications may use frequency shifts, or phase shifts.

The time between a shift defines the value of the data being transmitted. Data is not sent in traditional bits, but is instead sent in a configurable base. In this case we will use base 8.

In base 8, three bits can be sent at once. If we have the bits 010, in base 8 that would equal 2.

As a result the duration between the shifts will equal 2 times the usable interval.

The usable interval equals the amount of time each time-slot should exist during. The lower the usable interval the higher the data rate per channel, but the more your receiver will need to be synchronized.

A TBM signal can be thought of in sections. Each frequency is a carrier. Each data segment being carried each carrier is called a frame. And each interval in the usable interval is a window.

The carrier wave will shift (depending on the modulation technique) when the time the carrier has been at one state equals the midpoint of the window that will be sent.

If we have 010 to send, we send it in base 8, and we use the usable interval of 1 second. We will shift at 2 times the usable interval.

Which means our transmission can be anywhere between 0.5 seconds and 7.5 seconds. Since we are sending 2 the window we will need to signal is 2 times 1 or Window 2. This means our transmission will shift at 2.5 seconds.

Assuming our last transmission was "ON".

HHHHHL
-----_
|0|1|2|3|4|5|6|7|

Now lets assume we need to send 9 bits of data. And our data is

001 010 101

or in base 8:

1 2 5

Our signal would look something like this

HHHL
---_
|0|1| 

*Shift: So switching to next frame*

LLLLLH
_____-
|0|1|2|

*Shift again!*

HHHHHHHHHHHL
-----------_
|0|1|2|3|4|5|


Our transmission time would be 11 seconds for 9 bits.

# Transmission Rate

Using Base 8 to transmit 1000 bits or 334 frames

Usable Interval of 500ms: Maximum is 334 * 500ms = 167.0s / Minimum is 334 * 500ms = 3.34s
Usable Interval of 50ms: Maximum is 334 * 50ms = 16.7s / Minimum is 334 * 50ms = 0.334s
Usable Interval of 10ms: Maximum is 334 * 10ms = 3.34s / Minimum is 334 * 10ms = 0.334s
Usable Interval of 1ms: Maximum is 334 * 1ms = 0.334s / Minimum is 334 * 1ms = 0.0334s

Bitrate: 

bitrate(interval,frame_size) = 
    let result = [1000/(1000*interval*frame_size), 1000/(1000*(interval/2))] in
    let midpoint = (result[0] + result[1]) / 2 in
    [result[0], result[1], midpoint]


bitrate(1,3) => [0.3333, 2, 1.1667]
bitrate(10,3) => [0.0333, 0.2, 0.1167]
bitrate(50,3) => [0.0067, 0.04, 0.0233]
bitrate(500, 3) => [0.0007, 0.004, 0.0023]


# Multi carrier

multicarrier_bitrate(interval, frame_size, carriers) =
    bitrate(interval, frame_size) * carriers 

250 carriers in a 1000hz frame @ 500ms: 

multicarrier_bitrate(500, 3, 250) => [0.1667, 1, 0.5833]

250 carriers in a 1000hz frame @ 50ms:

multicarrier_bitrate(50, 3, 250) => [1.6667, 10, 5.8333]

500 carriers in a 1000hz frame @ 500ms:

multicarrier_bitrate(500, 3, 500) => [0.3333, 2, 1.1667]

# Multiplexed Multi Carrier

multiplexed_multicarrier_bitrate(interval, frame_size, carriers, multiplex) =
    multicarrier_bitrate(interval, frame_size, carriers) * multiplex 

250 carriers in a 1000hz frame @ 500ms with a multiplexed level of 2:

multiplexed_multicarrier_bitrate(500, 3, 250, 2) => [0.3333, 2, 1.1667]

250 carriers in a 1000hz frame @ 50ms with a multiplexed level of 2:

multiplexed_multicarrier_bitrate(50, 3, 250, 2) => [3.3333, 20, 11.6667]

500 carriers in a 1000hz frame @ 500ms with a multiplexed level of 2:

multiplexed_multicarrier_bitrate(500, 3, 500, 2) => [0.6667, 4, 2.3333]

500 carriers in a 1000hz frame @ 500ms with a multiplexed level of 4:

multiplexed_multicarrier_bitrate(500, 3, 500, 4) => [1.3333, 8, 4.6667]


# Application

Reasonable:
multiplexed_multicarrier_bitrate(500, 3, 250*2.5, 1) => [0.4167, 2.5, 1.4583]
multiplexed_multicarrier_bitrate(500, 3, 250*2.5, 2) => [0.8333, 5, 2.9167]

multiplexed_multicarrier_bitrate(500, 3, 250*8, 1) => [1.3333, 8, 4.6667]
multiplexed_multicarrier_bitrate(500, 3, 250*8, 2) => [2.6667, 16, 9.3333]


Potentially Reasonable:
multiplexed_multicarrier_bitrate(500, 3, 250*2.5, 4) => [1.6667, 10, 5.8333]
multiplexed_multicarrier_bitrate(500, 3, 500*2.5, 2) => [1.6667, 10, 5.8333]

multiplexed_multicarrier_bitrate(50, 3, 250*2.5, 1) => [4.1667, 25, 14.5833]
multiplexed_multicarrier_bitrate(50, 3, 500*2.5, 1) => [8.3333, 50, 29.1667]

multiplexed_multicarrier_bitrate(10, 3, 250*2.5, 1) => [20.8333, 125, 72.9167]
multiplexed_multicarrier_bitrate(10, 3, 500*2.5, 1) => [41.6667, 250, 145.8333]

multiplexed_multicarrier_bitrate(10, 3, 500*2.5, 4) => [166.6667, 1,000, 583.3333]

multiplexed_multicarrier_bitrate(10, 3, 500*2.5, 4) => [166.6667, 1,000, 583.3333]

Ludacris:

multiplexed_multicarrier_bitrate(1, 2, 500*2.5, 5) => [3,125, 12,500, 7,812.5]
multiplexed_multicarrier_bitrate(1, 2, 500*8, 5) => [10,000, 40,000, 25,000]



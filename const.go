package ads1x15

// Register addresses
const (
	ConvAddr = 0b00
	CfgAddr  = 0b01
	LoThresh = 0b10
	HiThresh = 0b11
)

// Operational Status
const (
	SingleConv uint16 = (1 << 15)
	OSMask     uint16 = (1 << 15)
)

// MuxSetting defines a multiplexer settings type.
type MuxSetting uint16

// Input Multiplexer
const (
	Mux0_1 MuxSetting = (iota << 12)
	Mux0_3
	Mux1_3
	Mux2_3
	Mux0
	Mux1
	Mux2
	Mux3

	MuxMask uint16 = (0b111 << 12)
)

// Programmable Gain Amplifier
const (
	FS6v144 uint16 = (iota << 9)
	FS4v096
	FS2v048
	FS1v024
	FS0v512
	FS0v256

	PGAMask uint16 = (0b111 << 9)
)

// Mode
const (
	ModeContinuous uint16 = (iota << 8)
	ModeSingleShot

	ModeMask uint16 = (1 << 8)
)

// Data Rate
const (
	SPS128 uint16 = (iota << 5)
	SPS250
	SPS490
	SPS920
	SPS1600
	SPS2400
	SPS3300

	DRMask uint16 = (0b111 << 5)
)

// Comparator Mode
const (
	CompTrad uint16 = (iota << 4)
	CompWindow

	CompModeMask uint16 = (1 << 4)
)

// Comparator Polarity
const (
	ActiveLow uint16 = (iota << 3)
	ActiveHigh

	CompPolMask uint16 = (1 << 3)
)

// Latch Comparator
const (
	CompNonLatch uint16 = (iota << 2)
	CompLatch

	CompLatchMask uint16 = (1 << 2)
)

// Comparator Queue
const (
	OneConv uint16 = (iota << 0)
	TwoConv
	FourConv
	DisableComp

	CompQueMask uint16 = (0b11 << 0)
)

// Device constants
const (
	Addr    = 0x48
	AddrGND = 0x48
	AddrVDD = 0x49
	AddrSDA = 0x4A
	AddrSCL = 0x4B
)

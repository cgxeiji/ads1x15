package ads1x15

import (
	"fmt"

	"github.com/cgxeiji/serial"
)

// Device defines a ADS1x15 device.
type Device struct {
	i2c *serial.I2C
}

// New returns a new ADS1x15 device.
func New(busName string, addr uint16) (*Device, error) {
	if addr == 0 {
		addr = Addr
	}

	i2c, err := serial.NewI2C(busName, addr)
	if err != nil {
		return nil, fmt.Errorf("ads1x15: could not initialize I2C: %w", err)
	}
	fmt.Println("connected")

	d := &Device{
		i2c: i2c,
	}

	// default options
	d.Options(
		DisableComparator(),
		NonLatch(),
		LowPolarity(),
		TraditionalMode(),
		SingleShot(),
		Rate(SPS1600),
		Gain(FS6v144),
	)

	return d, nil
}

// Close closes the device and cleans after itself.
func (d *Device) Close() {
	d.i2c.Close()
}

// ReadSingle returns the single ended result of the ADC on the requested
// channel.
func (d *Device) ReadSingle(channel int) (int, error) {
	if channel > 3 || channel < 0 {
		return 0, fmt.Errorf("ads1x15: channel outside the allowed range: (want: [0-3], got: %d)", channel)
	}

	var err error
	switch channel {
	case 0:
		_, err = d.Options(Mux(Mux0))
	case 1:
		_, err = d.Options(Mux(Mux1))
	case 2:
		_, err = d.Options(Mux(Mux2))
	case 3:
		_, err = d.Options(Mux(Mux3))
	}
	if err != nil {
		return 0, fmt.Errorf("ads1x15: could not select channel: %w", err)
	}

	if err := d.singleConv(); err != nil {
		return 0, fmt.Errorf("ads1x15: could not read channel %d: %w", channel, err)
	}

	if err := d.waitUntil(CfgAddr, OSMask, 1); err != nil {
		return 0, fmt.Errorf("ads1x15: could not read channel %d: %w", channel, err)
	}

	value, err := d.i2c.ReadBytes(ConvAddr, 2)
	if err != nil {
		return 0, fmt.Errorf("ads1x15: could not read channel %d: %w", channel, err)
	}

	v := (uint16(value[0])<<8 | uint16(value[1])) >> 4
	if v > 0x07FF {
		v |= 0xF000
	}

	return int(v), nil
}

// GetChannel returns a Channel struct attached to a specific multiplexer setting.
func (d *Device) GetChannel(channel MuxSetting) *Channel {
	return &Channel{
		d:   d,
		mux: channel,
	}
}

func (d *Device) singleConv() error {
	_, err := d.config(CfgAddr, OSMask, SingleConv)
	return err
}

func (d *Device) waitUntil(reg byte, flag uint16, bit byte) error {
	switch bit {
	case 1:
		for {
			state, err := d.i2c.ReadBytes(reg, 2)
			s := uint16(state[0])<<8 | uint16(state[1])
			if err != nil {
				return fmt.Errorf("could not wait for %#b in %#b to be %v", flag, reg, bit)
			} else if s&flag != 0 {
				return nil
			}
		}
	case 0:
		for {
			state, err := d.i2c.ReadBytes(reg, 2)
			s := uint16(state[0])<<8 | uint16(state[1])
			if err != nil {
				return fmt.Errorf("could not wait for %#b in %#b to be %v", flag, reg, bit)
			} else if s&flag == 0 {
				return nil
			}
		}
	}

	return fmt.Errorf("invalid bit %v, it should be 1 or 0", bit)
}

// Read reads a single byte from a register.
func (d *Device) Read(reg byte) (byte, error) {
	return d.i2c.Read(reg)
}

// ReadBytes reads n bytes from a register.
func (d *Device) ReadBytes(reg byte, n int) ([]byte, error) {
	return d.i2c.ReadBytes(reg, n)
}

// Write writes a byte or bytes to a register.
func (d *Device) Write(reg byte, data ...byte) error {
	return d.i2c.Write(reg, data...)
}

package ads1x15

import "fmt"

// Channel returns a Channel struct attached to a specific multiplexer setting.
type Channel struct {
	d   *Device
	mux MuxSetting
}

func (c *Channel) Read() (int, error) {
	old, err := c.d.Options(Mux(c.mux))
	if err != nil {
		return 0, fmt.Errorf("channel %v: could not select channel: %w", c.mux, err)
	}

	if err := c.d.singleConv(); err != nil {
		return 0, fmt.Errorf("channel %v: could not read channel: %w", c.mux, err)
	}
	if err := c.d.waitUntil(CfgAddr, OSMask, 1); err != nil {
		return 0, fmt.Errorf("channel %v: could not wait for channel: %w", c.mux, err)
	}

	value, err := c.d.i2c.ReadBytes(ConvAddr, 2)
	if err != nil {
		return 0, fmt.Errorf("channel %v: could not get value: %w", c.mux, err)
	}

	v := (uint16(value[0])<<8 | uint16(value[1])) >> 4
	if v > 0x07FF {
		v |= 0xF000
	}

	_, err = c.d.Options(old)
	if err != nil {
		return 0, fmt.Errorf("channel %v: could not return to previous configuration: %w", c.mux, err)
	}

	return int(v), nil
}

package ads1x15

import "fmt"

// Option defines a functional option for the device.
type Option func(d *Device) (Option, error)

// Options sets different configuration options and returns thte previous value
// of the last option passed.
func (d *Device) Options(options ...Option) (last Option, err error) {
	for _, opt := range options {
		last, err = opt(d)
		if err != nil {
			return nil, err
		}
	}

	return last, nil
}

func (d *Device) config(reg byte, mask, flag uint16) (uint16, error) {
	cfg, err := d.ReadBytes(reg, 2)
	if err != nil {
		return 0, fmt.Errorf("could not get %#b from %#x: %w", mask, reg, err)
	}
	c := uint16(cfg[0])<<8 | uint16(cfg[1])

	old := c & mask
	c &= ^mask
	flag = flag & mask
	c |= flag
	if err := d.Write(reg, byte(c>>8), byte(c&0xFF)); err != nil {
		return 0, fmt.Errorf("could not set %#b from %#x: %w", flag, reg, err)
	}

	return old, nil
}

// Mux configures the input multiplexer.
func Mux(mode MuxSetting) Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, MuxMask, uint16(mode))
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not configure multiplexer %#x: %w", mode, err)
		}
		return Mux(MuxSetting(old)), nil
	}
}

// Gain configures the programmable gain amplifier.
func Gain(mode uint16) Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, PGAMask, mode)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not configure gain amplifier %#x: %w", mode, err)
		}
		return Gain(old), nil
	}
}

// Rate configures the data rate.
func Rate(mode uint16) Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, DRMask, mode)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not configure data rate %#x: %w", mode, err)
		}
		return Rate(old), nil
	}
}

// TraditionalMode sets the comparator to traditional mode.
func TraditionalMode() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, CompModeMask, CompTrad)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set comparator mode to 'traditional': %w", err)
		}
		switch old {
		case CompTrad:
			return TraditionalMode(), nil
		case CompWindow:
			return WindowMode(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous comparator mode: got %#x", old)
	}
}

// WindowMode sets the comparator to traditional mode.
func WindowMode() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, CompModeMask, CompWindow)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set comparator mode to 'window': %w", err)
		}
		switch old {
		case CompTrad:
			return TraditionalMode(), nil
		case CompWindow:
			return WindowMode(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous comparator mode: got %#x", old)
	}
}

// LowPolarity sets the comparator to traditional mode.
func LowPolarity() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, CompPolMask, ActiveLow)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set comparator polarity to 'active low': %w", err)
		}
		switch old {
		case ActiveLow:
			return LowPolarity(), nil
		case ActiveHigh:
			return HighPolarity(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous comparator polarity: got %#x", old)
	}
}

// HighPolarity sets the comparator to traditional mode.
func HighPolarity() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, CompPolMask, ActiveHigh)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set comparator polarity to 'active low': %w", err)
		}
		switch old {
		case ActiveLow:
			return LowPolarity(), nil
		case ActiveHigh:
			return HighPolarity(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous comparator polarity: got %#x", old)
	}
}

// NonLatch sets the comparator to traditional mode.
func NonLatch() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, CompLatchMask, CompNonLatch)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set comparator latching mode to 'non-latching': %w", err)
		}
		switch old {
		case CompNonLatch:
			return NonLatch(), nil
		case CompLatch:
			return Latch(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous comparator latching mode: got %#x", old)
	}
}

// Latch sets the comparator to traditional mode.
func Latch() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, CompLatchMask, CompLatch)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set comparator latching mode to 'latching': %w", err)
		}
		switch old {
		case CompNonLatch:
			return NonLatch(), nil
		case CompLatch:
			return Latch(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous comparator latching mode: got %#x", old)
	}
}

// Queue configures the comparator queue.
func Queue(mode uint16) Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, CompQueMask, mode)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not configure comparator queue %#x: %w", mode, err)
		}
		return Queue(old), nil
	}
}

// DisableComparator disables the comparator.
func DisableComparator() Option {
	return func(d *Device) (Option, error) {
		old, err := Queue(DisableComp)(d)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not configure disable comparator: %w", err)
		}
		return old, nil
	}
}

// SingleShot sets the measurement mode to single shot.
func SingleShot() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, ModeMask, ModeSingleShot)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set mode to 'single shot': %w", err)
		}
		switch old {
		case ModeSingleShot:
			return SingleShot(), nil
		case ModeContinuous:
			return Continuous(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous measurement mode: got %#x", old)
	}
}

// Continuous sets the measurement mode to single shot.
func Continuous() Option {
	return func(d *Device) (Option, error) {
		old, err := d.config(CfgAddr, ModeMask, ModeContinuous)
		if err != nil {
			return nil, fmt.Errorf("ads1x15: could not set mode to 'continuous': %w", err)
		}
		switch old {
		case ModeSingleShot:
			return SingleShot(), nil
		case ModeContinuous:
			return Continuous(), nil
		}
		return nil, fmt.Errorf("ads1x15: invalid previous measurement mode: got %#x", old)
	}
}

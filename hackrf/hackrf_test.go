package hackrf

import (
	"testing"
	"time"
)

func testVer(dev *Device, t *testing.T) {
	ver, err := dev.Version()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Version: %s\n", ver)
}

func testRx(dev *Device, t *testing.T) {
	total := 0
	if err := dev.StartRX(func(buf []byte) error {
		total += len(buf)
		t.Logf("Rx: %d", len(buf))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if err := dev.StopRX(); err != nil {
		t.Fatal(err)
	}
	t.Logf("Rx total: %d bytes\n", total)
}

func TestHackRF(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Exit(); err != nil {
			t.Error(err)
		}
	}()

	var devs = make([]*DeviceInfo, 0)

	t.Run("TestDeviceList", func(t *testing.T) {
		var err error
		if devs, err = DeviceList(); err != nil {
			t.Fatal(err)
			return
		}
		for i, dev := range devs {
			t.Logf("Device %d: %s\n", i, dev.SerialNumber)
		}
	})

	if len(devs) == 0 {
		t.Skip("hackrf not found")
		return
	}

	t.Run("TestOpen", func(t *testing.T) {
		dev, err := Open()
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err = dev.Close(); err != nil {
				t.Error(err)
			}
		}()

		testVer(dev, t)
		testRx(dev, t)
	})

	t.Run("TestOpenBySerial", func(t *testing.T) {
		dev, err := OpenBySerial(devs[0].SerialNumber)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err = dev.Close(); err != nil {
				t.Error(err)
			}
		}()

		testVer(dev, t)
		testRx(dev, t)
	})
}

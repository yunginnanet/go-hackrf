package hackrf

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var openMutex sync.Mutex

func openHackRF(t *testing.T) *Device {
	openMutex.Lock()
	defer openMutex.Unlock()
	t.Helper()
	dev, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	return dev
}

func closeHackRF(dev *Device, t *testing.T) {
	openMutex.Lock()
	defer openMutex.Unlock()
	t.Helper()
	if err := dev.Close(); err != nil {
		t.Fatal(err)
	}
	return
}

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

func devList(t *testing.T) []*DeviceInfo {
	t.Helper()
	time.Sleep(500 * time.Millisecond)
	devs, err := DeviceList()
	if err != nil {
		t.Fatal(err)
	}
	return devs
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

	t.Run("TestGetVersion", func(t *testing.T) {
		dev := openHackRF(t)
		testVer(dev, t)
		closeHackRF(dev, t)
	})

	t.Run("Setup", func(t *testing.T) {
		dev := openHackRF(t)
		defer closeHackRF(dev, t)
		t.Run("TestSetSampleRate", func(t *testing.T) {
			if err := dev.SetSampleRate(10e6); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("TestSetFrequency", func(t *testing.T) {
			if err := dev.SetFreq(100e6); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("TestSetLNAGain", func(t *testing.T) {
			if err := dev.SetLNAGain(14); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("TestSetVGAGain", func(t *testing.T) {
			if err := dev.SetVGAGain(14); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("TestSetAmpEnable", func(t *testing.T) {
			if err := dev.SetAmpEnable(false); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("TestGetTransferBufferSize", func(t *testing.T) {
			size, err := dev.GetTransferBufferSize()
			if err != nil {
				t.Error(fmt.Errorf("GetTransferBufferSize: %v", err).Error())
			}
			t.Logf("Transfer buffer size: %d\n", size)
		})
		t.Run("TestRX", func(t *testing.T) {
			testRx(dev, t)
		})
	})

	t.Run("TestOpenBySerial", func(t *testing.T) {
		t.Run("TestDeviceList", func(t *testing.T) {
			devs := devList(t)
			for i, dev := range devs {
				t.Logf("Device %d: %s\n", i, dev.SerialNumber)
			}
		})

		time.Sleep(500 * time.Millisecond)

		devs := devList(t)

		dev, err := OpenBySerial(devs[0].SerialNumber)
		if err != nil {
			t.Fatal(err)
		}

		testVer(dev, t)

		if cErr := dev.Close(); cErr != nil {
			t.Error(cErr)
		}
	})
}

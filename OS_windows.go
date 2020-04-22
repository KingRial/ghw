package ghw

import (
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows/registry"
)

const wqlGetComputerInfo = "SELECT BuildNumber, BuildType, InstallDate, Manufacturer, Name, OSArchitecture, WindowsDirectory, SerialNumber, Version FROM Win32_OperatingSystem"

type win32GetComputerInfo struct {
	BuildNumber      string
	BuildType        string
	InstallDate      time.Time
	Manufacturer     string
	Name             string
	OSArchitecture   string
	WindowsDirectory string
	SerialNumber     string
	Version          string
}

func (ctx *context) osFillInfo(info *OSInfo) error {
	var win32GetComputerInfoDescriptions []win32GetComputerInfo
	if err := wmi.Query(wqlGetComputerInfo, &win32GetComputerInfoDescriptions); err != nil {
		return err
	}
	// Filling Info
	info.Name = strings.TrimSpace(strings.Split(win32GetComputerInfoDescriptions[0].Name, "|")[0])
	info.Version = win32GetComputerInfoDescriptions[0].Version
	info.Architecture = win32GetComputerInfoDescriptions[0].OSArchitecture
	info.Serial = win32GetComputerInfoDescriptions[0].SerialNumber
	info.License = "unknown"
	if err := fillLicense(info); err != nil {
		return err
	}
	return nil
}

/*
Tested On:
| OS                                        | Version    |
| ------------------------------------------|------------|
| Microsoft Windows 10 Pro                  | 10.0.18362 |
| Microsoft Windows Server 2008 R2 Standard | 6.1.7601   |
*/
func fillLicense(info *OSInfo) error {
	// Opening main key to find DigitalProductId
	mainKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return err
	}
	defer mainKey.Close()
	digitalProductID, _, _ := mainKey.GetBinaryValue("DigitalProductId")
	/*
		splittedVersion := strings.Split(info.Version, ".")
		major, err := strconv.Atoi(splittedVersion[0])
		if err != nil {
			return err
		}
		minor, err := strconv.Atoi(splittedVersion[1])
		if err != nil {
			return err
		}
		isWin8OrUp := (major > 6 || (major == 6 && minor >= 2))
	*/
	info.License = decodeProductKeyWin8AndUp(digitalProductID)
	/*if isWin8OrUp {
		info.License = decodeProductKeyWin8AndUp(digitalProductID)
	} else {
		info.License = decodeProductKeyUpToWindowsSeven(digitalProductID)
	}*/
	return nil
}

// See: https://github.com/mrpeardotnet/WinProdKeyFinder/blob/master/WinProdKeyFind/KeyDecoder.cs#L115
func decodeProductKeyWin8AndUp(key []byte) string {
	productKey := ""
	const alphabet = "BCDFGHJKMPQRTVWXY2346789"
	const keyOffset = 52
	// Check if OS is Windows 8
	isWin8 := byte((key[66] / 6) & 1)
	key[66] = byte((key[66] & 0xf7) | (isWin8&2)*4)
	last := 0
	for i := 24; i >= 0; i-- {
		current := 0
		for j := 14; j >= 0; j-- {
			current = current * 256
			current = int(key[j+keyOffset]) + current
			key[j+keyOffset] = byte(current / 24)
			current = current % 24
		}
		last = current
		productKey = string(alphabet[current]) + productKey
	}
	// Handling the special "N" character
	if isWin8 == byte(1) {
		if last > 0 { // Inserting in "last" position
			keypart1 := productKey[1 : last+1]
			keypart2 := productKey[last+1:]
			productKey = keypart1 + "N" + keypart2
		} else if last == 0 { // Removing the first character and using N
			productKey = "N" + productKey[1:]
		}
	}
	// Adding "-" separator
	productKey = productKey[0:5] + "-" + productKey[5:10] + "-" + productKey[10:15] + "-" + productKey[15:20] + "-" + productKey[20:]
	return productKey
}

/*// See: https://github.com/mrpeardotnet/WinProdKeyFinder/blob/master/WinProdKeyFind/KeyDecoder.cs#L69
func decodeProductKeyUpToWindowsSeven(key []byte) string {
	const keyStartIndex = 52
	const keyEndIndex = keyStartIndex + 15
	const digits = "BCDFGHJKMPQRTVWXY2346789"
	const decodeLength = 29
	const decodeStringLength = 15
	productKey := make([]rune, decodeLength)
	key = key[keyStartIndex : keyEndIndex+1]
	for i := decodeLength - 1; i >= 0; i-- {
		// Every sixth char is a separator.
		if (i+1)%6 == 0 {
			productKey[i] = rune('-')
		} else {
			// Do the actual decoding.
			var digitMapIndex = 0
			for j := decodeStringLength - 1; j >= 0; j-- {
				var byteValue = byte(digitMapIndex<<8) | key[j]
				key[j] = byteValue / 24
				digitMapIndex = int(byteValue % 24)
				productKey[i] = rune(digits[digitMapIndex])
			}
		}
	}
	return string(productKey)
}*/

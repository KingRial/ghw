package ghw

import (
	"strings"
	"time"

	"github.com/StackExchange/wmi"
)

const wqlGetComputerInfo = "SELECT BuildNumber, BuildType, InstallDate, Manufacturer, Name, OSArchitecture, WindowsDirectory, SerialNumber FROM Win32_OperatingSystem"

type win32GetComputerInfo struct {
	BuildNumber      string
	BuildType        string
	InstallDate      time.Time
	Manufacturer     string
	Name             string
	OSArchitecture   string
	WindowsDirectory string
	SerialNumber     string
	//Version          string
}

const wqlLicenseInfo = "SELECT * from SoftwareLicensingService"

type softwareLicensingService struct {
	OA3xOriginalProductKey            string
	OA3xOriginalProductKeyDescription string
	OA3xOriginalProductKeyPkPn        string
	Version                           string
}

func (ctx *context) osFillInfo(info *OSInfo) error {
	var win32GetComputerInfoDescriptions []win32GetComputerInfo
	if err := wmi.Query(wqlGetComputerInfo, &win32GetComputerInfoDescriptions); err != nil {
		return err
	}
	var softwareLicensingServiceDescriptions []softwareLicensingService
	if err := wmi.Query(wqlLicenseInfo, &softwareLicensingServiceDescriptions); err != nil {
		return err
	}
	// Filling Info
	info.Name = strings.TrimSpace(strings.Split(win32GetComputerInfoDescriptions[0].Name, "|")[0])
	info.Version = softwareLicensingServiceDescriptions[0].Version
	info.Architecture = win32GetComputerInfoDescriptions[0].OSArchitecture
	info.Serial = win32GetComputerInfoDescriptions[0].SerialNumber
	info.License = softwareLicensingServiceDescriptions[0].OA3xOriginalProductKey
	return nil
}

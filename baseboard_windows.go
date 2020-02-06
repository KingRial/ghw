// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import "github.com/StackExchange/wmi"

type win32_BaseBoard struct {
	Manufacturer string
	SerialNumber string
	Tag          string
	Version      string
}

func (ctx *context) baseboardFillInfo(info *BaseboardInfo) error {
	// Getting data from WMI
	var win32BaseboardDescriptions []win32_BaseBoard
	q1 := wmi.CreateQuery(&win32BaseboardDescriptions, "")
	if err := wmi.Query(q1, &win32BaseboardDescriptions); err != nil {
		return err
	}
	if len(win32BaseboardDescriptions) > 0 {
		info.AssetTag = win32BaseboardDescriptions[0].Manufacturer
		info.SerialNumber = win32BaseboardDescriptions[0].Manufacturer
		info.Vendor = win32BaseboardDescriptions[0].Manufacturer
		info.Version = win32BaseboardDescriptions[0].Manufacturer
	}

	return nil
}
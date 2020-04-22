//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
)

type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Serial       string `json:"serial"`
	License      string `json:"license"`
}

func (o *OSInfo) String() string {
	return fmt.Sprintf(
		"%s %s (Serial: %q, Product Key: %q)",
		o.Name,
		o.Version,
		o.Serial,
		o.License,
	)
}

// under a top-level "net:" key
func (o *OSInfo) YAMLString() string {
	return safeYAML(o)
}

// JSONString returns a string with the net information formatted as JSON
// under a top-level "net:" key
func (o *OSInfo) JSONString(indent bool) string {
	return safeJSON(o, indent)
}

// OS returns a pointer to a OSInfo collection containing information
// about the host's installed OS
func OS(opts ...*WithOption) (*OSInfo, error) {
	mergeOpts := mergeOptions(opts...)
	ctx := &context{
		chroot: *mergeOpts.Chroot,
	}
	info := &OSInfo{}
	if err := ctx.osFillInfo(info); err != nil {
		return nil, err
	}
	return info, nil
}

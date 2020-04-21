//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package commands

import (
	"github.com/jaypipes/ghw"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// osCmd represents the install command
var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Show installed OS information for the host system",
	RunE:  showOS,
}

// showMemory show memory information for the host system.
func showOS(cmd *cobra.Command, args []string) error {
	packages, err := ghw.OS()
	if err != nil {
		return errors.Wrap(err, "error getting OS info")
	}

	printInfo(packages)
	return nil
}

func init() {
	rootCmd.AddCommand(osCmd)
}

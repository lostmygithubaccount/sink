/*
 */
package cmd

// imports
import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "interact with the config file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// init
func init() {
	rootCmd.AddCommand(configCmd)
}

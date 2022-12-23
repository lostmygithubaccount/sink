/*
 */
package cmd

// imports
import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewConfigCmd represents the config view command
var viewConfigCmd = &cobra.Command{
	Use:   "view",
	Short: "view the config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		// print out the config file used
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			configFile = "not found"
		}
		fmt.Println("Config file:", configFile)

		// print the config
		fmt.Println("Config values:")
		keys := viper.AllKeys()
		sort.Slice(keys, func(i, j int) bool {
			return len(keys[i]) < len(keys[j])
		})
		for _, key := range keys {
			fmt.Printf("%*s: %v\n", 20, key, viper.Get(key))
		}
	},
}

func init() {
	configCmd.AddCommand(viewConfigCmd)
}

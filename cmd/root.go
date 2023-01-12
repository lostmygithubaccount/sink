/*
 */
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sink",
	Short: `A tool to help keep dbt OSS repositories in sync`,
	Long: `A tool to help keep dbt OSS repositories in sync.

Example usage:

	$ sink config view

	$ sink issue 6147 \
		--org dbt-labs \
		--source-repo dbt-core \
		--target-repos dbt-bigquery,dbt-snowflake,dbt-redshift,dbt-spark

	$ sink labels \
		--org dbt-labs \
		--source-repo dbt-core \
		--target-repos dbt-server

	$ sink labels \
		--org dbt-labs \
		--source-repo dbt-core \
		--team Core \
		--exclude-team-repos core-team,"schemas.getdbt.com" \
		--extra-repos dbt-server

It's recommended to use a config file instead of setting flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		// cmd.Help()
		// check for error
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// precedence: CLI flags > config file > defaults
// in the root command, set the defaults and populate viper with the config file
// in subcommands, allow the user to override the config file with CLI flags
func init() {
	// setup defaults
	viper.SetDefault("org", "")
	viper.SetDefault("team", "")
	viper.SetDefault("source-repo", "")
	viper.SetDefault("target-repos", []string{})
	viper.SetDefault("exclude-team-repos", []string{})
	viper.SetDefault("extra-repos", []string{})
	viper.SetDefault("dry-run", true)

	// setup config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			log.Println("Error reading config file. Exiting...")
			return
		}
	}
}

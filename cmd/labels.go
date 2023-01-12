/*
 */
package cmd

// imports
import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cli/go-gh"

	"sink/pkg/utils"
)

// labelsCmd represents the labels command
var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Sync labels from a source repository to many target repositories in the same organization",
	// TODO: add long description with examples
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// bind flags to viper
		viper.BindPFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {

		// setup the GitHub REST API client
		client, err := gh.RESTClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		// get the required config values
		org := viper.GetString("org")
		if org == "" {
			log.Fatal("org is required")
		}

		source := viper.GetString("source-repo")
		if source == "" {
			log.Fatal("source-repo is required")
		}

		// check if dry run
		dryRun := viper.GetBool("dry-run")
		if dryRun {
			log.Println("dry run...")
			defer log.Println("dry run completed, rerun with --dry-run=false to sync")
		}

		// get the options for target repos
		team := viper.GetString("team")
		targetRepos := viper.GetStringSlice("target-repos")
		excludeTeamRepos := viper.GetStringSlice("exclude-team-repos")
		extraRepos := viper.GetStringSlice("extra-repos")

		// convert strings to repos
		sourceRepo := utils.Repo{Name: source}

		// convert repos strings to Repo structs
		targetReposRepos := utils.ReposToRepos(targetRepos)
		excludeTeamReposRepos := utils.ReposToRepos(excludeTeamRepos)
		extraReposRepos := utils.ReposToRepos(extraRepos)

		// get the user TODO: unused, remove?
		user := utils.GetUser(client)

		// compute the target repos
		targetReposReposRepos := utils.GetTargetRepos(org, team, targetReposRepos, excludeTeamReposRepos, extraReposRepos, client)

		// debugging TODO: debug flag? remove?
		log.Println("----------------------------------------------------------------")
		log.Println("user:", user)
		log.Println("org:", org)
		log.Println("sourceRepo:", sourceRepo)
		log.Println("targetReposRepos:", targetReposReposRepos)
		log.Println("excludeTeamReposRepos:", excludeTeamReposRepos)
		log.Println("extraReposRepos:", extraReposRepos)
		log.Println("----------------------------------------------------------------")

		sourceLabels := utils.GetRepoLabels(org, sourceRepo, client)
		log.Println("sourceLabels:", sourceLabels)

		// for each target repo
		for _, target := range targetReposReposRepos {
			log.Println("\ttarget:", target)
			targetLabels := utils.GetRepoLabels(org, target, client)
			log.Println("\ttargetLabels:", targetLabels)
			log.Println("\tsyncing labels...")
			utils.SyncLabels(dryRun, org, sourceRepo, target, sourceLabels, targetLabels, client)
		}
	},
}

// init
func init() {

	// add the labels command to the root command
	rootCmd.AddCommand(labelsCmd)

	// create flags
	labelsCmd.Flags().StringP("org", "o", "", "GitHub organization")
	labelsCmd.Flags().StringP("team", "t", "", "team name")
	labelsCmd.Flags().StringP("source-repo", "s", "", "source repository")
	labelsCmd.Flags().StringSliceP("target-repos", "T", []string{}, "target repositories")
	labelsCmd.Flags().StringSliceP("exclude-team-repos", "e", []string{}, "exclude repositories")
	labelsCmd.Flags().StringSliceP("extra-repos", "E", []string{}, "extra repositories")
	labelsCmd.Flags().Bool("dry-run", true, "dry run")
}

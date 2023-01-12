/*
 */
package cmd

// imports
import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cli/go-gh"

	"sink/pkg/utils"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue [issue number]",
	Short: "Sync an issue from a source repository to many target repositories in the same organization",
	// TODO: add long description with examples
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Args: cobra.ExactArgs(1),
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

		// parse the issue number
		issueNumber, err := strconv.Atoi(args[0])
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

		dryRun := viper.GetBool("dry-run")
		if dryRun {
			log.Println("dry run...")
			defer log.Println("dry run completed, rerun with --dry-run=false to copy issues")
		}

		// get the options for target repos
		team := viper.GetString("team")
		targetRepos := viper.GetStringSlice("target-repos")
		excludeTeamRepos := viper.GetStringSlice("exclude-team-repos")
		extraRepos := viper.GetStringSlice("extra-repos")

		// convert the strings to repos
		sourceRepo := utils.Repo{Name: source}

		// convert the repos to repos
		targetReposRepos := utils.ReposToRepos(targetRepos)
		excludeTeamReposRepos := utils.ReposToRepos(excludeTeamRepos)
		extraReposRepos := utils.ReposToRepos(extraRepos)

		// get the user TODO: unused, remove?
		user := utils.GetUser(client)

		// compute the target repo
		targetReposReposRepos := utils.GetTargetRepos(org, team, targetReposRepos, excludeTeamReposRepos, extraReposRepos, client)

		// get the source issue
		sourceIssue := utils.GetIssue(org, sourceRepo, issueNumber, client)

		// debugging TODO: debug flag? remove?
		log.Println("----------------------------------------------------------------")
		log.Println("user:", user)
		log.Println("org:", org)
		log.Println("source-repo:", source)
		log.Println("sourceIssue:", sourceIssue)
		log.Println("targetRepos:", targetReposReposRepos)
		log.Println("excludeTeamRepos:", excludeTeamReposRepos)
		log.Println("extraRepos:", extraReposRepos)
		log.Println("----------------------------------------------------------------")

		// for each target repo
		for _, target := range targetReposReposRepos {
			log.Println("\ttarget:", target)
			if !dryRun {
				log.Println("\tcopying issue...")
				targetIssue := utils.CopyIssue(org, target, sourceIssue, client)
				log.Println("\ttargetIssue:", targetIssue)
			}
		}
	},
}

// init
func init() {

	// add the issues command to the root command
	rootCmd.AddCommand(issueCmd)

	// create flags
	issueCmd.Flags().StringP("org", "o", "", "GitHub organization")
	issueCmd.Flags().StringP("team", "t", "", "team name")
	issueCmd.Flags().StringP("source-repo", "s", "", "source repository")
	issueCmd.Flags().StringSliceP("target-repos", "T", []string{""}, "target repositories")
	issueCmd.Flags().StringSliceP("exclude-team-repos", "e", []string{}, "exclude repositories")
	issueCmd.Flags().StringSliceP("extra-repos", "E", []string{}, "extra repositories")
	issueCmd.Flags().Bool("dry-run", true, "dry run")
}

package cmd

import (
	"errors"
	"fmt"
	"github.com/gernest/front"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pjeziorowski/rollout/platforms"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "rollout",
	Short: "Distributes your markdown across multiple blogging platforms",
	Run: func(cmd *cobra.Command, args []string) {
		markdownFile := args[0]
		Publish(markdownFile)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a markdown file path argument. ")
		}
		return nil
	},
}

var hashnodeApiToken = ""
var devtoApiToken = ""
var mediumApiToken = ""
var hashnodePublicationId = ""
var mediumPublicationId = ""

func Publish(markdownFile string) {
	hashnode := platforms.NewHashnode(hashnodeApiToken, hashnodePublicationId)
	medium := platforms.NewMedium(mediumPublicationId, mediumApiToken)
	devto := platforms.NewDevto(devtoApiToken)

	p := []platforms.Platform{hashnode, medium, devto}

	bytes, err := ioutil.ReadFile(markdownFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	fileContent := string(bytes)

	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	frontmatter, body, err := m.Parse(strings.NewReader(fileContent))
	if err != nil {
		log.Fatal(err.Error())
	}

	markdown := body
	tags := cast.ToStringSlice(frontmatter["tags"])
	if tags == nil {
		log.Fatal("Forgot tags in markdown frontmatter? ")
	}
	title := cast.ToString(frontmatter["title"])
	if title == "" {
		log.Fatal("Forgot title in markdown frontmatter? ")
	}
	canonicalUrl := cast.ToString(frontmatter["canonical_url"])
	if title == "" {
		log.Fatal("Forgot canonical URL in markdown frontmatter? ")
	}

	for _, platform := range p {
		platform.Publish(title, markdown, tags, canonicalUrl)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&hashnodeApiToken, "HASHNODE_API_TOKEN", "", "")
	rootCmd.PersistentFlags().StringVar(&mediumApiToken, "MEDIUM_API_TOKEN", "", "")
	rootCmd.PersistentFlags().StringVar(&devtoApiToken, "DEVTO_API_TOKEN", "", "")
	rootCmd.PersistentFlags().StringVar(&hashnodePublicationId, "HASHNODE_PUBLICATION_ID", "", "")
	rootCmd.PersistentFlags().StringVar(&mediumPublicationId, "MEDIUM_PUBLICATION_ID", "", "")
	err := rootCmd.MarkPersistentFlagRequired("HASHNODE_API_TOKEN")
	err = rootCmd.MarkPersistentFlagRequired("MEDIUM_API_TOKEN")
	err = rootCmd.MarkPersistentFlagRequired("DEVTO_API_TOKEN")
	err = rootCmd.MarkPersistentFlagRequired("HASHNODE_PUBLICATION_ID")
	err = rootCmd.MarkPersistentFlagRequired("MEDIUM_PUBLICATION_ID")
	if err != nil {
		log.Fatal(err.Error())
	}
	cobra.OnInitialize(initConfig)
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".rollout" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".rollout")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		presetRequiredFlags(rootCmd)
	} else {
		fmt.Println("Please, provide a rollout config file in your $HOME directory: " + os.Getenv("HOME") + "/.rollout.yaml")
		fmt.Println("The file should contain configuration of platforms to publish your content to and be structured as follows:")
		fmt.Println("\n" +
			"HASHNODE_API_TOKEN: \"******\"\n" +
			"HASHNODE_PUBLICATION_ID: \"******\"\n" +
			"MEDIUM_API_TOKEN: \"******\"\n" +
			"MEDIUM_PUBLICATION_ID: \"******\"\n" +
			"DEVTO_API_TOKEN: \"******\"" +
			"")
		os.Exit(1)
	}
}

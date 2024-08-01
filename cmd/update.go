/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FrameworkConfig struct {
	Yaml struct {
		URL      string `mapstructure:"url"`
		FilePath string `mapstructure:"filePath"`
	} `mapstructure:"yaml"`
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update command is get swagger yaml from git repo.",
	Long:  `update command is get swagger yaml from git repo.`,
	Run:   updateYamlFromGit,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateYamlFromGit(cmd *cobra.Command, args []string) {
	configPath := filepath.Join(".", "data", "git-swagger-yaml-raw-url.yaml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf(Red+"Error reading config file, %s"+Reset, err)
	}

	configs := viper.AllSettings()

	for framework, _ := range configs {
		var fc FrameworkConfig
		err := viper.Sub(framework).Unmarshal(&fc)
		if err != nil {
			log.Fatalf(Red+"Unable to decode conf yaml, %v"+Reset, err)
		}
		fmt.Printf(Green+" > Update [%s] swagger.yaml from %s\n"+Reset, framework, fc.Yaml.URL)
		err = downloadFile(fc.Yaml.URL, fc.Yaml.FilePath)
		if err != nil {
			log.Printf(Red+"Unable to update [%s].swagger.yaml, %v"+Reset, framework, err)
		}
	}

}

func downloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

/*
Copyright © 2024 raccoon-mh
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert command is convert swagger yaml to csv file.",
	Long:  `convert command is convert swagger yaml to csv file.`,
	Run:   convertYaml,
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

func convertYaml(cmd *cobra.Command, args []string) {
	configPath := filepath.Join(".", "data", "git-swagger-yaml-raw-url.yaml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	fwInfos := viper.AllSettings()

	for framework, _ := range fwInfos {
		swaggerPath := filepath.Join(".", "data/basedata", framework+".swagger.yaml")
		viper.SetConfigFile(swaggerPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Println(Red+"Error reading config file, %s"+Reset, err)
		}
		swagger := viper.AllSettings()

		data := [][]string{
			{"tag", "opertaionID", "summary", "method", "path"},
		}

		for rootIndex, rootIndexvalue := range swagger {
			if rootIndex == "paths" {
				for path, pathData := range rootIndexvalue.(map[string]interface{}) {
					for method, apiSpec := range pathData.(map[string]interface{}) {
						row := make([]string, 5)
						row[3] = method
						row[4] = path
						for data, value := range apiSpec.(map[string]interface{}) {
							if data == "tags" {
								row[0] = convertAndJoinTags(value)
							} else if data == "operationid" {
								row[1] = value.(string)
							} else if data == "summary" {
								row[2] = value.(string)
							}
						}
						data = append(data, row)
					}
				}
			}
		}

		now := time.Now()
		formattedTime := now.Format("20060102.150405.")

		file, err := os.Create("./data/converted/" + formattedTime + framework + ".swagger.csv")
		if err != nil {
			log.Printf(Red+"CSV 파일을 생성하는데 실패했습니다: %s\n"+Reset, err.Error())
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		for _, record := range data {
			if err := writer.Write(record); err != nil {
				log.Printf(Red+"CSV 파일에 데이터를 기록하는데 실패했습니다: %s\n"+Reset, err.Error())
			}
		}

	}
}

func convertAndJoinTags(tags interface{}) string {
	var strTags []string
	switch v := tags.(type) {
	case []interface{}:
		for _, tag := range v {
			strTags = append(strTags, fmt.Sprintf("%v", tag))
		}
	default:
		fmt.Println("Unexpected type:", reflect.TypeOf(tags))
		return ""
	}
	return strings.Join(strTags, ", ")
}

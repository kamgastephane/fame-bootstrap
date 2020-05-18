/*
Copyright © 2020 NAME HERE <kamga.stephane@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)
var entity string
// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap -e [the name of an entity]",
	Short: "bootstrap an entity or a list of entity from fame",
	Long: `Bootstrap an entity of a list of entity from fame, 
		on the requested environment for the requested service bus queue`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		if len(url) == 0 {
			log.Fatalln("The Url parameter was not found in the configuration")
		}
		fameService := service{Url:url}
		entities, err := fameService.fetchEntities(entity)
		if err != nil {
			os.Exit(-1)
		}
		code, err := selectPrompt("Select the entity:", entities)
		if err != nil {
			log.Fatalln("An error occurred while selecting the entity")
		}
		fmt.Printf("bootstrapping entity with name %s and code=%s", entity, code)
		queue, err := inputPrompt("Enter the name of the destination queue:")
		if err != nil {
			log.Fatalln("invalid queue name")
		}
		taskId, err := fameService.createTask(queue, code)
		if err != nil {
			log.Fatalln("Failed to enqueue a bootstrap task")
		}
		err = fameService.dequeue(taskId)
		if err != nil {
			log.Fatalln("Failed to dequeue a bootstrap task")
		}
	},
}

func inputPrompt(question string) (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    question,
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}
func selectPrompt(question string, options []string) (string, error){
	var err error
	var result string
	for ok :=true; ok; ok = err != nil {
		prompt := promptui.Select{
			Label: question,
			Items: options,
		}
		_, result, err = prompt.Run()

	}
	return result, nil
}
func init() {
	rootCmd.AddCommand(bootstrapCmd)

	// Here you will define your flags and configuration settings.

	bootstrapCmd.Flags().StringVarP(&entity, "entity", "e", "", "The entity to bootstrap e.g. match (required)")
	_ = bootstrapCmd.MarkFlagRequired("entity")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bootstrapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

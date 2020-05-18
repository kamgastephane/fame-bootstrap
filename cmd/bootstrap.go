/*
Copyright © 2020 Kamga Stephane <kamga.stephane@gmail.com>

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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

const defaultApplicationId = "16"
const defaultProjectId = "0"

var environment string
var url string
var queue string

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap [the name of a resource] [the ID of the resource]",
	Short: "bootstrap a resource or a list of resources from fame",
	Long: `Bootstrap a resource of a list of resources from fame, 
		on the requested environment for the requested service bus queue. If the name of the resource is not specified, you can run a search to find it`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		context := getBaseUrl(url)
		context.ApplicationId, context.ProjectId = getBootstrapContext()

		resourceName := ""
		fameService := NewService(context)
		if len(args) > 0 {
			resourceName = args[0]
		} else {
			name, err := inputPrompt("Insert the name of the resource", false)
			resourceName = name
			if err != nil {
				log.Fatalln("Invalid name for the resource to bootstrap")
			}
		}
		resourceContracts, err := fameService.Find(resourceName)
		if err != nil {
			log.Fatalln(err)
		}
		options, resources := toMap(resourceContracts)
		var selected string
		if len(options) > 1 {
			selected, err = selectPrompt("Select the resource:", options)
			if err != nil {
				log.Fatalf("An error occurred while selecting the resource. exception:%s\n", err)
			}
		} else if len(options) == 1 {
			selected = options[0]
		} else {
			log.Fatalf("No resources found with name %s\n", resourceName)
		}

		resourceClassId := resources[selected]
		if queue == "" {
			queue, err = inputPrompt("Enter the name of the destination queue:", false)
			if err != nil {
				log.Fatalln("Invalid queue name")
			}
		}
		if queue == "" {
			log.Fatalln("Invalid queue name")
		}
		var entityId string = ""
		if len(args) > 1 {
			entityId = args[1]
		} else {
			entityId, err = inputPrompt("Enter the id of the specific entity to bootstrap (or leave empty to bootstrap the entire resourceClass)", true)
		}
		fmt.Printf("creating task for resourceClass %s entityId: %s on queue %s\n", selected, entityId, queue)

		response, err := fameService.CreateTask(queue, strconv.Itoa(resourceClassId), entityId)
		if err != nil {
			log.Fatalf("failed to create bootstrap task. exception %s\n", err)
		}
		fmt.Printf("Dequeing task for resourceClass %s entityId: %s on queue %s \n", selected, entityId, queue)

		dequeueResponse, err := fameService.Dequeue(response.CommandId)
		if err != nil {
			log.Fatalf("failed to launch the dequeue task. exception %s\n", err)
		}
		fmt.Printf("bootstrap task triggered for resourceClass %s entityId: %s on queue %s. operationId: %s\n", selected, entityId, queue, dequeueResponse.ResponseStatus.Meta.OperationId)
	},
}



func toMap(input [] ResourceContracts) ([] string, map[string]int) {
	result := make(map[string]int, len(input))
	var options []string
	for _, item := range input {
		key := fmt.Sprintf("%s (%d)", item.Name, item.ResourceClassId)
		options = append(options, key)
		result[key] = item.ResourceClassId
	}
	return options, result
}

func inputPrompt(question string, allowEmpty bool) (string, error) {
	validate := func(input string) error {
		if !allowEmpty {
			if len(input) == 0 {
				return errors.New("")
			}
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
func selectPrompt(question string, options []string) (string, error) {
	var result string
	prompt := promptui.Select{
		Label: question,
		Items: options,
	}
	_, result, err := prompt.Run()
	return result, err
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	// Here you will define your flags and configuration settings.
	bootstrapCmd.Flags().StringVarP(&environment, "environment", "e", "", "The environment to run the bootstrap on")
	bootstrapCmd.Flags().StringVarP(&url, "url", "u", "", "The base url to run the bootstrap on")
	bootstrapCmd.Flags().StringVarP(&queue, "queue", "q", "", "The destination queue to push the data on")
}

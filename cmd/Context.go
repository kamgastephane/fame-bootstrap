package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type Context struct {
	BaseUrl string
	Env string
	ApplicationId string
	ProjectId string
}
func getEnvironment(input string) string {
	if input == "" {
		selectedEnvironment, err := selectPrompt("Select the environment:", []string{"STG","GOLD","PROD"})
		if err != nil {
			os.Exit(-1)
		} else {
			return selectedEnvironment
		}
	}
	return input
}

func getBootstrapContext() (string, string) {
	appId := defaultApplicationId
	pId := defaultProjectId
	applicationId := viper.GetString("url.applicationId")
	if applicationId != "" {
		appId = applicationId
	}
	projectId := viper.GetString("url.projectId")
	if projectId != "" {
		pId = projectId
	}
	return appId,pId
}

func getBaseUrl(input string) Context {
	if input == "" {
		// no base url input, we get the requested environment
		env := getEnvironment(environment)
		if env == "" {
			log.Fatal("failed to get the environment")
		}
		url := viper.GetString(fmt.Sprintf("url.%s", strings.ToLower(env)))

		if url != "" {
			return Context{
				BaseUrl: url,
				Env:     env,
			}
		} else {

			log.Fatalf("the url parameter for environment %s was not found in the configuration\n", strings.ToUpper(env))
		}
	}
	return Context{BaseUrl:input}
}


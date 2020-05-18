package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type service struct {
	Url string
}

type fameEntity struct {
	key string
}

func(s * service) fetchEntities(entityName string) ([]string, error) {
	fameEntities := make([]fameEntity, 0)

	resp, err := http.Get(s.Url)
	if err != nil{
		return nil, errors.New(fmt.Sprintf("Server error while fetching resource with entity %s", entityName))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("invalid status code while fetching resource with entity %s", entityName))
	}
	defer resp.Body.Close()

	// read the payload, in this case
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid response body while fetching resource with entity %s", entityName))
	}
	err = json.Unmarshal(body, &fameEntities)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to deserialize body while fetching resource with entity %s", entityName))
	}
	entityNames := make([]string, len(fameEntities))
	for _, val := range fameEntities{
		entityNames = append(entityNames, val.key)
	}
	return entityNames, nil
}

func(s * service) createTask(queue string, entityId string) (string, error) {

}

func(s * service) dequeue(task string) error {

}

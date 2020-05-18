package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)
var client = &http.Client{Timeout: 10 * time.Second}

type Service struct {
	Url string
	ApplicationId string
	ProjectId string
}
func NewService(ctx Context) *Service {
	return &Service{ctx.BaseUrl, ctx.ApplicationId, ctx.ProjectId}
}

type findResponse struct {
	Resources [] ResourceContracts `json:"resourceContracts"`
}
type StatusResponse struct {
	RemainingMessagesSent int `json:"remainingMessagesToBeSent"`
	FailedMessages int `json:"failedMessages"`
}
type ResourceContracts struct {
	Name string `json:"name"`
	ResourceClassId int `json:"resourceClassId"`
	Properties interface{} `json:"properties"`
}
type enqueuePayload struct {
	To string `json:"to"`
	ApplicationId string `json:"applicationId"`
	ProjectId string `json:"projectId"`
	ResourceClassId string `json:"resourceClassId"`
	ResourceId string `json:"resourceId"`
}
type EnqueueResponse struct {
	EntityCount int `json:"commandsEnqueued"`
	CommandId string `json:"bootstrapCommandId"`
}
type DequeueResponse struct {
	ResponseStatus ResponseStatus `json:"ResponseStatus"`
}
type ResponseStatus struct {
	Meta Meta `json:"meta"`
}
type Meta struct {
	OperationId string `json:"operationId"`
}
type DequeuePayload struct {
	CommandId string `json:"bootstrapCommandId"`
	Async bool `json:"async"`

}

func(s * Service) Find(entityName string) ([]ResourceContracts, error) {
	findUrl := s.Url + fmt.Sprintf("resources/contracts/find?resourceClassName=%s&format=json", entityName)
	response := findResponse{}
	resp, err := client.Get(findUrl)
	if err != nil{
		fmt.Print(err)
		return nil, errors.New(fmt.Sprintf("Server error while fetching resource with entity %s\n", entityName))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("invalid status code while fetching resource with entity %s\n", entityName))
	}
	defer resp.Body.Close()
	// read the payload, in this case
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialize body while fetching resource with entity %s, exception: %s\n", entityName, err)
	}
	return response.Resources, nil

}

func(s * Service) CreateTask(queue string, resourceClassId string, resourceId string) (*EnqueueResponse, error) {
	enqueueUrl := s.Url +
		fmt.Sprintf("resources/bootstrap/enqueue?to=%s&applicationId=%s&projectId=%s&resourceClassId=%s&resourceId=%s&format=json",
			queue, s.ApplicationId, s.ProjectId, resourceClassId, resourceId)
	response, err := http.Post(enqueueUrl, "application/json", strings.NewReader("{}"))
	if err != nil{
		return nil, fmt.Errorf("server error while enqueing resource with entity %s on queue %s\n", resourceId, queue)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code while enqueing resource with entity %s on queue %s\n", resourceId, queue)
	}
	defer response.Body.Close()
	var enqueue EnqueueResponse
	err = json.NewDecoder(response.Body).Decode(&enqueue)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialize body while enqueing resource with entity %s on queue %s. exception:%s\n", resourceId, queue, err)
	}
	return &enqueue, nil
}

func(s * Service) Dequeue(commandId string) (*DequeueResponse,error) {
	dequeueUrl := s.Url + fmt.Sprintf("resources/bootstrap/dequeue?bootstrapCommandId=%s&async=true&format=json", commandId)

	response, err := http.Post(dequeueUrl, "application/json", strings.NewReader("{}"))
	if err != nil{
		return nil, fmt.Errorf("server error while dequeing resource with commandId %s\n", commandId)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code while dequeing resource with commandId %s\n", commandId)
	}
	defer response.Body.Close()

	var dequeueResp DequeueResponse
	err = json.NewDecoder(response.Body).Decode(&dequeueResp)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialize body while dequeing resource with commandId %s\n", commandId)
	}
	return &dequeueResp, nil
}

func(s * Service) Status(commandId string) (*StatusResponse,error) {
	url := s.Url + fmt.Sprintf("resources/bootstrap/%s/status?format=json", commandId)
	response := StatusResponse{}
	resp, err := client.Get(url)
	if err != nil{
		fmt.Print(err)
		return nil, errors.New(fmt.Sprintf("Server error while fetching status with commandId %s\n", commandId))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("invalid status code while fetching status with command %s\n", commandId))
	}
	defer resp.Body.Close()
	// read the payload, in this case
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialize body while fetching status with commandId %s, exception: %s\n", commandId, err)
	}
	return &response, nil
}
package tooling

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/owid-metadata/common"
	"github.com/margostino/owid-metadata/configuration"
	"github.com/margostino/owid-metadata/model"
	"github.com/margostino/owid-metadata/utils"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func GenerateMetadata() {
	config := configuration.GetConfig()
	var accessToken = config.GithubAccessToken
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	_, datasets, _, err := client.Repositories.GetContents(context.Background(), "owid", "owid-datasets", "datasets", nil)
	common.Check(err)
	for _, dataset := range datasets {
		var metadata model.Metadata
		var dataPackage map[string]interface{}
		datasetName := utils.NormalizeName(*dataset.Name)
		path := *dataset.Path + "/datapackage.json"
		encodedData, _, _, err := client.Repositories.GetContents(context.Background(), "owid", "owid-datasets", path, nil)
		common.Check(err)
		content, err := encodedData.GetContent()
		common.Check(err)
		json.Unmarshal([]byte(content), &dataPackage)
		dataResources := dataPackage["resources"].([]interface{})

		metadata.Name = datasetName
		metadata.DataFile = dataResources[0].(map[string]interface{})["path"].(string) // TODO: check index
		metadata.DataBaseUrl = strings.ReplaceAll(encodedData.GetDownloadURL(), "datapackage.json", metadata.DataFile)
		metadata.Description = dataPackage["description"].(string)
		metadata.Arguments = make([]*model.Variable, 0)
		metadata.Variables = make([]*model.Variable, 0)

		fmt.Printf("\nOriginal: %s  -  Normalized: %s\n", *dataset.Name, datasetName)
		for _, resource := range dataResources {
			schema := resource.(map[string]interface{})["schema"]
			fieldsMap := schema.(map[string]interface{})
			fields := fieldsMap["fields"].([]interface{})
			for _, fieldMap := range fields {
				field := fieldMap.(map[string]interface{})
				fieldName := utils.NormalizeName(field["name"].(string))
				fieldType := utils.NormalizeType(field["type"].(string))
				variable := model.Variable{
					Name: fieldName,
					Type: fieldType,
				}
				if fieldName == "entity" || fieldName == "year" {
					metadata.Arguments = append(metadata.Arguments, &variable)
				} else if len(fieldName) > 0 && fieldName != "\r" {
					metadata.Variables = append(metadata.Variables, &variable)
				}
				fmt.Printf("FieldName: %s  -  FieldType: %s\n", fieldName, fieldType)
			}
		}

		fileName := fmt.Sprintf("%s/%s.yml", config.MetadataPath, metadata.Name)
		if len(metadata.Variables) > 0 {
			yamlData, err := yaml.Marshal(&metadata)
			common.Check(err)
			err = ioutil.WriteFile(fileName, yamlData, 0644)
			common.Check(err)
			fmt.Printf("\nNew file: %s\n", fileName)
		} else {
			fmt.Printf("\nNo variables for %s\n", fileName)
		}
	}

	// TODO: add and report sanity checks at the end
}

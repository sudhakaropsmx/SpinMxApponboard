// Copyright (c) 2019, OpsMx, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package application

import (
	"encoding/json"
	//"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	//"text/template"
	"github.com/spf13/cobra"
	"github.com/spinnaker/spin/cmd/gateclient"
	"github.com/spinnaker/spin/cmd/orca-tasks"
	"github.com/spinnaker/spin/util"
	appconfig "github.com/sudhakaropsmx/spinmx/config/application"
	"gopkg.in/yaml.v2"
)

type CreateApplicationOptions struct {
	*applicationOptions
	output                string
	applicationConfigFile string
}

var (
	createapplicationShort = "Create the application with the provided config values and then save "
	createapplicationLong  = "Create the application with the provided config values and then save "
)

func NewCreateApplicationCmd(applicationOptions applicationOptions) *cobra.Command {
	options := CreateApplicationOptions{
		applicationOptions: &applicationOptions,
	}
	cmd := &cobra.Command{
		Use:   "create",
		Short: createapplicationShort,
		Long:  createapplicationLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCreateApplication(cmd, options)
		},
	}

	cmd.PersistentFlags().StringVarP(&options.applicationConfigFile, "file", "f", "", "path to the application config values file")

	return cmd
}

func getCreateApplication(cmd *cobra.Command, options CreateApplicationOptions) error {

	gateClient, err := gateclient.NewGateClient(cmd.InheritedFlags())

	if err != nil {
		return err
	}
	applicationConfig := &appconfig.AplicationOpsConfig{}

	yamlFile, err := ioutil.ReadFile(options.applicationConfigFile)

	if yamlFile != nil {

		err = yaml.UnmarshalStrict([]byte(os.ExpandEnv(string(yamlFile))), &applicationConfig)

		if err != nil {
			util.UI.Error(fmt.Sprintf("Could not deserialize config file with contents: %s, failing.", yamlFile))
			return err
		}
	}
	//valid := trueboarding
	application := applicationConfig.Application
	for i := 0; i < len(application); i++ {
		fmt.Printf("application : %s\n", application[i].Name)
		app, resp, err := gateClient.ApplicationControllerApi.GetApplicationUsingGET(gateClient.Context, application[i].Name, map[string]interface{}{"expand": false})

		if resp.StatusCode == http.StatusNotFound && len(app) == 0 {
			app := application[i]
			permissions := strings.Split(app.Permissions, ",")

			permissionsMap := map[string][]string{
				"READ":    permissions,
				"WRITE":   permissions,
				"EXECUTE": permissions,
			}
			taksApp := map[string]interface{}{
				"accounts": app.Accounts,
				"name":    app.Name,
				"email":  app.OwnerEmail,
				"permissions": permissionsMap,
			}
			jsonString, err := json.Marshal(taksApp)
			fmt.Println(string(jsonString))

			createAppTask := map[string]interface{}{
				"job":         []interface{}{map[string]interface{}{"type": "createApplication", "application": taksApp}},
				"application": taksApp["name"],
				"description": fmt.Sprintf("Create Application: %s", taksApp["name"]),
			}
			
			ref, _, err := gateClient.TaskControllerApi.TaskUsingPOST1(gateClient.Context, createAppTask)
			if err != nil {
				return err
			}

			err = orca_tasks.WaitForSuccessfulTask(gateClient, ref, 5)
			if err != nil {
				return err
			}
			util.UI.Info(util.Colorize().Color(fmt.Sprintf("[reset][bold][green]Application save succeeded")))

		} else {
			if len(app) > 0 {
				return fmt.Errorf("Encountered an error create application, application already exist with the name :%d\n", application[i].Name)
			} else {
				return fmt.Errorf("Encountered an error querying application verify exist or not, status code: %d -- %d \n", resp.StatusCode, err)
			}
		}

	}
	return nil
}

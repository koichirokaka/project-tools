// Copyright © 2018 Okaka <koichirokaka.gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// mlCmd represents the ml command
var mlCmd = &cobra.Command{
	Use:   "ml",
	Short: "Create machine learning environment",
	Long: `Create machine learing environment.
	* environment.yml for conda
	* pylint config
	* vscode settings
	* gitignore`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		} else if name == "" {
			log.Fatal("name flag is required")
		}
		if err := create(name, mlproject); err != nil {
			log.Fatal(err)
		}
	},
}

var mlproject *project

func initml() {
	rootCmd.AddCommand(mlCmd)
	p, err := setupProject("ml")
	if err != nil {
		panic(err)
	}
	mlproject = p
}

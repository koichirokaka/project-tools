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

// mlEngineCmd represents the mlEngine command
var mlEngineCmd = &cobra.Command{
	Use:   "mlEngine",
	Short: "Create Cloud ML engine starter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := projectName(cmd)
		if err != nil {
			log.Fatal(err)
		}
		if err := create(name, mlproject); err != nil {
			log.Fatal(err)
		}
	},
}

var mlengineproject *project

func initmlengine() {
	rootCmd.AddCommand(mlEngineCmd)
	p, err := setupProject("ml-engine")
	if err != nil {
		panic(err)
	}
	mlproject = p
}

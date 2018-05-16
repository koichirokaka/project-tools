// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "project-tools",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Order is important because can't use viper config without calling initConfig.
	initConfig()
	initml()

	rootCmd.PersistentFlags().StringP("name", "n", "", "project name")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile("viper.yaml")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

type project struct {
	baseDir string
	dirs    []string
	files   map[string]*embededFile
	cmds    [][]string
}

type embededFile struct {
	content   string
	isEmbeded bool
}

func setupProject(baseDir string) (*project, error) {
	// Create emebeded map
	root := filepath.Join("template", baseDir)
	enames := viper.GetStringSlice("embeded." + baseDir)
	embeded := make(map[string]bool)
	for _, ename := range enames {
		embeded[filepath.Join(root, ename)] = true
	}
	// Walk base directory.
	p := &project{baseDir: baseDir, files: make(map[string]*embededFile)}
	existdir := make(map[string]bool)
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		realpath := strings.Split(path, root)[1][1:]
		if info.IsDir() {
			if realpath == "" {
				return nil
			}
			realpath = strings.Trim(realpath, "/")
			p.dirs = append(p.dirs, realpath)
			existdir[realpath] = true
		} else {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			dir, _ := filepath.Split(realpath)
			dir = strings.Trim(dir, "/")
			if dir != "" && !existdir[dir] {
				p.dirs = append(p.dirs, strings.Trim(dir, "/"))
			}
			// Remove root directory and slash.
			p.files[realpath] = &embededFile{string(b), embeded[path]}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	cmds := viper.GetStringSlice("cmd." + baseDir)
	for _, cmd := range cmds {
		p.cmds = append(p.cmds, strings.Split(cmd, " "))
	}
	return p, nil
}

func create(name string, p *project) error {
	if err := os.MkdirAll(name, 0755); err != nil {
		return fmt.Errorf("fail to create project directory: %v", err)
	}
	for _, dir := range p.dirs {
		if err := os.MkdirAll(filepath.Join(name, dir), 0755); err != nil {
			return fmt.Errorf("fail to create directory %s: %v", dir, err)
		}
	}

	for path, file := range p.files {
		if file.isEmbeded {
			file.content = fmt.Sprintf(file.content, name)
		}
		if err := ioutil.WriteFile(filepath.Join(name, path), []byte(file.content), 0755); err != nil {
			return fmt.Errorf("fail to create file %s: %v", path, err)
		}
	}

	for _, cmd := range p.cmds {
		mc := cmd[0]
		shc := exec.Command(mc, cmd[1:]...)
		shc.Stdin = os.Stdin
		shc.Stdout = os.Stdout
		shc.Stderr = os.Stderr
		if err := shc.Run(); err != nil {
			return err
		}
	}

	return nil
}

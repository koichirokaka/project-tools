package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

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

	cd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(name); err != nil {
		return err
	}
	defer os.Chdir(cd)

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

func projectName(cmd *cobra.Command) (string, error) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", err
	} else if name == "" {
		return "", errors.New("name flag is required")
	}

	if strings.Contains(name, "~") {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		name = strings.Replace(name, "~", home, 1)
	}
	return name, nil
}

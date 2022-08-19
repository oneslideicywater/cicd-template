package main

import (
	"awesomeProject/generator"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/vifraa/gopom"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// parse from maven parent pom.xml, return modules name
func parseMavenModulesFromPom(pomPath string) (string, []string) {
	pom, err := gopom.Parse(pomPath)
	if err != nil {
		log.Fatal(err)
	}
	parent := pom.ArtifactID
	return parent, pom.Modules
}

// parse from npm package.json,return module name
func parsePackageJson(pomPath string) string {
	content, err := ioutil.ReadFile(filepath.Join(pomPath, "package.json"))
	if err != nil {
		log.Fatal(err.Error())
		return ""
	}
	var packageJson generator.PackageJson
	err = json.Unmarshal(content, &packageJson)
	if err != nil {
		log.Fatal(err.Error())
		return ""
	}
	return packageJson.Name
}

// JudgeProjectProfile detect project type in pwd, typically maven has pom.xml, npm has package.json
func JudgeProjectProfile(path string) (string, error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, file := range dir {
		if !file.IsDir() && file.Name() == "pom.xml" {
			return "maven", nil
		}
		if !file.IsDir() && file.Name() == "package.json" {
			return "npm", nil
		}
	}
	return "", errors.New("only maven and npm are supported")
}

func main() {

	var namespace = "default"
	var maintainers = []string{"test"}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err.Error())
	}

	// cli parameters parse
	app := &cli.App{
		Name:  "cicd-template",
		Usage: "ci/cd configuration generator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "namespace",
				Aliases:  []string{"n"},
				Usage:    "project namespace",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "maintainers",
				Aliases:  []string{"m"},
				Usage:    "project maintainers",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "project path,default is current directory",
			},
		},
		Action: func(c *cli.Context) error {
			// set namespace
			if c.String("namespace") != "" {
				namespace = c.String("namespace")
			}
			// get maintainers list
			if c.String("maintainers") != "" {
				maintainers = strings.Split(c.String("maintainers"), ",")
			}
			// get maintainers list
			if c.String("path") != "" {
				if _, err := os.Stat(path); err != nil {
					// if not exist, return error
					log.Fatal(err.Error())
				} else {
					path = c.String("path")
				}
			}
			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// 1. detect project type, maven or npm
	profile, err := JudgeProjectProfile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("project type: %s \n", profile)

	// 2. process maven project if profile is maven
	if profile == "maven" {

		// 1. force generate cicd.yaml
		// all modules has the sample format cicd.yaml, cicd.yaml is not maintained by developer
		parent, modules := parseMavenModulesFromPom(filepath.Join(path, "pom.xml"))
		// simple maven project without recursive module
		mtmp := generator.GenerateCICDSketch(parent, modules, namespace, maintainers, profile)
		result, err := yaml.Marshal(&mtmp)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = ioutil.WriteFile(filepath.Join(path, "cicd.yaml"), result, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}

		// 2. force generate Jenkinsfile
		jenkinsFile, err := generator.Asset("templates/maven/Jenkinsfile")
		if err != nil {
			fmt.Println(err.Error())
		}

		err = ioutil.WriteFile(filepath.Join(path, "Jenkinsfile"), jenkinsFile, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}

		// 3. multi-module maven project
		if len(modules) != 0 {
			// generate Dockerfile and .helm
			for _, module := range modules {
				// calc each submodule path
				modulePath := filepath.Join(path, module)
				// generate Dockerfile if not exist
				if _, err = os.Stat(filepath.Join(modulePath, "Dockerfile")); err != nil {
					dockerFile, err := generator.Asset("templates/maven/Dockerfile")
					if err != nil {
						log.Fatal(err.Error())
					}
					err = ioutil.WriteFile(filepath.Join(modulePath, "Dockerfile"), dockerFile, 0644)
					if err != nil {
						log.Fatal(err.Error())
					}
				}
				// generate helm chart if not exist
				if _, err = os.Stat(filepath.Join(modulePath, ".helm")); err != nil {
					err = generator.RestoreAssets(modulePath, "templates/maven/.helm")
					if err != nil {
						log.Fatal(err.Error())
					}
					// move .helm to upper
					err = os.Rename(filepath.Join(modulePath, "templates", "maven", ".helm"), filepath.Join(modulePath, ".helm"))
					if err != nil {
						log.Fatal(err.Error())
					}
					// delete template dir
					err = os.RemoveAll(filepath.Join(modulePath, "templates"))
					if err != nil {
						log.Fatal(err.Error())
					}
				}

			}
		}
		// 4. simple maven project
		if len(modules) == 0 {
			// generate Dockerfile and .helm
			modulePath := path
			// generate Dockerfile if not exist
			if _, err = os.Stat(filepath.Join(modulePath, "Dockerfile")); err != nil {
				dockerFile, err := generator.Asset("templates/maven/Dockerfile")
				if err != nil {
					log.Fatal(err.Error())
				}
				err = ioutil.WriteFile(filepath.Join(modulePath, "Dockerfile"), dockerFile, 0644)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
			// generate helm chart if not exist
			if _, err = os.Stat(filepath.Join(modulePath, ".helm")); err != nil {
				err = generator.RestoreAssets(modulePath, "templates/maven/.helm")
				if err != nil {
					log.Fatal(err.Error())
				}
				// move .helm to upper
				err = os.Rename(filepath.Join(modulePath, "templates", "maven", ".helm"), filepath.Join(modulePath, ".helm"))
				if err != nil {
					log.Fatal(err.Error())
				}
				// delete template dir
				err = os.RemoveAll(filepath.Join(modulePath, "templates"))
				if err != nil {
					log.Fatal(err.Error())
				}
			}

		}

	}

	// 3. process npm project if profile is npm
	if profile == "npm" {
		// generate cicd.yaml
		parent := parsePackageJson(path)
		mtmp := generator.GenerateCICDSketch(parent, nil, namespace, maintainers, profile)
		result, err := yaml.Marshal(&mtmp)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = ioutil.WriteFile(filepath.Join(path, "cicd.yaml"), result, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}

		// generate Jenkinsfile
		// 2. force generate Jenkinsfile
		jenkinsFile, err := generator.Asset("templates/npm/Jenkinsfile")
		if err != nil {
			fmt.Println(err.Error())
		}
		err = ioutil.WriteFile(filepath.Join(path, "Jenkinsfile"), jenkinsFile, 0644)
		if err != nil {
			log.Fatal(err.Error())
		}

		// 3. generate Dockerfile if not exist
		if _, err = os.Stat(filepath.Join(path, "Dockerfile")); err != nil {
			dockerFile, err := generator.Asset("templates/npm/Dockerfile")
			if err != nil {
				log.Fatal(err.Error())
			}
			err = ioutil.WriteFile(filepath.Join(path, "Dockerfile"), dockerFile, 0644)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		// 4. generate nginx.conf.template if not exist
		if _, err = os.Stat(filepath.Join(path, "nginx.conf.template")); err != nil {
			dockerFile, err := generator.Asset("templates/npm/nginx.conf.template")
			if err != nil {
				log.Fatal(err.Error())
			}
			err = ioutil.WriteFile(filepath.Join(path, "nginx.conf.template"), dockerFile, 0644)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		// generate helm chart if not exist
		if _, err = os.Stat(filepath.Join(path, ".helm")); err != nil {
			err = generator.RestoreAssets(path, "templates/npm/.helm")
			if err != nil {
				log.Fatal(err.Error())
			}
			// move .helm to upper
			err = os.Rename(filepath.Join(path, "templates", "npm", ".helm"), filepath.Join(path, ".helm"))
			if err != nil {
				log.Fatal(err.Error())
			}
			// delete template dir
			err = os.RemoveAll(filepath.Join(path, "templates"))
			if err != nil {
				log.Fatal(err.Error())
			}
		}

	}

	fmt.Println("generate success,happy ci/cd :) !")
}

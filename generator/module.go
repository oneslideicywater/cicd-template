package generator

type Module struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Maintainers []string          `yaml:"maintainers,omitempty"`
	Path        string            `yaml:"path"`
	BuildTool   string            `yaml:"build-tool"`
	Env         map[string]string `yaml:"env,omitempty"`
	Stages      []Stage           `yaml:"stages,omitempty"`
	Modules     []Module          `yaml:"modules,omitempty"`
	Resources   []Resource        `yaml:"resources,omitempty"`
}

type Stage struct {
	Name  string   `yaml:"name"`
	Shell []string `yaml:"shell"`
}

type Resource struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
	Dist   string `yaml:"dist"`
}

// GenerateCICDSketch generate cicd.yaml file content
func GenerateCICDSketch(parent string, modules []string, namespace string,
	maintainers []string, buildTool string) Module {

	var root Module

	// maven stages
	var stages = []Stage{
		{
			Name: "docker",
			Shell: []string{
				"docker build -t ${CURRENT_DOCKER_IMAGE_TAG} .",
				"docker push ${CURRENT_DOCKER_IMAGE_TAG}",
				"docker save ${CURRENT_DOCKER_IMAGE_TAG} > ${CURRENT_ARTIFACT_FILENAME}.tar",
				"mc cp ${CURRENT_ARTIFACT_FILENAME}.tar ${CURRENT_DOCKER_IMAGE_OBJECT_PATH}",
			},
		},
		{
			Name: "helm",
			Shell: []string{
				"(test ${GIT_BRANCH} = 'master' && kubectl config use-context prod) || true",
				"(test ${GIT_BRANCH} = 'develop' && kubectl config use-context dev) || true",
				"cd .helm; helm package .",
				"cd .helm; curl --data-binary \"@${CHART_NAME}-${CHART_VERSION}.tgz\" ${HELM_REGISTRY}/api/charts",
				"helm repo update geoway;helm upgrade -f ${COMMON_VALUE} ${APP_NAME} geoway/${CHART_NAME} --version ${CHART_VERSION} --install",
				"kubectl rollout restart deployment ${CHART_NAME} -n ${NAMESPACE}",
			},
		},
	}

	// simple maven project or npm
	if len(modules) == 0 {
		if buildTool == "maven" {
			root = Module{
				Name:        parent,
				Path:        ".",
				Namespace:   namespace,
				Maintainers: maintainers,
				Env: map[string]string{
					"COMMON_VALUE": "<common values link here>",
				},
				BuildTool: buildTool,
				Stages:    stages,
				Resources: []Resource{
					{
						Name:   "jar",
						Source: "target/*encrypted.jar",
						Dist:   "",
					},
					{
						Name:   "lic",
						Source: "target/*encrypted.lic",
						Dist:   "",
					},
					{
						Name:   "yml-application",
						Source: "target/classes/config/application.yml",
						Dist:   "config",
					},
					{
						Name:   "yml-bootstrap",
						Source: "target/classes/config/bootstrap.yml",
						Dist:   "config",
					},
					{
						Name:   "bash",
						Source: "deploy/bash/*.sh",
						Dist:   "",
					},
					{
						Name:   "contrib",
						Source: "deploy/contrib/*.service",
						Dist:   "",
					},
					{
						Name:   "wrapper",
						Source: "deploy/wrapper/*",
						Dist:   "",
					},
					{
						Name:   "cwd",
						Source: "deploy/cwd/*",
						Dist:   "",
					},
				},
			}
		}

		if buildTool == "npm" {

			// maven stages
			var npmStages = []Stage{
				{
					Name: "package",
					Shell: []string{
						"npm install -g yarn --registry=https://registry.npm.taobao.org",
						"yarn config set registry https://registry.npm.taobao.org -g",
						"yarn config set sass_binary_site http://cdn.npm.taobao.org/dist/node-sass -g",
						"yarn && yarn run build",
					},
				},
				{
					Name: "docker",
					Shell: []string{
						"docker build -t ${CURRENT_DOCKER_IMAGE_TAG} .",
						"docker push ${CURRENT_DOCKER_IMAGE_TAG}",
						"docker save ${CURRENT_DOCKER_IMAGE_TAG} > ${CURRENT_ARTIFACT_FILENAME}.tar",
						"mc cp ${CURRENT_ARTIFACT_FILENAME}.tar ${CURRENT_DOCKER_IMAGE_OBJECT_PATH}",
					},
				},
				{
					Name: "helm",
					Shell: []string{
						"(test ${GIT_BRANCH} = 'master' && kubectl config use-context prod) || true",
						"(test ${GIT_BRANCH} = 'develop' && kubectl config use-context dev) || true",
						"cd .helm; helm package .",
						"cd .helm; curl --data-binary \"@${CHART_NAME}-${CHART_VERSION}.tgz\" ${HELM_REGISTRY}/api/charts",
						"helm repo update geoway;helm upgrade ${APP_NAME} geoway/${CHART_NAME} --version ${CHART_VERSION} --install",
						"kubectl rollout restart deployment ${CHART_NAME} -n ${NAMESPACE}",
					},
				},
			}

			root = Module{
				Name:        parent,
				Path:        ".",
				Namespace:   namespace,
				Maintainers: maintainers,
				BuildTool:   buildTool,
				Stages:      npmStages,
				Resources: []Resource{
					{
						Name:   "dist",
						Source: "dist",
						Dist:   "",
					},
				},
			}
		}

	}

	// recursive modules maven project, multi-module npm is not supported
	if len(modules) != 0 && buildTool == "maven" {
		var submodules = make([]Module, 0)

		// init sub modules
		for _, name := range modules {
			submodules = append(submodules, Module{
				Name:      name,
				Path:      name,
				BuildTool: buildTool,
				Resources: []Resource{
					{
						Name:   "jar",
						Source: "target/*encrypted.jar",
						Dist:   "",
					},
					{
						Name:   "lic",
						Source: "target/*encrypted.lic",
						Dist:   "",
					},
					{
						Name:   "yml-application",
						Source: "target/classes/config/application.yml",
						Dist:   "config",
					},
					{
						Name:   "yml-bootstrap",
						Source: "target/classes/config/bootstrap.yml",
						Dist:   "config",
					},
					{
						Name:   "bash",
						Source: "deploy/bash/*.sh",
						Dist:   "",
					},
					{
						Name:   "contrib",
						Source: "deploy/contrib/*.service",
						Dist:   "",
					},
					{
						Name:   "wrapper",
						Source: "deploy/wrapper/*",
						Dist:   "",
					},
					{
						Name:   "cwd",
						Source: "deploy/cwd/*",
						Dist:   "",
					},
				},
				Stages: stages,
			})
		}
		root = Module{
			Name:        parent,
			Path:        ".",
			Namespace:   namespace,
			Maintainers: maintainers,
			Modules:     submodules,
			Env: map[string]string{
				"COMMON_VALUE": "<common values link here>",
			},
			BuildTool: buildTool,
		}
	}

	return root
}

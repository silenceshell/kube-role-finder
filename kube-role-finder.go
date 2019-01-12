/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func hasRule(clusterRole v1.ClusterRole, apiGroup, resource, verb string) bool {
	for _, rule := range clusterRole.Rules {
		for _, a := range rule.APIGroups {
			// check apiGroup first, if match then check resource
			if a == "*" || a == apiGroup {
				for _, r := range rule.Resources {
					if r == "*" || r == resource {
						// if not specify verb, skip verb checking
						if verb == "" {
							return true
						}
						// else check verb
						for _, v := range rule.Verbs {
							if v == verb {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func getTerminalSize() (high, length int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		high = 40
		length = 80
	} else {
		fmt.Sscanf(string(out), "%d %d", &high, &length)
	}
	return high, length
}

// In order to sort by a custom function in Go, we need a
// corresponding type. Here we've created a `byLength`
// type that is just an alias for the builtin `[]string`
// type.
type byLength []string

// We implement `sort.Interface` - `Len`, `Less`, and
// `Swap` - on our type so we can use the `sort` package's
// generic `Sort` function. `Len` and `Swap`
// will usually be similar across types and `Less` will
// hold the actual custom sorting logic. In our case we
// want to sort in order of increasing string length, so
// we use `len(s[i])` and `len(s[j])` here.
func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func printRoles(array []string) {
	var parser string
	_, width := getTerminalSize()
	width -= 6
	fmtStringLen := 0
	total := 0

	sort.Sort(byLength(array))
	base := 10
	i := 0

	for _, v := range array {
		var fmtString string
		f := ""

		//array have been sorted
		if len(v) < base {
			f = fmt.Sprintf("%%-%ds", base)
			fmtString = fmt.Sprintf(f, v)
		} else {
			var pad string
			if i % 2 == 1 {
				pad = fmt.Sprintf(fmt.Sprintf("%%%ds", base), " ")
			}
			base += base
			f += fmt.Sprintf(pad+"%%-%ds", base)
			i = 0
		}
		fmtString = fmt.Sprintf(f, v)
		i += 1

		fmtStringLen = len(fmtString)
		total += fmtStringLen
		if total > width {
			parser += "\r\n"
			total = fmtStringLen
			i = 0
		}
		parser += fmtString
	}
	fmt.Println(parser)
}

func main() {
	var apiGroup string
	var resource string
	var verb string
	var kubeconfig *string

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.StringVar(&apiGroup, "apiGroup", "", "(optional)apiGroup, default is \"\"")
	flag.StringVar(&resource, "resource", "", "resource")
	flag.StringVar(&verb, "verb", "", "(optional)verb, default match all")

	flag.Parse()

	if resource == "" {
		flag.Usage()
		return
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clusterRolesNames := make([]string, 0)
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(metav1.ListOptions{})
	for _, clusterRole := range clusterRoles.Items {
		if hasRule(clusterRole, apiGroup, resource, verb) {
			clusterRolesNames = append(clusterRolesNames, clusterRole.Name)
		}
	}

	if len(clusterRolesNames) != 0 {
		fmt.Printf("Those clusterRole has resource %s:\n", resource)
		printRoles(clusterRolesNames)
	} else {
		fmt.Printf("No clusterRole has resource %s.", resource)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

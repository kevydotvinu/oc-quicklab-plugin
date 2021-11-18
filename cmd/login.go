/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate the quicklab shared cluster using oc binary.",
	Long:  `This command fetches the credentials from the portal and authenticates the quicklab shared cluster using oc binary. The username will be kubeadmin.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkArg()
		a := clusterNameAndLinks()
		fmt.Println(a)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringP("cluster", "c", "", "Pass quicklab shared cluster name")
	loginCmd.MarkFlagRequired("cluster")
}

func checkArg() {
	fmt.Println(os.Args[3])
}

func clusterNameAndLinks() map[string]string {
	name, links, _ := getClustersList()
	rowLength := len(name)
	nameLinkMap := make(map[string]string)
	for i := 1; i < rowLength; i++ {
		j := name[i][0]
		k := links[i][0]
		nameLinkMap[j] = k
	}
	return nameLinkMap
}

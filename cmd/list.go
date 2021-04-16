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
	"log"
	"os"

	"github.com/sebsebmc/ctf-watcher/lib"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays a list of challenges",
	Long:  `Lists all of the challenges in the CTF`,
	Run:   pull,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func pull(cmd *cobra.Command, args []string) {
	creds, err := lib.GetCreds()
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("No credentials found, please 'init' first")
		} else {
			log.Fatalf("Unable to list: %v", err)
		}
	}

	inst := lib.MakeCtfdInstance(creds.Url, creds.Username, creds.Password)
	// fmt.Println("Logging in")
	err = inst.LoginToSite()
	if err != nil {
		log.Fatalf("Unable to list: %v", err)
	}
	// fmt.Println("Fetching challenges")
	err = inst.GetLatestChallenges()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Printing...")
	inst.PrintChallenges()
}

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
	"log"
	"os"
	"path"

	"github.com/sebsebmc/ctf-watcher/lib"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init ctf-url username password target-dir [test]",
	Short: "Setup a directory for syncing from the CTF",
	Long: `cobra init to write settings for syncing from the CTF.
The target-directory can be '.' and the CTF name will be inferred from the current folder`,
	Run:  writeCtfFile,
	Args: cobra.ExactArgs(4),
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().BoolP("test", "t", false, "Attempt to login to confirm credentials")
	initCmd.Flags().BoolP("force", "f", false, "Replace an existing .ctf credentials file")
}

func writeCtfFile(cmd *cobra.Command, args []string) {
	os.MkdirAll(args[3], 0750)
	shouldTest, _ := cmd.Flags().GetBool("test")
	shouldReplace, _ := cmd.Flags().GetBool("force")

	if shouldTest {
		inst := lib.MakeCtfdInstance(args[0], args[1], args[2])
		err := inst.LoginToSite()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Credentials valid, proceeding")
	}

	settingsFile := path.Join(args[3], ".ctf")
	flags := os.O_RDWR | os.O_CREATE
	if !shouldReplace {
		flags |= os.O_EXCL
	} else {
		flags |= os.O_TRUNC
	}
	fh, err := os.OpenFile(settingsFile, flags, 0666)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("Unable to create file %s: %v", settingsFile, err)
		} else if os.IsExist(err) {
			log.Fatalln(".ctf file already found, pass -f or --force to replace the old credentials")
		}
	}
	defer fh.Close()

	settings := fmt.Sprintf("%s\n%s\n%s\n", args[0], args[1], args[2])
	_, err = fh.Write([]byte(settings))
	if err != nil {
		log.Fatalf("Failed to write settings file: %v", err)
	}
}

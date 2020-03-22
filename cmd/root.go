/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/techniumlabs/cinit/pkg/proc"
	"github.com/techniumlabs/cinit/pkg/secrets"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cinit",
	Short: "Init process for containers",
	Long: `Init process for containers which does the following
1. Proper Signal Forwarding
2. Orphaned Zombies Reaping
3. Fetch Secrets and expose it as Environment Variables`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Println("No Main Command Provided")
			os.Exit(1)
		}
		// Get and expose any secrets
		client := secrets.NewSecretsClient()
		envs := client.GetParsedEnvs()

		// Replace Template Files

		// Execute any pre commands and post commands on exit

		// Execute the provided command
		// Routine to reap zombies (it's the job of init)
		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		wg.Add(1)
		go proc.RemoveZombies(ctx, &wg)

		var mainRC int
		var argsSlice []string
		command := args[0]
		if len(args) > 1 {
			argsSlice = args[1:]
		}
		err := proc.Run(command, argsSlice, envs)
		if err != nil {
			log.Println("Main command failed with error", err.Error())
			mainRC = 1
		} else {
			log.Printf("Main command exited")
		}

		proc.CleanQuit(cancel, &wg, mainRC)
	},
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cinit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/")
		viper.SetConfigName(".cinit")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("%s", err.Error())
	}
}

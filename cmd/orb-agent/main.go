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
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ns1labs/orb/agent"
)

var (
	cfgFiles      []string
	defaultConfig = "/etc/orb/agent.yml"

	rootCmd = &cobra.Command{
		Use:   "orb-agent",
		Short: "orb-agent connects to orb control plane",
		Long:  "orb-agent connects to orb control plane",
		Run:   Run,
	}
)

func Run(cmd *cobra.Command, args []string) {
	var config agent.Config
	viper.Unmarshal(&config)
	fmt.Printf("%+v\n", config)
	s, err := agent.New(config)
	cobra.CheckErr(err)
	cobra.CheckErr(s.Start())
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringSliceVarP(&cfgFiles, "config", "c", []string{}, "Path to config files (can be specified multiple times)")
}

func mergeOrError(path string) {
	v := viper.New()
	v.SetConfigFile(path)
	cobra.CheckErr(v.ReadInConfig())

	var fZero float64

	// check that version of config files are all matched up
	if versionNumber1 := viper.GetFloat64("version"); versionNumber1 != fZero {
		versionNumber2 := v.GetFloat64("version")
		if versionNumber2 == fZero {
			cobra.CheckErr("Failed to parse config vesrion in: " + path)
		}
		if versionNumber2 != versionNumber1 {
			cobra.CheckErr("Config file version mismatch in: " + path)
		}
	}
	fmt.Fprintln(os.Stderr, "Using config file:", v.ConfigFileUsed())

	cobra.CheckErr(viper.MergeConfigMap(v.AllSettings()))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("ORB_AGENT")
	viper.AutomaticEnv() // read in environment variables that match

	if len(cfgFiles) == 0 {
		mergeOrError(defaultConfig)
	} else {
		for _, conf := range cfgFiles {
			mergeOrError(conf)
		}
	}
}

func main() {
	rootCmd.Execute()
}

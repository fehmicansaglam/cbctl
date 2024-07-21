package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fehmicansaglam/cbctl/cmd/config"
	"github.com/fehmicansaglam/cbctl/cmd/get"
	"github.com/fehmicansaglam/cbctl/cmd/query"
	"github.com/fehmicansaglam/cbctl/constants"
	"github.com/fehmicansaglam/cbctl/shared"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cbctl",
	Short: "cbctl is CLI for Couchbase",
	Long:  `cbctl is a read-only CLI for Couchbase that allows users to manage and monitor their Couchbase clusters.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initialize)

	initProtocolFlag()
	initHostFlag()
	initPortFlag()
	initUsernameFlag()
	initPasswordFlag()

	rootCmd.PersistentFlags().StringVar(&shared.Context, "context", "", "Override context")
	rootCmd.PersistentFlags().BoolVar(&shared.Debug, "debug", false, "Enable debug mode")

	rootCmd.AddCommand(get.Cmd())
	rootCmd.AddCommand(query.Cmd())
}

func initialize() {
	if shared.CouchbaseHost == "" {
		conf := config.ParseConfigFile()
		readContextFromConfig(conf)
	}
}

func readContextFromConfig(conf config.Config) {
	if len(conf.Contexts) == 0 {
		fmt.Println("Error: No contexts defined in the configuration.")
		os.Exit(1)
	}

	var context string

	if shared.Context != "" {
		context = shared.Context
	} else if conf.CurrentContext != "" {
		context = conf.CurrentContext
	} else {
		context = conf.Contexts[0].Name
	}

	clusterFound := false
	for _, cluster := range conf.Contexts {
		if cluster.Name == context {
			shared.CouchbaseProtocol = cluster.Protocol
			if shared.CouchbaseProtocol == "" {
				shared.CouchbaseProtocol = constants.DefaultCouchbaseProtocol
			}
			shared.CouchbasePort = cluster.Port
			if shared.CouchbasePort == 0 {
				shared.CouchbasePort = constants.DefaultCouchbasePort
			}
			shared.CouchbaseUsername = cluster.Username
			shared.CouchbasePassword = cluster.Password
			shared.CouchbaseHost = cluster.Host
			if shared.CouchbaseHost == "" {
				fmt.Println("Error: 'host' field is not specified in the configuration for the current cluster.")
				os.Exit(1)
			}
			clusterFound = true
			break
		}
	}

	if !clusterFound {
		fmt.Printf("Error: No cluster found with the name '%s' in the configuration.\n", conf.CurrentContext)
		os.Exit(1)
	}
}

func initProtocolFlag() {
	defaultProtocol := constants.DefaultCouchbaseProtocol
	defaultProtocolEnv := os.Getenv(constants.CouchbaseProtocolEnvVar)
	if defaultProtocolEnv != "" {
		defaultProtocol = defaultProtocolEnv
	}
	rootCmd.PersistentFlags().StringVar(&shared.CouchbaseProtocol, "protocol", defaultProtocol, "Couchbase protocol")
}

func initHostFlag() {
	defaultHost := os.Getenv(constants.CouchbaseHostEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.CouchbaseHost, "host", defaultHost, "Couchbase host")
}

func initPortFlag() {
	defaultPort := constants.DefaultCouchbasePort
	defaultPortStr := os.Getenv(constants.CouchbasePortEnvVar)
	if defaultPortStr != "" {
		parsedPort, err := strconv.Atoi(defaultPortStr)
		if err != nil || parsedPort <= 0 {
			fmt.Printf("Invalid value for %s environment variable: %s\n", constants.CouchbasePortEnvVar, defaultPortStr)
			os.Exit(1)
		}
		defaultPort = parsedPort
	}
	rootCmd.PersistentFlags().IntVar(&shared.CouchbasePort, "port", defaultPort, "Couchbase port")
}

func initUsernameFlag() {
	defaultUsername := os.Getenv(constants.CouchbaseUsernameEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.CouchbaseUsername, "username", defaultUsername, "Couchbase username")
}

func initPasswordFlag() {
	defaultPassword := os.Getenv(constants.CouchbasePasswordEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.CouchbasePassword, "password", defaultPassword, "Couchbase password")
}

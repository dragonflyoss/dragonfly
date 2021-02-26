package cmd

import (
	logger "d7y.io/dragonfly/v2/pkg/dflog"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/server"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	// SupernodeEnvPrefix is the default environment prefix for Viper.
	// Both BindEnv and AutomaticEnv will use this prefix.
	SchedulerEnvPrefix = "scheduler"
)

var (
	cdnList string
	schedulerViper = viper.GetViper()
)

// supernodeDescription is used to describe supernode command in details.
var supernodeDescription = `scheduler is a long-running process with two primary responsibilities:
It's the tracker and scheduler in the P2P network that choose appropriate downloading net-path for each peer.`

var rootCmd = &cobra.Command{
	Use:               "scheduler",
	Short:             "the central control server of Dragonfly used for scheduling",
	Long:              supernodeDescription,
	Args:              cobra.NoArgs,
	DisableAutoGenTag: true, // disable displaying auto generation tag in cli docs
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := logger.InitScheduler()
		if err != nil {
			return errors.Wrap(err, "init scheduler logger")
		}

		// load config file.
		if err = readConfigFile(schedulerViper, cmd); err != nil {
			return errors.Wrap(err, "read config file")
		}

		// get config from viper.
		cfg, err := getConfigFromViper(schedulerViper)
		if err != nil {
			return errors.Wrap(err, "get config from viper")
		}

		logger.Debugf("get scheduler config: %+v", cfg)
		logger.Infof("start to run scheduler")

		svr := server.NewServer()
		return svr.Start()
	},
}

func init() {
	setupFlags(rootCmd)
}

// setupFlags setups flags for command line.
func setupFlags(cmd *cobra.Command) {
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	// flagSet := cmd.PersistentFlags()

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	flagSet := cmd.Flags()

	defaultBaseProperties := config.GetConfig()

	flagSet.String("config", config.DefaultConfigFilePath,
		"the path of scheduler's configuration file")

	flagSet.Int("port", defaultBaseProperties.Server.Port,
		"port is the port that scheduler server listens on")

	flagSet.Int("worker-num", defaultBaseProperties.Worker.WorkerNum,
		"worker-num is used for scheduler and do not change it")

	flagSet.Int("worker-job-pool-size", defaultBaseProperties.Worker.WorkerJobPoolSize,
		"worker-job-pool-size is used for scheduler and do not change it")

	flagSet.Int("sender-num", defaultBaseProperties.Worker.SenderNum,
		"sender-num is used for scheduler and do not change it")

	flagSet.Int("sender-job-pool-size", defaultBaseProperties.Worker.WorkerJobPoolSize,
		"sender-job-pool-size is used for scheduler and do not change it")

	flagSet.Var(config.NewCdnValue(&defaultBaseProperties.CDN), "cdn-list",
		"cdn list with format of [CdnName1]:[ip1]:[rpcPort1]:[downloadPort1]|[CdnName2]:[ip2]:[rpcPort2]:[downloadPort2]")

	exitOnError(bindRootFlags(schedulerViper), "bind root command flags")
}

// bindRootFlags binds flags on rootCmd to the given viper instance.
func bindRootFlags(v *viper.Viper) error {
	flags := []struct {
		key  string
		flag string
	}{
		{
			key:  "config",
			flag: "config",
		},
		{
			key:  "server.port",
			flag: "port",
		},
		{
			key:  "worker.worker-num",
			flag: "worker-num",
		},
		{
			key:  "worker.worker-job-pool-size",
			flag: "worker-job-pool-size",
		},
		{
			key:  "worker.sender-num",
			flag: "sender-num",
		},
		{
			key:  "worker.sender-job-pool-size",
			flag: "sender-job-pool-size",
		},
	}

	for _, f := range flags {
		if err := v.BindPFlag(f.key, rootCmd.Flag(f.flag)); err != nil {
			return err
		}
	}

	v.SetEnvPrefix(SchedulerEnvPrefix)
	v.AutomaticEnv()

	return nil
}

// readConfigFile reads config file into the given viper instance. If we're
// reading the default configuration file and the file does not exist, nil will
// be returned.
func readConfigFile(v *viper.Viper, cmd *cobra.Command) error {
	v.SetConfigFile(v.GetString("config"))
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		// when the default config file is not found, ignore the error
		if os.IsNotExist(err) && !cmd.Flag("config").Changed {
			return nil
		}
		return err
	}

	return nil
}

// getDefaultConfig returns the default configuration of scheduler
func getDefaultConfig() (interface{}, error) {
	return getConfigFromViper(viper.GetViper())
}

// getConfigFromViper returns supernode config from the given viper instance
func getConfigFromViper(v *viper.Viper) (*config.Config, error) {
	cfg := config.GetConfig()

	if err := v.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
		dc.Squash = true
		dc.DecodeHook = decodeWithYAML(
			reflect.TypeOf(time.Second),
		)
	}); err != nil {
		return nil, errors.Wrap(err, "unmarshal yaml")
	}

	return cfg, nil
}

// decodeWithYAML returns a mapstructure.DecodeHookFunc to decode the given
// types by unmarshalling from yaml text.
func decodeWithYAML(types ...reflect.Type) mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data interface{}) (interface{}, error) {
		for _, typ := range types {
			if t == typ {
				b, _ := yaml.Marshal(data)
				v := reflect.New(t)
				return v.Interface(), yaml.Unmarshal(b, v.Interface())
			}
		}
		return data, nil
	}
}

// Execute will process supernode.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
}

func exitOnError(err error, msg string) {
	if err != nil {
		logger.Errorf("%s: %v", msg, err)
	}
}

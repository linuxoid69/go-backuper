package cmd

import (
	"flag"
	"fmt"
	"github.com/linuxoid69/go-backuper/config"
	"github.com/linuxoid69/go-backuper/cron"
	"github.com/linuxoid69/go-backuper/encryption"
	"log"
	"os"
	"time"
)

// CmdFlags - all flags of command line
var (
	flagFile     string
	flagPassword string
	flagConfig   string
	flagDaemon   bool
	// help         bool
)

func cmdInit() {
	// flag.BoolVar(&help, "h", false, "Help")
	flag.StringVar(&flagFile, "f", "", "Name of a file is for decrypt")
	flag.BoolVar(&flagDaemon, "d", false, "Run as daemon")
	flag.StringVar(&flagPassword, "p", "", "Password is for decrypt")
	flag.StringVar(&flagConfig, "c", "./config.yaml", "config")
	flag.Parse()
}

const (
	exitError   = 1
	exitSuccess = 0
)

// Run - main function for use flags
func Run() {
	cmdInit()
	checkFlags()
	cmdDecryptFile()
	RunDaemon()
}

func cmdExit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

func cmdDecryptFile() {
	if flagFile != "" && flagPassword == "" {
		cmdExit("You must set a password!", exitError)
	} else if flagPassword != "" && flagFile == "" {
		cmdExit("You must set name of a file!", exitError)
	} else if flagPassword == "" && flagFile == "" {
		return
	} else {
		err := encryption.DencryptFile(flagPassword, flagFile)
		if err != nil {
			cmdExit("Password is invalid!", exitError)
		}
		cmdExit(fmt.Sprintf("File %s was seccesfully decrypted", flagFile), exitSuccess)
	}
	log.Println(flagFile)
}

func checkFlags() {
	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(0)
	}
}

// GetConfig - function get config
func GetConfig() (cfg *config.Config, err error) {
	log.Println("Start load config")
	_, envConfig := os.LookupEnv("CONFIG")

	if envConfig == true {
		log.Println("Get config from env CONFIG")
		cfg, err = config.LoadConfig(os.Getenv("CONFIG"))
	} else if envConfig == false && flagConfig != "" {
		log.Println("Get config from file")
		cfg, err = config.LoadConfig(flagConfig)
	}

	if err != nil {
		log.Fatalf("Can't load config %v: ", err)
	}

	log.Println("Config was successfully loaded")
	return cfg, nil
}

// RunDaemon - function run application in mode daemon
func RunDaemon() {
	if flagDaemon == true {
		log.Println("Start daemon")
		cfg, _ := GetConfig()
		cron.TasksCron(cfg)
		for {
			time.Sleep(time.Millisecond * 10)
		}
	}
}

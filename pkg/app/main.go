package app

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/techniumlabs/cinit/pkg/config"
	"github.com/techniumlabs/cinit/pkg/proc"
	"github.com/techniumlabs/cinit/pkg/secrets"
	"github.com/techniumlabs/cinit/pkg/templates"
)

type App struct {
	Config *config.Config
}

type FileTemplate struct {
	Source string
	Dest   string
}

func NewApp(cfgFile string) (*App, error) {
	c, err := config.Load(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config %x", err)
	}
	return &App{Config: c}, nil
}

func (a *App) RunInit(args []string) {
	var err error

	if len(args) == 0 {
		log.Fatal("No Main Command Provided")
	}
	// Get and expose any secrets
	client := secrets.NewSecretsClient(a.Config)
	envs := client.GetParsedEnvs()

	// Replace Template Files
	tclient, _ := templates.NewTemplateClient("default")

	for _, elem := range a.Config.Templates {
		err = tclient.Provider.ResolveTemplates(elem.Source, elem.Dest, envs)
		if err != nil {
			log.Error("Template could not be Resolved")
		}
	}
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
	err = proc.Run(command, argsSlice, envs)
	if err != nil {
		log.Println("Main command failed with error", err.Error())
		mainRC = 1
	} else {
		log.Printf("Main command exited")
	}

	proc.CleanQuit(cancel, &wg, mainRC)

}

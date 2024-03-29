package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaulYakow/metrics-track/cmd/agent/config"
	"github.com/PaulYakow/metrics-track/internal/app/client"
)

/*
  Build version: <buildVersion> (или "N/A" при отсутствии значения)
  Build date: <buildDate> (или "N/A" при отсутствии значения)
  Build commit: <buildCommit> (или "N/A" при отсутствии значения)
*/

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	cfg, err := config.NewClientConfig()
	if err != nil {
		log.Printf("agent - create config: %v", err)
		return
	}

	printInfo()
	c := client.New(cfg)
	c.Run()
}

func printInfo() {
	var sb strings.Builder

	sb.WriteString("Build version: ")
	sb.WriteString(buildVersion)
	sb.WriteString("\n")
	sb.WriteString("Build date: ")
	sb.WriteString(buildDate)
	sb.WriteString("\n")
	sb.WriteString("Build commit: ")
	sb.WriteString(buildCommit)
	sb.WriteString("\n")
	fmt.Println(sb.String())
}

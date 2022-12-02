package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/app/client"
)

/*
  Build version: <buildVersion> (или "N/A" при отсутствии значения)
  Build date: <buildDate> (или "N/A" при отсутствии значения)
  Build commit: <buildCommit> (или "N/A" при отсутствии значения)
*/

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	cfg, err := config.NewClientConfig()
	if err != nil {
		log.Printf("agent - create config: %v", err)
		return
	}

	bi := buildInfo{}
	bi.printInfo()

	client.Run(cfg)
}

type buildInfo struct {
	buildVersion string
	buildDate    string
	buildCommit  string
}

func (bi *buildInfo) printInfo() {
	var sb strings.Builder

	bi.buildVersion = "N/A"
	bi.buildDate = "N/A"
	bi.buildCommit = "N/A"

	if buildVersion != "" {
		bi.buildVersion = buildVersion
	}

	if buildDate != "" {
		bi.buildDate = buildDate
	}

	if buildCommit != "" {
		bi.buildCommit = buildCommit
	}

	sb.WriteString("Build version: ")
	sb.WriteString(bi.buildVersion)
	sb.WriteString("\n")
	sb.WriteString("Build date: ")
	sb.WriteString(bi.buildDate)
	sb.WriteString("\n")
	sb.WriteString("Build commit: ")
	sb.WriteString(bi.buildCommit)
	sb.WriteString("\n")
	fmt.Println(sb.String())
}

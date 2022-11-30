package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaulYakow/metrics-track/config"
	"github.com/PaulYakow/metrics-track/internal/app/server"
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
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("server - create config: %v\n", err)
	}

	bi := buildInfo{}
	bi.printInfo()

	server.Run(cfg)
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

	switch {
	case buildVersion != "":
		bi.buildVersion = buildVersion
	case buildDate != "":
		bi.buildDate = buildDate
	case buildCommit != "":
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

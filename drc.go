package main

import (
	"fmt"
	"github.com/pbergman/docker-registry-cli/api"
	"github.com/pbergman/docker-registry-cli/config"
	"github.com/pbergman/docker-registry-cli/helpers"
	"github.com/pbergman/docker-registry-cli/http"
	"github.com/pbergman/docker-registry-cli/logger"
	"io"
	"os"
	"os/exec"
	"time"
)

func init() {
	api.ApiCheck()
}

func main() {
	switch config.Config.Command {
	case config.DELETE:
		api.Delete(
			*config.Config.Input["delete.repository"].(*string),
			*config.Config.Input["delete.tag"].(*string),
			*config.Config.Input["delete.force"].(*bool),
		)
		fmt.Printf("Image %s:%s removed\n", *config.Config.Input["delete.repository"].(*string), *config.Config.Input["delete.tag"].(*string))
	case config.TOKEN:
		token, err := config.Config.TokenManager.GetToken(&http.AuthChallenge{
			Service: *config.Config.Input["token.service"].(*string),
			Realm:   *config.Config.Input["token.realm"].(*string),
			Scope:   *config.Config.Input["token.scope"].(*string),
		})
		logger.Logger.CheckError(err)
		fmt.Println("")
		fmt.Println("Expires: " + token.ExpireTime().Format("2006-01-02 15:04:05 UTC"))
		fmt.Println(token.Token)
		fmt.Println("")
	case config.LIST:
		table := helpers.NewTable("REPOSITORY", "TAG", "DATE", "AUTHOR", "SIZE(MB)")
		list := api.GetList()
		list.Sort()
		for _, info := range *list {
			for _, tag := range info.Tags {
				manifest := api.GetManifest(info.Name, tag, true)
				author := ""
				for _, history := range manifest.History {
					if value, exist := history.Unpack()["author"]; exist {
						author = value.(string)
						break
					}
				}

				time := time.Time{}
				time.UnmarshalText([]byte(manifest.History[len(manifest.History)-1].Unpack()["created"].(string)))
				table.AddRow(
					info.Name,
					tag,
					time.Format("2006-01-02 15:04:05.000000"),
					author,
					api.GetSize(info.Name, tag)/(1024*1024),
				)
			}
		}
		table.Print()
	case config.SIZE:
		size := api.GetSize(
			*config.Config.Input["size.repository"].(*string),
			*config.Config.Input["size.tag"].(*string),
		)
		fmt.Printf("\n%dMB\n", size/(1024*1024))
	case config.HISTORY:
		manifest := api.GetManifest(
			*config.Config.Input["history.repository"].(*string),
			*config.Config.Input["history.tag"].(*string),
			true,
		)
		if config.Config.Verbose { // Print json set to less
			path, err := exec.LookPath("less")
			logger.Logger.CheckError(err)
			cmd := exec.Command(path, "-egR")
			stdin, err := cmd.StdinPipe()
			logger.Logger.CheckError(err)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			for i := len(manifest.History) - 1; i >= 0; i-- {
				io.WriteString(stdin, string(manifest.History[i].Print().Bytes()))
				io.WriteString(stdin, "\n")
			}
			stdin.Close()
			err = cmd.Run()
			logger.Logger.CheckError(err)
		} else { // Print only date and command
			for i := len(manifest.History) - 1; i >= 0; i-- {
				data := manifest.History[i].Unpack()
				cmd := data["container_config"].(map[string]interface{})["Cmd"]
				if cmd != nil {
					line := ""
					for _, part := range cmd.([]interface{}) {
						line += part.(string) + " "
					}
					time := time.Time{}
					time.UnmarshalText([]byte(data["created"].(string)))
					fmt.Printf("[%s] %s\n", time.Format("2006-01-02 15:04:05.000000"), line)
				}
			}
		}

	case config.REPOSITORIES:
		repositories := api.GetRepositories()
		table := helpers.NewTable("REPOSITORIES")
		for _, name := range repositories.Images {
			table.AddRow(name)
		}
		table.Print()
	case config.TAGS:
		tags := api.GetTags(*config.Config.Input["tag.repository"].(*string))
		if tags != nil {
			table := helpers.NewTable("TAGS(" + *config.Config.Input["tag.repository"].(*string) + ")")
			for _, name := range tags.Tags {
				table.AddRow(name)
			}
			table.Print()
		}
	}
}

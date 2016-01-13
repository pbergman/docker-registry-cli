package main

import (
	"fmt"
	"github.com/pbergman/docker-registery-cli/api"
	"github.com/pbergman/docker-registery-cli/config"
	"github.com/pbergman/docker-registery-cli/http"
	"github.com/pbergman/docker-registery-cli/logger"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func init() {
	api.ApiCheck()
}

func main() {
	switch config.Config.Command {
	case config.DELETE:
		api.Delete(*config.Config.Input["delete.repository"], *config.Config.Input["delete.tag"])
		fmt.Printf("Image %s:%s removed\n", *config.Config.Input["delete.repository"], *config.Config.Input["delete.tag"])
	case config.TOKEN:
		token, err := config.Config.TokenManager.GetToken(&http.AuthChallenge{
			Service: *config.Config.Input["token.service"],
			Realm:   *config.Config.Input["token.realm"],
			Scope:   *config.Config.Input["token.scope"],
		})
		logger.Logger.CheckError(err)
		fmt.Println("")
		fmt.Println("Expires: " + token.ExpireTime().Format("2006-01-02 15:04:05 UTC"))
		fmt.Println(token.Token)
		fmt.Println("")
	case config.LIST:
		type tag struct {
			name string
			size int
			time time.Time
		}
		list := make(map[string][]*tag, 0)
		sizes := make([]int, 2)
		sizes[0] = len("REPOSITORY")
		sizes[1] = len("TAG")
		for _, reposName := range api.GetRepositories().Images {
			if len(reposName) > sizes[0] {
				sizes[0] = len(reposName)
			}
			tags := api.GetTags(reposName)
			for _, tagName := range tags.Tags {
				if len(tagName) > sizes[1] {
					sizes[1] = len(reposName)
				}
				manifest := api.GetManifest(reposName, tagName)
				time := time.Time{}
				time.UnmarshalText([]byte(manifest.History[len(manifest.History)-1].Unpack()["created"].(string)))
				list[reposName] = append(list[reposName], &tag{
					name: tagName,
					size: api.GetSize(reposName, tagName),
					time: time,
				})
			}
		}

		pad := func(string, padding string, width int) string {
			if left := width - len(string); left > 0 {
				return string + strings.Repeat(padding, left)
			}
			return string
		}

		format := "2006-01-02 15:04:05.000000"

		fmt.Printf(
			"%s\t%s\t%s\tSIZE\n",
			pad("REPOSITORY", " ", sizes[0]),
			pad("TAG", " ", sizes[1]),
			pad("DATE", " ", len(format)),
		)

		for repository, tags := range list {
			for _, tag := range tags {
				fmt.Printf(
					"%s\t%s\t%s\t%dMB\n",
					pad(repository, " ", sizes[0]),
					pad(tag.name, " ", sizes[1]),
					tag.time.Format(format),
					tag.size/(1024*1024),
				)
			}
		}

	case config.SIZE:
		size := api.GetSize(*config.Config.Input["size.repository"], *config.Config.Input["size.tag"])
		fmt.Printf("\n%dMB\n", size/(1024*1024))
	case config.HISTORY:
		manifest := api.GetManifest(*config.Config.Input["history.repository"], *config.Config.Input["history.tag"])
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
		fmt.Println("")
		fmt.Println("REPOSITORIES")
		for _, name := range repositories.Images {
			fmt.Println(name)
		}
		fmt.Println("")
	case config.TAGS:
		tags := api.GetTags(*config.Config.Input["tag.repository"])
		if tags != nil {
			fmt.Println("")
			fmt.Printf("TAGS <%s>\n", tags.Name)
			for _, name := range tags.Tags {
				fmt.Println(name)
			}
			fmt.Println("")
		}
	}
}

package config

import (
	"github.com/pbergman/docker-registry-cli/helpers"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"regexp"
)

const (
	LIST = 1 << iota
	REPOSITORIES
	TAGS
	TOKEN
	HISTORY
	SIZE
	DELETE
)

func (config *config) ParseInput() {

	// Add to args to bypass defaults field
	if config.RegistryHost != "" {

		hasHost := false

		for _, arg := range os.Args {
			if reg, _ := regexp.Compile(`^-?-(h|registry-host)(=|\s)`); reg.MatchString(arg) {
				hasHost = true
				break
			}
		}
		if !hasHost && config.RegistryHost != "" {
			os.Args = append(os.Args, "--registry-host="+config.RegistryHost)
		}
	}

	// Global Config Command
	app := kingpin.New("drc", "A cli that can communicate with a docker private register.")
	app.Flag("verbose", "Verbose mode.").Short('v').BoolVar(&config.Verbose)
	app.Flag("username", "Username for login.").Short('u').StringVar(&config.User.Username)
	app.Flag("password", "Password for login.").Short('p').StringVar(&config.User.Password)
	app.Flag("registry-host", "Host of registry.").Short('h').Default("https://index.docker.io").StringVar(&config.RegistryHost)

	// List Command
	list := app.Command("list", "Retrieve a list of repositories and tags available in the registry.")

	// Catalog Command
	repositories := app.Command("repositories", "Retrieve a list of repositories available in the registry.")

	// Tags Command
	tags := app.Command("tags", "List all of the tags under the given repository.")
	config.Input["tag.repository"] = tags.Arg("repository", "Repository to fetch tags from.").Required().String()

	// History Command
	history := app.Command("history", "Get history infomation from given repository.")
	config.Input["history.repository"] = history.Arg("repository", "Repository to fetch tags from.").Required().String()
	config.Input["history.tag"] = history.Arg("tag", "Tag of repository").Default("latest").String()

	// Delete Command
	delete := app.Command("delete", "Delete tagged repository")
	config.Input["delete.repository"] = delete.Arg("repository", "Repository to fetch tags from.").Required().String()
	config.Input["delete.tag"] = delete.Arg("tag", "Tag of repository").Default("latest").String()
	config.Input["delete.dry"] = delete.Flag("dry-run", "Check which layers are getting queues for remooval").Bool()

	// Size Command
	size := app.Command("size", "Get size infomation from given repository.")
	config.Input["size.repository"] = size.Arg("repository", "Repository to fetch tags from.").Required().String()
	config.Input["size.tag"] = size.Arg("tag", "Tag of repository").Default("latest").String()

	// Token Command
	token := app.Command("token", "Create api token for docker register server")
	config.Input["token.service"] = token.Arg("service", "The name of the service which hosts the resource.").Required().String()
	config.Input["token.realm"] = token.Arg("realm", "Token authentication server.").Required().String()
	config.Input["token.scope"] = token.Arg("scope", "space-delimited list of case-sensitive scope values indicating the required scope of the access token for accessing the requested resource, exmaple: repository:busybox:push").Required().String()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case repositories.FullCommand():
		config.Command = REPOSITORIES
	case list.FullCommand():
		config.Command = LIST
	case tags.FullCommand():
		config.Command = TAGS
	case token.FullCommand():
		config.Command = TOKEN
	case history.FullCommand():
		config.Command = HISTORY
	case size.FullCommand():
		config.Command = SIZE
	case delete.FullCommand():
		config.Command = DELETE
	}
}

func (config *config) CheckUser() {
	if config.User.Username == "" {
		config.User.Username = helpers.Ask("username: ")
	}

	if config.User.Password == "" {
		config.User.Password = helpers.Password("password: ")
	}
}

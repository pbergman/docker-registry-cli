##DOCKER PRIVATE REGISTRY CLI

This is a small tool i made to communicate with a docker private registry. It has some commands that the normal docker
doesn`t have like listing, deleting of image on a private server. This version is only tested with a private server 
configured with token authentication and docker registry 2.1 but it should work with server without authentication.

####Commands:

#####list
    Retrieve a list of repositories and tags available in the registry.

#####repositories
    Retrieve a list of repositories available in the registry.

#####tags \<repository\>
    List all of the tags under the given repository.

#####history \<repository\> [\<tag\>]
    Get history infomation from given repository.

#####delete \<repository\> [\<tag\>]
    Delete tagged repository

#####size \<repository\> [\<tag\>]
    Get size infomation from given repository.

#####token \<service\> \<realm\> \<scope\>
    Create api token for docker register server

####Commands:

When the application is started it will check if there is config file in the $HOME/.docker-registry/conf.json folder. 
here you can define the registry host, username, password and verbosity. Password can be left out anf the application 
will aks for it while verifying the api.

#####Example:

{
	"user" : {
		"name": "some_user",
		"pass": "*******"
	},
	"registry-host":  "http://registry.example.com",
	"verbose": true
}

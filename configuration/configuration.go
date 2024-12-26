package configuration

import (
	"os"
)

type Configuration struct {
	Auth authenticationCredentials
	Domain string
	SpaceID string
}

type authenticationCredentials struct {
	Username string
	ApiKey   string
}

// Parses the .env for use and stores all info into a handy struct
func Init() Configuration {
	spaceID := readConfluenceSpaceIDFromEnvironmentVariable()
	auth := readAuthenticationCredentialsFromEnvironmentVariables()
	domain := readConfluenceDomainFromEnvironmentVariable()
	return Configuration{Auth: auth, Domain: domain, SpaceID: spaceID}
}

func readAuthenticationCredentialsFromEnvironmentVariables() authenticationCredentials {
	return authenticationCredentials{
		Username: os.Getenv("CONFLUENCE_USERNAME"),
		ApiKey:   os.Getenv("CONFLUENCE_API_KEY"),
	}
}

func readConfluenceDomainFromEnvironmentVariable() string {
	return os.Getenv("CONFLUENCE_DOMAIN")
}

func readConfluenceSpaceIDFromEnvironmentVariable() string {
	return os.Getenv("CONFLUENCE_SPACEID")
}

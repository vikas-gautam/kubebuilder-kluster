package doutils

import (
	"github.com/digitalocean/godo"
)

func generateDoClient(tokenValue string) *godo.Client {
	//do client with tokenValue
	doClient := godo.NewFromToken(tokenValue)
	return doClient
}

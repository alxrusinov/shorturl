package handler

import "fmt"

func createShortLink(host string, shorten string) string {
	return fmt.Sprintf("%s/%s", host, shorten)
}

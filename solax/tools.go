package solax

import (
	"net/url"
)

func UrlValid(apiAddr string) bool {
  _, err := url.Parse(apiAddr)
  return err == nil
}

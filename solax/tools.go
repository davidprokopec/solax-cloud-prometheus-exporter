package solax

import (
	"net/url"
)

func UrlValid(apiAddr string) bool {
  return true
  _, err := url.Parse(apiAddr)
  return err == nil
}

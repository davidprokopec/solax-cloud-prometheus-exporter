package solax

import (
	"net/url"
)

func UrlValid(apiAddr string) (bool, error){
  _, err := url.Parse(apiAddr)
  if(err != nil){
    return false, err
  }
  return true, nil
}

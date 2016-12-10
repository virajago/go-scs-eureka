package scseureka

import (
  "fmt"
  //"encoding/xml"
  //"crypto/tls"
)

func Register(scsInstance string) error {
  fmt.Println("Registering to SCS")
  /*
  serverURI:= []string{"https://eureka-9325402d-e8e0-49e6-95eb-35d877be7be8.apps.pcf.thy.com/eureka"}
  clientSecret := "F0MDnbdI26VW"
  clientID := "p-service-registry-df85da3d-c536-45bb-9351-d535a3b32035"
  tokenURI := "https://p-spring-cloud-services.uaa.system.pcf.thy.com/oauth/token"

  oAuth2Options:= eureka.Oauth2ClientCredentials(clientID, clientSecret, tokenURI)
  tlsconfig:= eureka.TLSConfig(&tls.Config{InsecureSkipVerify: true})

  c:=eureka.NewClient(serverURI,oAuth2Options,tlsconfig)




	data, err := xml.Marshal(instance)
	if err != nil {
		return err
	}

	return c.retry(c.do("POST", c.appPath(instance.AppName), data, http.StatusNoContent))
  */
  return nil
}

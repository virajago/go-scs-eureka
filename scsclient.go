package eureka

import (
  "fmt"
  "os"
  "encoding/json"
  "encoding/xml"
  "crypto/tls"
  "strconv"
  "time"
)

type CFApp struct {
		ID   string `json:"application_id"`
		Name string `json:"application_name"`
    URIs []string `json:"application_uris"`
    Port uint16
    IPAddr string
    InstanceID string
    Hostname string
}

type CFService struct {
		ServerURI string
    ClientSecret string
    ClientID string
    TokenURI string
}

const (
    ServiceRegistry = "p-service-registry"
)

func getAppInstanceInfo() (*CFApp, *Instance, error) {

  app:= new(CFApp)
  regInstance:= new(Instance)
  p,err := strconv.ParseUint(os.Getenv("PORT"),10,16)
  app.Port = uint16(p)
  app.IPAddr= os.Getenv("CF_INSTANCE_IP")
  app.InstanceID = os.Getenv("CF_INSTANCE_GUID")
  //hostname := os.Getenv("CF_INSTANCE_ADDR")
  app.Hostname = app.IPAddr+":"+os.Getenv("PORT")

  VCAPApplicationEnv := os.Getenv("VCAP_APPLICATION")

  fmt.Println("VCAP port",app.Port)
  fmt.Println("VCAP IPAddr",app.IPAddr)
  fmt.Println("VCAP instanceID",app.InstanceID)
  fmt.Println("Hostname",app.Hostname)
  fmt.Println("VCAP APPLICATION",VCAPApplicationEnv)

  if err = json.Unmarshal([]byte(VCAPApplicationEnv), app); err != nil {
    panic(err)
  }

  fmt.Println("ApplID",app.ID)
  fmt.Println("ApplName",app.Name)
  fmt.Println("ApplURI",app.URIs[0])
  fmt.Println("Hostname",app.Hostname)

  regInstance = &Instance{
      XMLName:        xml.Name{Local: "instance"},
      ID:             app.InstanceID,
      HostName:       app.Hostname,
      AppName:        app.Name,
      IPAddr:         app.IPAddr,
      VIPAddr:        app.IPAddr,
      SecureVIPAddr:  app.IPAddr,
      Status:         StatusUp,
      StatusOverride: StatusUnknown,
      Port:           Port(app.Port),
      SecurePort:     Port(app.Port),
      HomePageURL:    app.URIs[0],
      StatusPageURL:  app.URIs[0]+"/status",
      HealthCheckURL: app.URIs[0]+"/health",
      DataCenterInfo: DataCenter{
        Type: DataCenterTypePrivate,
      },
  }

  return app,regInstance,nil

}

func getServiceInfo() (*CFService,error){

  var services map[string]interface{}
  service:= new(CFService)

  VCAPServicesEnv := os.Getenv("VCAP_SERVICES")
  fmt.Println("VCAP SERVICES",VCAPServicesEnv)

  if err := json.Unmarshal([]byte(VCAPServicesEnv), &services); err != nil {
  panic(err)
  }
  fmt.Println(services)
  registryCred := services[ServiceRegistry].([]interface{})[0].(map[string]interface{})["credentials"].(map[string]interface{})
  service.ServerURI = registryCred["uri"].(string)+"/eureka"
  service.ClientSecret = registryCred["client_secret"].(string)
  service.ClientID = registryCred["client_id"].(string)
  service.TokenURI = registryCred["access_token_uri"].(string)

  fmt.Println("serverURI",service.ServerURI)
  fmt.Println("clientSecret",service.ClientSecret)
  fmt.Println("clientID",service.ClientID)
  fmt.Println("tokenURI",service.TokenURI)

  return service,nil
}

func GetClientSCS(skip_ssl bool) *Client {

  serviceCred,err := getServiceInfo()
  if err!=nil {
    fmt.Println("Error getting CF Service Inst env",err)
  }

  oAuth2Options := Oauth2ClientCredentials(serviceCred.ClientID, serviceCred.ClientSecret, serviceCred.TokenURI)
  var tlsOption Option
  if skip_ssl == true {
    tlsOption = TLSConfig(&tls.Config{InsecureSkipVerify: true})
  }

  c := NewClient([]string{serviceCred.ServerURI},oAuth2Options,tlsOption)

  if c==nil {
    errMsg := "Error creating HTTP client"
    fmt.Println(errMsg)
    return nil
  }

  //Use this block to test successful connectivity to SCS server
  /*regApps,err:=c.Apps()
  if err!=nil {
  fmt.Println("Error getting list of apps:",err)
  return nil
  }
  fmt.Println("No of apps: ",len(regApps))
  */

  return c
}

func RegisterSCS(skip_ssl bool) error {
  fmt.Println("Getting HTTP Client for SCS")

  c:=GetClientSCS(skip_ssl)

  _,regInstance,err:=getAppInstanceInfo()
  if err!=nil {
    fmt.Println("Error getting CF app env",err)
    return err
  }

  fmt.Println("Registering application",regInstance.AppName,"/",regInstance.ID)
  err= c.Register(regInstance)
  if err!=nil {
    fmt.Println("Error registering app:",err)
    return err
  }

  return nil
}

func SendHearbeatSCS(skip_ssl bool) {
  fmt.Println("Getting HTTP Client for SCS")

  c:=GetClientSCS(skip_ssl)

  _,regInstance,err:=getAppInstanceInfo()
  if err!=nil {
    fmt.Println("Error getting CF app env",err)
    return
  }

  for {
    fmt.Println("Sending Heartbeat for",regInstance.AppName,"/",regInstance.ID)
    err = c.Heartbeat(regInstance)
    if err!=nil {
      fmt.Println("Error sending Heartbeat:",err)
      return
    }
    time.Sleep(time.Second * 30)
  }

}

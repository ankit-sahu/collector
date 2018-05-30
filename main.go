package main

import (
	"net/http"
	"github.com/pkg/errors"
	"github.com/golang/glog"
	"flag"
	"crypto/tls"
	"io/ioutil"
	"bytes"
	"strings"
	"time"
	"github.com/collector/config"
)

type Resource struct{
	Name string
	Path string
}

func SendRequest(req *http.Request) ([]byte, error) {

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		msg := "invalid cluster url provided or cluster is not operational." +
			"Ensure that the cluster url is valid and cluster is operational. " +
			"More info : " + err.Error()
		err := errors.New(msg)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		err := errors.New("unauthorized. make sure your token is correct and valid")
		return nil, err
	}
	if resp.StatusCode == 403 {
		err := errors.New("forbidden. user does not have access to the project or project does not exist")
		return nil, err
	}
	if resp.StatusCode == 404 {
		err := errors.New("resource not found. please ensure that the resource exists")
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func PrepareUrl(url, path string) string{

	if !strings.HasSuffix(url,"/"){
		url = url + "/"
	}

	if strings.HasPrefix(path,"/"){
		path = strings.Trim(path, "/")
	}
	return url + path
}

func CollectData(url, token string) ([]byte, error){

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil,err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	data, err := SendRequest(req)
	if err != nil{
		return nil, err
	}
	return data, err
}

func PostData(url string,resourceName string, data []byte) error{

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		glog.Error("Unable to create request.")
		return err
	}
	req.Header.Add("Resource-Type", resourceName)
	req.Header.Add("Content-Type", "application/json")
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return err
}

func main(){

	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	appConfig,err := config.Config("/config/appConfig.yaml")
	if err != nil{
		glog.Error(err.Error())
	}

	for _, item := range appConfig.MetricResources{
		resource := Resource{Name:item.Name, Path:item.Path}
		stopChan := make(chan bool)
		go func(resource Resource){
			defer close(stopChan)
			ticker := time.NewTicker(time.Duration(appConfig.Interval) * time.Second)
			for {
				select{
					case <- stopChan:
						glog.Error("terminating goroutine")
					case <- ticker.C:
						data, err := CollectData(PrepareUrl(appConfig.ClusterURL,resource.Path), appConfig.Token)
						if err != nil {
							glog.Info(err)
						} else {
							err := PostData(appConfig.DestinationURL,resource.Name, data)
							if err != nil {
								glog.Info(err)
							}
						}
				}
			}
		}(resource)
		<-stopChan
	}
}


package nxrm

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/terraform/states/remote"
)

type NXRMClient struct {
	userName        string
	password        string
	url             string
	subpath         string
	timeoutSeconds  int
	tfStateArtifact string
	httpClient      *http.Client
}

func (n *NXRMClient) getNXRMURL() string {
	return fmt.Sprintf("%s/%s/%s", n.url, n.subpath, n.tfStateArtifact)
}

func (n *NXRMClient) getHTTPClient() *http.Client {
	if n.httpClient == nil {
		n.httpClient = &http.Client{
			Timeout: time.Second * time.Duration(n.timeoutSeconds),
		}
	}
	return n.httpClient
}

func (n *NXRMClient) getRequest(method string, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, n.getNXRMURL(), data)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(n.userName, n.password)

	return req, nil
}

func (n *NXRMClient) Get() (*remote.Payload, error) {
	req, err := n.getRequest(http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	res, err := n.getHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, nil
	}

	hash := md5.Sum(output)

	payload := &remote.Payload{
		Data: output,
		MD5:  hash[:md5.Size],
	}

	return payload, nil
}

func (n *NXRMClient) Put(data []byte) error {
	req, err := n.getRequest(http.MethodPut, bytes.NewReader(data))
	if err != nil {
		return err
	}

	_, err = n.getHTTPClient().Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (n *NXRMClient) Delete() error {
	req, err := n.getRequest(http.MethodDelete, nil)
	if err != nil {
		return err
	}

	n.getHTTPClient().Do(req)
	if err != nil {
		return err
	}

	return nil
}

package ssm

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

var (
	defaultRegion = "ap-southeast-2"
)

type Client struct {
	ssmiface.SSMAPI
}

func New() (*Client, error) {
	svc := ssm.New(Sess())
	return &Client{
		svc,
	}, nil
}

//Sess init new config session
func Sess() *session.Session {
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_DEFAULT_REGION")
	}
	if awsRegion == "" {
		awsRegion = defaultRegion
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(awsRegion)},
	}))
	return sess
}

// GetSecrets queries ssm parameter store for a given path and returns of map of key values.
func (c *Client) GetSecrets(parampath string) (map[string]string, error) {
	results := make(map[string]string)

	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(parampath),
		Recursive:      aws.Bool(false),
		WithDecryption: aws.Bool(true),
	}

	for {

		resp, err := c.GetParametersByPath(input)
		if err != nil {
			return nil, err
		}

		for _, param := range resp.Parameters {
			_, key := path.Split(*param.Name)
			results[key] = *param.Value
		}

		if resp.NextToken == nil {
			break
		}

	}
	return results, nil

}

// DecodeSecrets will convert from gzipped, base64 encoded values to strings.
func (c *Client) DecodeSecrets(secrets map[string]string) (map[string]string, error) {
	for k, v := range secrets {

		decoded, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			fmt.Printf("warn: base64 decode error: %s\n", err)
			continue
		}
		buf := bytes.NewBuffer(decoded)
		zr, err := gzip.NewReader(buf)
		if err != nil {
			fmt.Printf("warn: gzip reader error: %s\n", err)
			secrets[k] = string(decoded)
			continue
		}
		defer zr.Close()

		uncompressed, err := ioutil.ReadAll(zr)
		if err != nil {
			fmt.Printf("warn: gzip decompress error: %s\n", err)
			secrets[k] = string(decoded)
			continue
		}
		secrets[k] = string(uncompressed)

	}
	return secrets, nil
}

// EncodeSecrets will convert from strings to gzipped, base64 encoded values.
func (c *Client) EncodeSecrets(secrets map[string]string) (map[string]string, error) {
	for k, v := range secrets {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		defer zw.Close()

		_, err := zw.Write([]byte(v))
		if err != nil {
			fmt.Printf("warn: gzip compress error: %s\n", err)
			continue
		}
		if err := zw.Close(); err != nil {
			fmt.Printf("warn: gzip close error: %s\n", err)
			continue
		}
		secrets[k] = base64.StdEncoding.EncodeToString(buf.Bytes())

	}
	return secrets, nil
}

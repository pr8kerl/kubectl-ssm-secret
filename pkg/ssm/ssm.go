package ssm

import (
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

/*
type SSMClient interface {
	GetSecrets(string) (map[string]string, error)
}
*/

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

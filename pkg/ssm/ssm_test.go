package ssm

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/stretchr/testify/assert"
)

var (
	mockParameters []*ssm.Parameter = []*ssm.Parameter{
		&ssm.Parameter{
			ARN: aws.String("arn:aws:ssm:ap-southeast-2:012345678901:parameter/foo/passwd"),
			// LastModifiedDate: 2019-06-12 02:40:26 +0000 UTC,
			LastModifiedDate: aws.Time(time.Now()),
			Name:             aws.String("/foo/passwd"),
			Type:             aws.String("SecureString"),
			Value:            aws.String("SecretSquirrel"),
			Version:          aws.Int64(1),
		},
		&ssm.Parameter{
			ARN: aws.String("arn:aws:ssm:ap-southeast-2:012345678901:parameter/foo/username"),
			//                LastModifiedDate: 2019-06-12 02:40:27 +0000 UTC,
			LastModifiedDate: aws.Time(time.Now()),
			Name:             aws.String("/foo/username"),
			Type:             aws.String("String"),
			Value:            aws.String("Gerald"),
			Version:          aws.Int64(1),
		},
	}
)

func (m *Client) GetParametersByPath(i *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	// mock response/functionality
	return &ssm.GetParametersByPathOutput{
		Parameters: mockParameters,
		NextToken:  nil,
	}, nil
}

func TestGetSecrets(t *testing.T) {
	// Setup Test
	mockssm := Client{}

	t.Run("test GetSecrets returns expected results", func(t *testing.T) {
		secrets, err := mockssm.GetSecrets("/foo")
		assert.Nil(t, err)
		assert.Equal(t, 2, len(secrets))
		assert.Equal(t, "SecretSquirrel", secrets["passwd"], "SecretSquirrel")
		assert.Equal(t, "Gerald", secrets["username"])
	})
}

func TestSession(t *testing.T) {

	t.Run("setting region takes precedence", func(t *testing.T) {
		os.Setenv("AWS_REGION", "us-east-1")
		os.Setenv("AWS_DEFAULT_REGION", "eu-central-1")
		s := Sess()
		assert.Equal(t, "us-east-1", aws.StringValue(s.Config.Region))
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_DEFAULT_REGION")
	})

	t.Run("setting default region takes effect when region not set", func(t *testing.T) {
		os.Setenv("AWS_DEFAULT_REGION", "eu-central-1")
		s := Sess()
		assert.Equal(t, "eu-central-1", aws.StringValue(s.Config.Region))
		os.Unsetenv("AWS_DEFAULT_REGION")
	})

	t.Run("setting no region gives default fallback region", func(t *testing.T) {
		s := Sess()
		assert.Equal(t, "ap-southeast-2", aws.StringValue(s.Config.Region))
	})

}

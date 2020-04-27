package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	log "github.com/sirupsen/logrus"
)

// SecretsProvider AWS secrets provider
type SecretsProvider struct {
	session *session.Session
	sm      secretsmanageriface.SecretsManagerAPI
	ssm     ssmiface.SSMAPI
}

func NewAwsSecretsProvider() (*SecretsProvider, error) {
	var err error
	sp := SecretsProvider{}
	// create AWS session
	sp.session, err = session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return nil, err
	}
	// init AWS Secrets Manager client
	sp.sm = secretsmanager.New(sp.session)
	// init AWS SSM client
	sp.ssm = ssm.New(sp.session)
	return &sp, nil
}

// ResolveSecrets replaces all passed variables values prefixed with 'aws:' and 'aws:arn:'
// by corresponding secrets from AWS Secret Manager and AWS Parameter Store
func (sp *SecretsProvider) ResolveSecrets(vars map[string]string) map[string]string {
	parsedString := make(map[string]string)
	for key, value := range vars {
		if strings.HasPrefix(value, "aws:") {
			if strings.HasPrefix(value, "arn:aws:ssm") && strings.Contains(value, ":parameter/") {
				tokens := strings.Split(value, ":")
				// valid parameter ARN arn:aws:ssm:REGION:ACCOUNT:parameter/PATH
				if len(tokens) == 6 {
					// get SSM parameter name (path)
					paramName := strings.TrimPrefix(tokens[5], "parameter")
					// get AWS SSM API
					withDecryption := true
					param, err := sp.ssm.GetParameter(&ssm.GetParameterInput{
						Name:           &paramName,
						WithDecryption: &withDecryption,
					})
					if err != nil {
						log.Printf("Could not resolv %s. Err %s", paramName, err)
						parsedString[key] = value
					} else {
						parsedString[key] = *param.Parameter.Value
					}
				}

			} else if strings.HasPrefix(value, "arn:aws:ssm") && strings.Contains(value, ":parameter/") {
				// get secret value
				secretKeyPath := strings.TrimPrefix(value, "aws:")
				secret, err := sp.sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretKeyPath})
				if err != nil {
					log.Printf("Could not resolv %s. Err %s", secretKeyPath, err)
					parsedString[key] = value
				} else {
					parsedString[key] = *secret.SecretString
				}

			}
		} else {
			parsedString[key] = value
		}
	}

	return parsedString
}

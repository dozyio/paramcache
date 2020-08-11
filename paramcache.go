package paramcache

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

//SSMParameterStoreCache setup
type SSMParameterStoreCache struct {
	Value        *ssm.GetParameterOutput
	CacheExpires int64
}

const (
	cacheDefaultTimeout int64 = 300
)

var (
	useSSMCache    string           = "true" //override via environment var USE_SSM_CACHE
	parameterStore                  = make(map[string]SSMParameterStoreCache)
	cacheTimeout   int64            = cacheDefaultTimeout //override via environment var SSM_CACHE_TIMEOUT
	Sess           *session.Session = nil
)

func GetParameterStoreValue(param string) (*ssm.GetParameterOutput, error) {
	if val, ok := os.LookupEnv("USE_SSM_CACHE"); ok {
		useSSMCache = val
	}

	if val, ok := os.LookupEnv("SSM_CACHE_TIMEOUT"); ok {
		if timeout, err := strconv.ParseInt(val, 10, 64); err == nil {
			if timeout >= 0 {
				cacheTimeout = timeout
			}
		}
	}

	if useSSMCache == "true" || useSSMCache == "TRUE" {
		if parameter, ok := parameterStore[param]; ok {
			if time.Now().Unix() < parameter.CacheExpires {
				log.Printf("SSM Param: %s - from cache\n", param)
				return parameter.Value, nil
			}
		}
	}

	if Sess == nil {
		Sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(os.Getenv("AWS_REGION")),
			},
		}))
	}

	ssmService := ssm.New(Sess)

	paramOutput, err := ssmService.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(param),
	})
	if err != nil {
		return nil, err
	}

	if useSSMCache == "true" || useSSMCache == "TRUE" {
		log.Printf("SSM Param: %s - not from cache\n", param)
		t := time.Now()
		cacheExpires := t.Unix() + cacheTimeout
		p := &SSMParameterStoreCache{
			Value:        paramOutput,
			CacheExpires: cacheExpires,
		}
		parameterStore[param] = *p
	}
	return paramOutput, nil
}

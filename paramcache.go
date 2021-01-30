package paramcache

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	//override via environment var SSM_CACHE_ENABLED
	ssmCacheEnabled string = "true"

	//override via environment var SSM_CACHE_TIMEOUT
	cacheTimeout int64 = cacheDefaultTimeout

	//override via environment var SSM_VERBOSE
	verbose string = "false"

	//the cache store
	parameterStore = make(map[string]SSMParameterStoreCache)

	//shared session
	sess *session.Session = nil
)

//AWSSession create AWS Session for reuse
func AWSSession(s *session.Session) *session.Session {
	if s == nil {
		sessionNew := session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(os.Getenv("AWS_REGION")),
			},
		}))
		return sessionNew
	}
	return s
}

//Setup configures paramcache via environment variables
//Configurable variables are as follows
//SSM_CACHE_ENABLED (default: true)
//SSM_VERBOSE (default: false)
//SSM_CACHE_TIMEOUT (default: true)
func setup() {
	//get environment vars
	if val, ok := os.LookupEnv("SSM_CACHE_ENABLED"); ok {
		ssmCacheEnabled = val
	}

	if val, ok := os.LookupEnv("SSM_VERBOSE"); ok {
		verbose = val
	}

	if val, ok := os.LookupEnv("SSM_CACHE_TIMEOUT"); ok {
		if timeout, err := strconv.ParseInt(val, 10, 64); err == nil {
			if timeout >= 0 {
				cacheTimeout = timeout
			}
		}
	}

	//create aws session
	if sess == nil {
		sess = AWSSession(nil)
	}

}

//GetParameterStoreValue returns a string, stringlist or securestring from SSM and caches the value if configured to do so.
func GetParameterStoreValue(param string, paramType string) (*ssm.GetParameterOutput, error) {
	setup()

	//return value if already cached
	if ssmCacheEnabled == "true" || ssmCacheEnabled == "TRUE" {
		if parameter, ok := parameterStore[param]; ok {
			if time.Now().Unix() < parameter.CacheExpires {
				if verbose == "true" || verbose == "TRUE" {
					log.Printf("SSM ParamCache: %s - from cache", param)
				}
				return parameter.Value, nil
			}
		}
	}

	ssmService := ssm.New(sess)

	paramOutput, err := ssmService.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(param),
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Printf("Error: SSM ParamCache: %v", aerr.Error())
		} else {
			log.Printf("Error: SSM ParamCache: %v", err)
		}
		return nil, err
	}

	//store value in cache
	if ssmCacheEnabled == "true" || ssmCacheEnabled == "TRUE" {
		if verbose == "true" || verbose == "TRUE" {
			log.Printf("SSM ParamCache: %s - not from cache, caching for %v seconds", param, cacheTimeout)
		}

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

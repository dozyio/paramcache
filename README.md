# ParamCache - A Go SSM Parameter Store Cache

ParamCache is used in AWS Go Lambdas when cold booting to save repeated lookups on subsequent invocations.
Assumes SSM / AWS Systems Manager is in the same region as the lambda.

Cache configurable via lambda environment variables:

SSM_CACHE_ENABLE string (default "TRUE")

SSM_CACHE_VERBOSE string (default "FALSE")

SSM_CACHE_TIMEOUT int (default 300 seconds)


# Usage
```go
import "github.com/dozyio/paramcache"

...

tableName, err := paramcache.GetParameterStoreValue("dynamodb_table_name")
if err != nil {
	//handle error
}
log.Printf("%s\n", *tableName.Parameter.Value)
```

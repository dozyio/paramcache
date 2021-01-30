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

# Example: Cache timeout per parameter
It is possible to set a cache timeout value per parameter. The example below sets a 20 second cache for "dynamodb_table_name". 
```go
tableName, err := paramcache.GetParameterStoreValue("dynamodb_table_name", 20)
```

# Example: Don't cache a parameter
Setting a cache timeout of 0 will not add the parameter to the cache
```go
tableName, err := paramcache.GetParameterStoreValue("dynamodb_table_name", 0)
```

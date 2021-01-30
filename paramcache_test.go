package paramcache

import (
	"os"
	"strconv"
	"testing"
)

//use eu-west-2 for testing

func TestSetupSSMCacheDefault(t *testing.T) {
	os.Unsetenv("SSM_CACHE_ENABLED")
	setup()
	if cacheEnabled != "TRUE" {
		t.Errorf("want TRUE, got %v", cacheEnabled)
	}
}

func TestSetupSSMCacheDisabled(t *testing.T) {
	os.Unsetenv("SSM_CACHE_ENABLED")
	want := "FALSE"
	os.Setenv("SSM_CACHE_ENABLED", want)
	setup()
	if cacheEnabled != want {
		t.Errorf("want %v, got %v", want, cacheEnabled)
	}
}

func TestSetupSSMCacheEnabled(t *testing.T) {
	os.Unsetenv("SSM_CACHE_ENABLED")
	want := "TRUE"
	os.Setenv("SSM_CACHE_ENABLED", want)
	setup()
	if cacheEnabled != want {
		t.Errorf("want %v, got %v", want, cacheEnabled)
	}
}

func TestSetupVerboseDefault(t *testing.T) {
	os.Unsetenv("SSM_VERBOSE")
	want := "FALSE"
	setup()
	if verbose != want {
		t.Errorf("want %v, got %v", want, verbose)
	}
}

func TestSetupVerboseEnabled(t *testing.T) {
	os.Unsetenv("SSM_VERBOSE")
	want := "TRUE"
	os.Setenv("SSM_VERBOSE", want)
	setup()
	if verbose != want {
		t.Errorf("want %v, got %v", want, verbose)
	}
}

func TestSetupVerboseDisabled(t *testing.T) {
	os.Unsetenv("SSM_VERBOSE")
	want := "FALSE"
	os.Setenv("SSM_VERBOSE", want)
	setup()
	if verbose != want {
		t.Errorf("want %v, got %v", want, verbose)
	}
}

func TestSetupCacheTimeoutDefault(t *testing.T) {
	os.Unsetenv("SSM_CACHE_TIMEOUT")
	var want int64 = 300
	setup()
	if cacheTimeout != want {
		t.Errorf("want %v, got %v", want, cacheTimeout)
	}
}

func TestSetupTimeout(t *testing.T) {
	os.Unsetenv("SSM_CACHE_TIMEOUT")
	var want int64 = 100
	os.Setenv("SSM_CACHE_TIMEOUT", strconv.FormatInt(int64(want), 10))
	setup()
	if cacheTimeout != want {
		t.Errorf("want %v, got %v", want, cacheTimeout)
	}
}

func TestAWSSessionDefault(t *testing.T) {
	setup()
	sess := AWSSession(sess)
	if sess == nil {
		t.Errorf("want != nil, got %v", sess)
	}
}

//Requires ssm parameter paramcache_test_string = teststring
func TestGetParameterStoreValueString(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")

	want := "teststring"
	got, err := GetParameterStoreValue("paramcache_test_string")

	if err != nil {
		t.Errorf("Error %v", err)
	} else {
		if *got.Parameter.Value != want {
			t.Errorf("want %v, got %v", want, *got.Parameter.Value)
		}
	}
}

//Requires ssm parameter paramcache_test_string_list = teststring1,teststring2
func TestGetParameterStoreValueStringList(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")

	want := "teststring1,teststring2"
	got, err := GetParameterStoreValue("paramcache_test_string_list")

	if err != nil {
		t.Errorf("Error %v", err)
	} else {
		if *got.Parameter.Value != want {
			t.Errorf("want %v, got %v", want, *got.Parameter.Value)
		}
	}
}

//Requires ssm parameter paramcache_test_secure_string = testsecurestring
func TestGetParameterStoreValueSecureString(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")

	want := "testsecurestring"
	got, err := GetParameterStoreValue("paramcache_test_secure_string")

	if err != nil {
		t.Errorf("Error %v", err)
	} else {
		if *got.Parameter.Value != want {
			t.Errorf("want %v, got %v", want, *got.Parameter.Value)
		}
	}
}

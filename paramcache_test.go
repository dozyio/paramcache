package paramcache

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

//use eu-west-2 for testing

func TestSetupSSMCacheDefault(t *testing.T) {
	os.Unsetenv("SSM_CACHE_ENABLED")

	setup()
	if cacheEnabled != "TRUE" {
		t.Errorf("want TRUE, got %v", cacheEnabled)
	}

	os.Unsetenv("SSM_CACHE_ENABLED")
}

func TestSetupSSMCacheDisabled(t *testing.T) {
	os.Unsetenv("SSM_CACHE_ENABLED")

	want := "FALSE"
	os.Setenv("SSM_CACHE_ENABLED", want)
	setup()
	if cacheEnabled != want {
		t.Errorf("want %v, got %v", want, cacheEnabled)
	}

	os.Unsetenv("SSM_CACHE_ENABLED")
}

func TestSetupSSMCacheEnabled(t *testing.T) {
	os.Unsetenv("SSM_CACHE_ENABLED")

	want := "TRUE"
	os.Setenv("SSM_CACHE_ENABLED", want)
	setup()
	if cacheEnabled != want {
		t.Errorf("want %v, got %v", want, cacheEnabled)
	}

	os.Unsetenv("SSM_CACHE_ENABLED")
}

func TestSetupVerboseDefault(t *testing.T) {
	os.Unsetenv("SSM_VERBOSE")

	want := "FALSE"
	setup()
	if verbose != want {
		t.Errorf("want %v, got %v", want, verbose)
	}

	os.Unsetenv("SSM_VERBOSE")
}

func TestSetupVerboseEnabled(t *testing.T) {
	os.Unsetenv("SSM_VERBOSE")

	want := "TRUE"
	os.Setenv("SSM_VERBOSE", want)
	setup()
	if verbose != want {
		t.Errorf("want %v, got %v", want, verbose)
	}

	os.Unsetenv("SSM_VERBOSE")
}

func TestSetupVerboseDisabled(t *testing.T) {
	os.Unsetenv("SSM_VERBOSE")

	want := "FALSE"
	os.Setenv("SSM_VERBOSE", want)
	setup()
	if verbose != want {
		t.Errorf("want %v, got %v", want, verbose)
	}

	os.Unsetenv("SSM_VERBOSE")
}

func TestSetupCacheTimeoutDefault(t *testing.T) {
	os.Unsetenv("SSM_CACHE_TIMEOUT")

	var want int64 = 300
	setup()
	if cacheTimeout != want {
		t.Errorf("want %v, got %v", want, cacheTimeout)
	}

	os.Unsetenv("SSM_CACHE_TIMEOUT")
}

func TestSetupTimeout(t *testing.T) {
	os.Unsetenv("SSM_CACHE_TIMEOUT")

	var want int64 = 100
	os.Setenv("SSM_CACHE_TIMEOUT", strconv.FormatInt(int64(want), 10))
	setup()
	if cacheTimeout != want {
		t.Errorf("want %v, got %v", want, cacheTimeout)
	}

	os.Unsetenv("SSM_CACHE_TIMEOUT")
}

func TestAWSSessionDefault(t *testing.T) {
	setup()
	sess := AWSSession(sess)
	if sess == nil {
		t.Errorf("want != nil, got %v", sess)
	}
	sess = nil
}

//Requires ssm parameter paramcache_test_string = teststring
func TestGetParameterStoreValueString(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")
	os.Setenv("SSM_VERBOSE", "TRUE")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	want := "teststring"
	got, err := GetParameterStoreValue("paramcache_test_string")

	if err != nil {
		t.Errorf("Error %v", err)
	} else {
		if *got.Parameter.Value != want {
			t.Errorf("want %v, got %v", want, *got.Parameter.Value)
		}
	}

	os.Unsetenv("AWS_REGION")
}

//Requires ssm parameter paramcache_test_string_list = teststring1,teststring2
func TestGetParameterStoreValueStringList(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")
	os.Setenv("SSM_VERBOSE", "TRUE")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

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
func TestGetParameterStoreValueNotFound(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")

	got, err := GetParameterStoreValue("paramcache_does_not_exist")
	if err != nil {
		if !strings.Contains(buf.String(), "ParameterNotFound") {
			t.Errorf("ParameterNotFound not found in error. %v", buf.String())
		}
	} else {
		t.Errorf("want error, got %v", got)
	}
}

//Requires ssm parameter paramcache_test_timeout = timeouttest
func TestGetParameterStoreValueTimeout(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")
	os.Setenv("SSM_VERBOSE", "TRUE")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	want := "timeouttest"
	timeout := 20
	got, err := GetParameterStoreValue("paramcache_test_timeout", timeout)

	if err != nil {
		t.Errorf("Error %v", err)
	} else {
		if *got.Parameter.Value != want {
			t.Errorf("want %v, got %v", want, *got.Parameter.Value)
		}
		if !strings.Contains(buf.String(), "20 seconds") {
			t.Errorf("output didn't contain '%v seconds', got %v", timeout, buf.String())
		}
	}

	os.Unsetenv("AWS_REGION")
}

//Requires ssm parameter paramcache_test_string = teststring
func TestGetParameterStoreValueCache(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")
	os.Setenv("SSM_VERBOSE", "FALSE")

	want := "teststring"
	_, err := GetParameterStoreValue("paramcache_test_string")

	os.Setenv("SSM_VERBOSE", "TRUE")
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	got2, err := GetParameterStoreValue("paramcache_test_string")

	if err != nil {
		t.Errorf("Error %v", err)
	} else {
		if *got2.Parameter.Value != want {
			t.Errorf("want %v, got %v", want, *got2.Parameter.Value)
		}
		if !strings.Contains(buf.String(), "paramcache_test_string - from cache") {
			t.Errorf("want 'paramcache_test_string - from cache', got %v", buf.String())
		}
	}

	os.Unsetenv("AWS_REGION")
}

//Requires ssm parameter paramcache_test_string = teststring
func TestGetParameterStoreValue0Cache(t *testing.T) {
	sess = nil
	os.Setenv("AWS_REGION", "eu-west-2")
	os.Setenv("SSM_VERBOSE", "TRUE")

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	want := "teststring"
	got, err := GetParameterStoreValue("paramcache_test_string", 0)

	if err != nil {
		t.Errorf("Error %v", err)
	} else {
		if *got.Parameter.Value != want {
			t.Errorf("want %v, got %v", want, *got.Parameter.Value)
		}
		if strings.Contains(buf.String(), "from cache") {
			t.Errorf("output includes 'from cache', got %v", buf.String())
		}
	}

	os.Unsetenv("AWS_REGION")
}

func TestGetTimeoutDefault(t *testing.T) {
	os.Unsetenv("SSM_CACHE_TIMEOUT")
	var want int64 = 300
	os.Setenv("SSM_CACHE_TIMEOUT", strconv.FormatInt(int64(want), 10))
	setup()

	got := getTimeout()

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGetTimeout(t *testing.T) {
	want := 20
	got := getTimeout(want)

	if int64(want) != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

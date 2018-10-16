package lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mbndr/logo"
	"strings"
	"testing"
	"time"
)

func TestNewAssumeRoleProvider(t *testing.T) {
	t.Run("ProfileNil", func(t *testing.T) {
		defer func() {
			if x := recover(); x == nil {
				t.Errorf("Did not receive expected panic calling NewAssumeRoleProvider with nil profile")
			}
		}()
		NewAssumeRoleProvider(nil, new(CachedCredentialsProviderOptions))
	})

	t.Run("OptionsNil", func(t *testing.T) {
		p := NewAssumeRoleProvider(new(AWSProfile), nil)
		cache := p.(*assumeRoleProvider).cacher.CacheFile()
		if !strings.HasSuffix(cache, "_") {
			t.Errorf("Unexpected value returned calling NewAssumeRoleProvider with nil options: %v", cache)
		}
	})
}

func TestAssumeRoleProvider_IsExpired(t *testing.T) {
	opts := &CachedCredentialsProviderOptions{LogLevel: logo.DEBUG}
	t.Run("CredsNil", func(t *testing.T) {
		p := NewAssumeRoleProvider(new(AWSProfile), opts)
		if !p.IsExpired() {
			t.Errorf("Expected IsExpired() to be true for nil creds")
		}
	})

	t.Run("True", func(t *testing.T) {
		p := NewAssumeRoleProvider(new(AWSProfile), opts)
		p.(*assumeRoleProvider).cacher = &credentialsCacher{file: "config/test/cached_creds_expired.json"}
		if !p.IsExpired() {
			t.Errorf("Expected IsExpired() to be true for expired creds")
		}
	})

	t.Run("False", func(t *testing.T) {
		p := NewAssumeRoleProvider(new(AWSProfile), opts)
		p.(*assumeRoleProvider).cacher = &credentialsCacher{file: "config/test/cached_creds_valid.json"}
		if p.IsExpired() {
			t.Errorf("Expected IsExpired() to be false for non-expired creds")
		}
	})
}

func TestAssumeRoleProvider_ValidateSessionName(t *testing.T) {
	p := NewAssumeRoleProvider(new(AWSProfile), new(CachedCredentialsProviderOptions)).(*assumeRoleProvider)

	t.Run("NameEmpty", func(t *testing.T) {
		v := p.validateSessionName("")
		if !strings.HasPrefix(*v, "AWS-RUNAS-") {
			t.Errorf("Did not receive expected session name, expected: AWS-RUNAS-, got: %s", *v)
		}
	})

	t.Run("NameSet", func(t *testing.T) {
		n := "mock"
		v := p.validateSessionName(n)
		if *v != n {
			t.Errorf("Did not receive expected session name, expected: %s, got: %s", n, *v)
		}
	})
}

func TestAssumeRoleProvider_ValidateDuration(t *testing.T) {
	p := NewAssumeRoleProvider(new(AWSProfile), new(CachedCredentialsProviderOptions)).(*assumeRoleProvider)

	t.Run("DurationZero", func(t *testing.T) {
		i := p.validateDuration(0)
		if *i != *aws.Int64(int64(ASSUME_ROLE_DEFAULT_DURATION.Seconds())) {
			t.Errorf("Expected default duration, got %d", *i)
		}
	})

	t.Run("DurationHigh", func(t *testing.T) {
		i := p.validateDuration(72 * time.Hour)
		if *i != *aws.Int64(int64(ASSUME_ROLE_MAX_DURATION.Seconds())) {
			t.Errorf("Expected max duration, got %d", *i)
		}
	})

	t.Run("DurationLow", func(t *testing.T) {
		i := p.validateDuration(5 * time.Minute)
		if *i != *aws.Int64(int64(ASSUME_ROLE_MIN_DURATION.Seconds())) {
			t.Errorf("Expected min duration, got %d", *i)
		}
	})

	t.Run("DurationNormal", func(t *testing.T) {
		d := 2 * time.Hour
		i := p.validateDuration(d)
		if *i != *aws.Int64(int64(d.Seconds())) {
			t.Errorf("Expected original duration, got %d", *i)
		}
	})
}

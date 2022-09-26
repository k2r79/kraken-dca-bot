package domain

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	err := os.Setenv("KRAKEN_KEY", "fake_kraken_key")
	if err != nil {
		t.Errorf("%v", err)
	}
	err = os.Setenv("KRAKEN_SECRET", "fake_kraken_secret")
	if err != nil {
		t.Errorf("%v", err)
	}
	err = os.Setenv("EMAIL_PASSWORD", "fake_email_password")
	if err != nil {
		t.Errorf("%v", err)
	}

	config, err := ParseConfig("../../test/data/sample-config.yaml")
	if err != nil {
		t.Errorf("An unexpected error occurred : %v", err)
	}

	expectedConfig := Config{
		Kraken: Kraken{
			Key:    "fake_kraken_key",
			Secret: "fake_kraken_secret",
		},
		Smtp: Smtp{
			Host:     "smtp.google.com",
			Port:     587,
			User:     "smtp_user",
			Password: "fake_email_password",
			From:     "sender@gmail.com",
		},
		Notify:    "recipient@gmail.com",
		Frequency: "1ms",
		Currency:  "ZEUR",
		Pairs: []DCAPair{
			{"XETHZEUR", 20.00},
			{"XXBTZEUR", 10.00},
			{"XXRPZEUR", 10.00},
			{"ADAEUR", 10.00},
			{"USDTEUR", 10.00},
		},
	}

	if !reflect.DeepEqual(*config, expectedConfig) {
		t.Errorf("The parsed configuration is %v instead of %v", *config, expectedConfig)
	}
}

func TestParseConfigOpenFileFail(t *testing.T) {
	_, err := ParseConfig("../../test/data/does-not-exist.yaml")
	if err == nil || !strings.HasPrefix(err.Error(), "an error occurred while trying to open the configuration file :") {
		t.Errorf("An unexpected error occurred : %v", err)
	}
}

func TestParseConfigReadFileFail(t *testing.T) {
	_, err := ParseConfig("../../test/data/invalid-yaml.yaml")
	if err == nil || !strings.HasPrefix(err.Error(), "an error occurred while unmarshalling the configuration Yaml :") {
		t.Errorf("An unexpected error occurred : %v", err)
	}
}

func TestParseConfigEmptyKrakenKeyFail(t *testing.T) {
	_, err := ParseConfig("../../test/data/empty-kraken-key.yaml")
	if err == nil || err.Error() != "the kraken key is not specified" {
		t.Errorf("An unexpected error occurred : %v", err)
	}
}

func TestParseConfigEmptyKrakenSecretFail(t *testing.T) {
	_, err := ParseConfig("../../test/data/empty-kraken-secret.yaml")
	if err == nil || err.Error() != "the kraken secret is not specified" {
		t.Errorf("An unexpected error occurred : %v", err)
	}
}

package appsettings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

// ReadFromFileAndEnv reads the settings from a local file and applies any existing environment variables to it.
func ReadFromFileAndEnv(settings interface{}) error {
	file, err := os.Open("anton.json")
	log.Print(file)
	if err != nil {
		return err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)

	if err != nil {
		return errors.Wrap(err, "Failed to read appsettings")
	}

	err = json.Unmarshal(data, settings)

	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal appsettings")
	}

	err = envconfig.Init(settings)

	if err != nil {
		err = errors.Wrap(err, "Failed to update with env vars")
	}
	return err
}

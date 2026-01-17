package env

import (
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	envparse "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const (
	appEnvironment = "APP_ENV"
	devEnvironment = "development"
)

var (
	goVersion string

	revision         string
	revisionShort    string
	revisionTime     time.Time
	revisionModified bool
)

func init() {
	_ = godotenv.Load()

	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion = info.GoVersion

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
				if len(setting.Value) > 7 {
					revisionShort = setting.Value[:7]
				} else {
					revisionShort = setting.Value
				}
			case "vcs.time":
				revisionTime, _ = time.Parse(time.RFC3339, setting.Value)
			case "vcs.modified":
				revisionModified = setting.Value == "true"
			}
		}
	}
}

// Parse parses environment variables into the provided struct.
// See github.com/caarlos0/env for tag documentation.
//
// Additionally supports buildinfo tags for build information:
//
//	type Config struct {
//	    GoVersion string    `buildinfo:"goversion"`
//	    Revision  string    `buildinfo:"revision"`
//	    Short     string    `buildinfo:"revisionshort"`
//	    Time      time.Time `buildinfo:"revisiontime"`
//	    Modified  bool      `buildinfo:"revisionmodified"`
//	}
func Parse(v any) error {
	if err := envparse.Parse(v); err != nil {
		return err
	}
	return parseBuildInfo(v)
}

// Get returns the value of the specified environment variable
func Get(key string) string {
	return os.Getenv(key)
}

// GetDefault returns the value of the specified environment variable
// If the variable is not defined or is empty, return the default value
func GetDefault(key, defaultValue string) string {
	value := Get(key)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

// AppEnvironment returns the name of the current environment the app is in
func AppEnvironment() string {
	return strings.ToLower(Get(appEnvironment))
}

// IsDevelopment returns true if the current environment is development
func IsDevelopment() bool {
	return AppEnvironment() == devEnvironment
}

// IsProduction returns if the current environment is production
func IsProduction() bool {
	return !IsDevelopment()
}

// GoVersion returns the Go version used to build the application
func GoVersion() string {
	return goVersion
}

// Revision returns the full VCS revision used to build the application
func Revision() string {
	return revision
}

// RevisionShort returns the short VCS revision used to build the application
func RevisionShort() string {
	return revisionShort
}

// RevisionTime returns the VCS revision time used to build the application
func RevisionTime() time.Time {
	return revisionTime
}

// RevisionModified returns true if there were uncommitted changes
// when the application was built
func RevisionModified() bool {
	return revisionModified
}

func parseBuildInfo(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return nil
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	for i := range val.NumField() {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := typ.Field(i).Tag.Get("buildinfo")
		if tag == "" {
			continue
		}

		switch strings.ToLower(tag) {
		case "goversion":
			if field.Kind() == reflect.String {
				field.SetString(goVersion)
			}
		case "revision":
			if field.Kind() == reflect.String {
				field.SetString(revision)
			}
		case "revisionshort":
			if field.Kind() == reflect.String {
				field.SetString(revisionShort)
			}
		case "revisiontime":
			if field.Type() == reflect.TypeFor[time.Time]() {
				field.Set(reflect.ValueOf(revisionTime))
			}
		case "revisionmodified":
			if field.Kind() == reflect.Bool {
				field.SetBool(revisionModified)
			}
		}
	}
	return nil
}

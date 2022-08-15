package env

import "os"

const Prefix = "BTH_"

// Get retrieves the value of the environment variable named
// by the key with prefix (PREFIX + key). If the variable is present in the environment the
// value (which may be empty) is returned.
// Otherwise, the returned value will be fallback value.
func Get(key string, fallback string) string {
	val, found := os.LookupEnv(Prefix + key)
	if !found {
		return fallback
	}
	return val
}

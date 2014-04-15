package gofig

import (
	"os"
	"encoding/gob"
)

// IgnoreCache will force re-parsing of the configuration file
//
func (config Configuration ) IgnoreCache( flag bool ) Configuration {
	config.ignoreCache = flag
	return config
}
// SetCache Set the cache filename. This will notalter the data.
// To change the data, call LoadCache()
func (config Configuration) SetCache(cache string) Configuration {
	config.cacheFile = cache
	config.IsCache = false
	return config
}

// SaveCache save the contents of the configuration, unconditionally
// This will take the cacheFile entry (if set) and write the contents
// of the configuration out.
//
func (config *Configuration) SaveCache() error {
	if config.cacheFile != "" && config.IsLoaded {
		cache, err := os.Create(config.cacheFile)
		if err != nil {
			return err
		}
		defer cache.Close()
		enc := gob.NewEncoder(cache)
		enc.Encode(config)
	}
	return nil
}

// LoadCache using the configuration, load unconditionally load the GOB
// config file
func (config *Configuration) LoadCache() (*Configuration, error) {
	var newC Configuration

	cache, err := os.Open(config.cacheFile)
	if err != nil {
		return config, err
	}
	defer cache.Close()
	dec := gob.NewDecoder(cache)
	err = dec.Decode(&newC)
	if err != nil {
		return config, err
	}
	config.ConfigMap = newC.ConfigMap
	config.IsLoaded = true
	config.IsCache = true
	return config, nil

}

// Copyright 2014 Charles Gentry All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gofig implements a simple configuration parsing system.
// It allows you to load a section-based configuration into a simple map.
// all acceses should occur through the routines within this package
//
// This package uses the convention of 'functionName' to describe a function
// and, if it has additional options available, it may have an function
// named 'functionNameWithOptions'. So, NewConfigurationFromIniFile will
// load the configuration file from a standard INI-style file.
// NewConfigurationFromIniFileWithCache will include a try at cache loading.
//
// Example ini file:
//   ; Comment
//   # Also a comment. Blank lines will also be ignored
//   #
//   [ db ]
//   db="postgres"
//   host = remotehost
//   #
//   [testdb : db ]
//   host=localhost
//   #
// testdb will inherit from db. You can have as man inheritence as you like. They are
// evaluated from left to right.
// Spaces and quotes are stripped and ignored.
// if you want to have quotes, you can double them: ""value"" will give "value"
//
// To get an option, you would call GetString( "testdb" , "db" )
// To get a numeric option, you would use GetInt. For booleans, use GetBool
//
package gofig

import (
	"bufio"
	"errors"
	//"fmt"
	"io"
	"os"
	"strings"
	"strconv"
)

const (
	defaultPreAllocate = 10
)

// ConfigOption is a single map level for key => value pair
type ConfigOption map[string]string // Single line config

// Configuration
// This structure holds all of the configuration information for
// the system. Each config is a simple map that can be 'backstored'
// into a GOB cache file.
//
// GOB cache files are automatically created when requested and used
// if the GOB is newer than the human-readable file

//
type Configuration struct {

	// Contains the section => option map
	ConfigMap map[string]ConfigOption

	// How many sections are filled in
	Sections int

	// True if we want options that are 'defaulted' to be added into each section
	OnDefaultAddToSection bool

	// True if the contents of the ConfigMap are valid and loaded
	IsLoaded bool

	// True if a cache was used for the data
	IsCache bool

	// True if you want the cache contents to always be ignored
	ignoreCache bool

	// Source of the configuration
	ConfigFile string

	// Name of the cache file used
	cacheFile string
}

// OnDefaultAddToSection will set the flag to determine if we should add values into each section
// when we use a default value.
func ( config *Configuration ) SetAddOnDefault( flag bool ) *Configuration {
	config.OnDefaultAddToSection = flag
	return config
}
// IsCacheFileNewer will check to see if a cache file is newer than the main
// file. If the file doesnt exist, it will be considered 'older'
func ( config Configuration ) IsCacheFileNewer( ) bool {
	if "" != config.cacheFile {
		return false
	}
	return isCacheFileNewer( config.ConfigFile , config.cacheFile);
}

// isCacheFileNewer will return true IF the cache file is newer or equal to file.
func isCacheFileNewer(file, cacheFile string) bool {
	fileInfo, fileErr := os.Stat(file)
	cacheInfo, cacheErr := os.Stat(cacheFile)
	if fileErr != nil || cacheErr != nil || cacheInfo.Size() == 0 {
		return false
	}
	return (cacheInfo.ModTime().After( fileInfo.ModTime() ) )
}

// NewConfigurationFromCache Create a cache file from the cache string
// return error when the file cannot be opened
func NewConfigurationFromCache(cacheFile string) (*Configuration, error) {

	newConfig := NewConfigurationWithCache(cacheFile)

	_, err := newConfig.LoadCache()   // Load from the file...
	newConfig.IsLoaded = (err == nil) // Set the flag for load status
	return newConfig, err             // return the configuration
}

// New ConfigurationWithCache Create a new configuration with a cache
// file name set. This does not load the cache, but only creates it
func NewConfigurationWithCache(file string) *Configuration {
	config := Configuration{
		ConfigMap:  make(map[string]ConfigOption, defaultPreAllocate ),
		IsLoaded:   false,
		IsCache:    false,
		ConfigFile: "",
		cacheFile:  file,
	}
	config.AddSection( "_default")			// Where we store defaulted values
	return &config
}

// NewConfiguration will return a completely initialised configuration
// but with no cache setup
func NewConfiguration() *Configuration {
	return NewConfigurationWithCache("")
}

// configFromReader is the internal function that does the actual
// parsing required for the sections and options.
func configFromReader(reader io.Reader, cacheName string) (*Configuration, error) {
	scanner := bufio.NewScanner(reader)
	var line string
	section := "default"
	config := NewConfigurationWithCache(cacheName)

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		lenLine := len(line)
		// If we have a non-comment, non-blank line...
		if lenLine > 0 && line[0:1] != "#" && line[0:1] != ";" {
			if line[0:1] == "[" {
				if line[lenLine-1:] != "]" {
					return nil, errors.New("Invalid section marker in line: " + line)
				}
				section = strings.TrimSpace(line[1 : lenLine-1])
				// Find out if there are any subsections (inheritance)
				if strings.Contains(section, ":") {

					for i, name := range strings.Split(section, ":") {
						sectionName := strings.TrimSpace(name)
						if i == 0 {
							section = sectionName
						} else {
							config.MergeOptions(section, sectionName)
						}
					}

				}
			} else {
				parts := strings.SplitN(line, "=" , 2)
				if len(parts) != 2 {
					return nil, errors.New("Invalid key/value pair in line: " + line)
				}
				config.SetString(section, parts[0], parts[1] )
			}
		}
	}
	config.IsLoaded = true
	config.SaveCache()
	return config, nil
}

// NewConfigurationFromIniString will create a new configuration from a
// string rather than using a file. Caching is not used with strings
func NewConfigurationFromIniString(input string) (*Configuration, error) {
	if input == "" {
		return nil, errors.New("String cannot be empty")
	}
	return configFromReader(strings.NewReader(input), "")
}

// NewConfigurationFromIniFile will create a new configuration, read in the
// the standard ini-style configuration file and return a configuration
func NewConfigurationFromIniFileWithCache(filename, cache string) (*Configuration, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if cache != "" && isCacheFileNewer(filename, cache) {
		config, err := NewConfigurationFromCache(cache)
		if err == nil {
			return config, err
		}
	}
	return configFromReader(file, cache)
}

// NewConfigurationFromIniFile will open up a filename and parse the ini-style
// strings from each line found.
func NewConfigurationFromIniFile(filename string) (*Configuration, error) {
	return NewConfigurationFromIniFileWithCache(filename, "")
}

// String will convert the configuration into a nicely printable,
// indented format
func (config Configuration) String() string {
	var out string

	out = out + ";\n;  OnDefaultAddToSection is " + strconv.FormatBool( config.OnDefaultAddToSection ) + "\n;\n"
	for key := range config.ConfigMap {
		out = out + "[" + key + "]\n"
		for subkey := range config.ConfigMap[key] {
			out = out + "\t'" + subkey + "' = '" + config.ConfigMap[key][subkey] + "'\n"
		}
	}

	return out

}




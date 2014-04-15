gofig
=====

Simple configuration package for INI-type files for the GO language.

PACKAGE DOCUMENTATION

package gofig
    import "github.com/CGentry/gofig"

    Package gofig implements a simple configuration parsing system. It
    allows you to load a section-based configuration into a simple map. all
    acceses should occur through the routines within this package

    This package uses the convention of 'functionName' to describe a
    function and, if it has additional options available, it may have an
    function named 'functionNameWithOptions'. So,
    NewConfigurationFromIniFile will load the configuration file from a
    standard INI-style file. NewConfigurationFromIniFileWithCache will
    include a try at cache loading.

    Example ini file:

	; Comment
	# Also a comment. Blank lines will also be ignored
	#
	[ db ]
	db="postgres"
	host = remotehost
	#
	[testdb : db ]
	host=localhost
	#

    testdb will inherit from db. You can have as man inheritence as you
    like. They are evaluated from left to right. Spaces and quotes are
    stripped and ignored. if you want to have quotes, you can double them:
    ""value"" will give "value"

    To get an option, you would call GetString( "testdb" , "db" ) To get a
    numeric option, you would use GetInt. For booleans, use GetBool

    Copyright 2014 Charles Gentry All Rights Reserved. Use of this source
    code is governed by a BSD-style license that can be found in the LICENSE
    file.


TYPES

type ConfigOption map[string]string // Single line config

    ConfigOption is a single map level for key => value pair



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

    // Source of the configuration
    ConfigFile string
    // contains filtered or unexported fields
}


func NewConfiguration() *Configuration
    NewConfiguration will return a completely initialised configuration but
    with no cache setup


func NewConfigurationFromCache(cacheFile string) (*Configuration, error)
    NewConfigurationFromCache Create a cache file from the cache string
    return error when the file cannot be opened


func NewConfigurationFromIniFile(filename string) (*Configuration, error)
    NewConfigurationFromIniFile will open up a filename and parse the
    ini-style strings from each line found.


func NewConfigurationFromIniFileWithCache(filename, cache string) (*Configuration, error)
    NewConfigurationFromIniFile will create a new configuration, read in the
    the standard ini-style configuration file and return a configuration


func NewConfigurationFromIniString(input string) (*Configuration, error)
    NewConfigurationFromIniString will create a new configuration from a
    string rather than using a file. Caching is not used with strings


func NewConfigurationWithCache(file string) *Configuration
    New ConfigurationWithCache Create a new configuration with a cache file
    name set. This does not load the cache, but only creates it


func (config *Configuration) AddSection(sectionName string) ConfigOption
    AddSection will return a map for a given section. If the section doesn't
    exist, it will be created

func (config *Configuration) DeleteOption(sectionName, optionName string) *Configuration
    DeleteOption will delete an option within a given section. If the
    section or option doesn't exist, the request will be ignored

func (config *Configuration) DeleteSection(sectionName string) *Configuration
    DeleteSection Delete all of the entries in a section. If the map doesn't
    exist, ignore it

func (config *Configuration) GetBool(sectionName, optionName string) (bool, error)
    GetBool will return an boolean value of the string, converted Boolean
    values must conform to the strconv.ParseBool values (1/0,
    true/false,etc) The values 'yes' and 'no' do not work

func (config *Configuration) GetBoolWithDefault(sectionName, optionName string, defaultValue bool) (bool, error)
    GetBoolWithDefault will return an bool or the default value if nothing
    is available

func (config *Configuration) GetInt(sectionName, optionName string) (int64, error)
    GetInt will return an int64 value of the number, convertered

func (config *Configuration) GetIntWithDefault(sectionName, optionName string, defaultValue int64) (int64, error)
    GetIntWithDefault will return an int64 or the default value if nothing
    is available

func (config *Configuration) GetSection(sectionName string) (ConfigOption, bool)
    GetSection will return a map for a given section and true or return nil
    and false

func (config *Configuration) GetSectionNames() []string
    GetSectionNames will return a complete list of all of the sections
    defined

func (config *Configuration) GetString(sectionName, optionName string) (string, error)
    GetString will search a section for a specific option. If the option or
    section doesn't exist, an error will be returned.

func (config *Configuration) GetStringWithDefault(sectionName, optionName, defaultValue string) string
    GetStringWithDefault will search a section for a specific option. If the
    section or option doesn't exist, a default value (passed into the
    routine) will be returned instead.

func (config Configuration) IgnoreCache(flag bool) Configuration
    IgnoreCache will force re-parsing of the configuration file

func (config Configuration) IsCacheFileNewer() bool
    IsCacheFileNewer will check to see if a cache file is newer than the
    main file. If the file doesnt exist, it will be considered 'older'

func (config *Configuration) IsSection(sectionName string) bool
    IsSection will return true if a section name exists in the map

func (config *Configuration) LoadCache() (*Configuration, error)
    LoadCache using the configuration, load unconditionally load the GOB
    config file

func (config *Configuration) MergeOptions(targetSection, sourceSection string) *Configuration
    MergeOptions This will merge a source section into a target section.
    This is normally only used with the cascade-style functions in a
    section: [A : B : C]

func (config *Configuration) SaveCache() error
    SaveCache save the contents of the configuration, unconditionally This
    will take the cacheFile entry (if set) and write the contents of the
    configuration out.

func (config *Configuration) SetAddOnDefault(flag bool) *Configuration
    OnDefaultAddToSection will set the flag to determine if we should add
    values into each section when we use a default value.

func (config Configuration) SetCache(cache string) Configuration
    SetCache Set the cache filename. This will notalter the data. To change
    the data, call LoadCache()

func (config *Configuration) SetString(sectionName, optionName, value string)
    SetString will insert an option and value into a section. If the section
    doesn't exist, it will be created

func (config Configuration) String() string
    String will convert the configuration into a nicely printable, indented
    format




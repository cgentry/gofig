// Copyright 2014 Charles Gentry All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
package gofig

import (
	"errors"
	"strings"
	"strconv"
	//"fmt"
)


func conformOption( option string ) string {
	// Get rid of spaces surrounding option
	option = strings.TrimSpace(option)
	optLen := len(option)

	// Remove matching quote marks ("......")
	if optLen > 0 {
		if option[0:1] == "'" || option[0:1] == `"` {
			if  option[0:1] == option[optLen-1:]{
				option = option[1:optLen-1]
			}
		}
	}
	return option
}

// DeleteOption will delete an option within a given section. If the section or
// option doesn't exist, the request will be ignored
func (config *Configuration) DeleteOption(sectionName, optionName string) *Configuration {
	mm, ok := config.ConfigMap[sectionName]
	if ok {
		delete(mm, optionName)
	}
	return config
}


// MergeOptions This will merge a source section into a target section. This is
// normally only used with the cascade-style functions in a section: [A : B : C]
func (config *Configuration) MergeOptions(targetSection, sourceSection string) *Configuration {
	ts := config.AddSection(targetSection)
	if ss, found := config.GetSection(sourceSection); found {
		for key, value := range ss {
			ts[key] = value
		}
	}
	return config

}

// SetString will insert an option and value into a section. If the section
// doesn't exist, it will be created
func (config *Configuration) SetString(sectionName, optionName, value string) {

	mm := config.AddSection(sectionName)
	mm[conformOption(optionName)] = conformOption(value)
}

// IsOption return true if a section and option exists in the config
func (config *Configuration ) IsOption( sectionName, optionName  string ) bool {
        mm , found := config.ConfigMap[sectionName]
	if found {
		_ , found = mm[optionName]
	}
        return found
}

// GetString will search a section for a specific option. If the option
// or section doesn't exist, an error will be returned.
func (config *Configuration) GetString(sectionName, optionName string) (string, error) {

	mm, ok := config.GetSection( sectionName )
	if !ok {
		return "", errors.New("Section '" + sectionName + "' not found")
	}
	value, ok := mm[optionName]
	if !ok {
		return "", errors.New("Option '" + optionName + "' not found")
	}

	return value, nil
}

// GetStringWithDefault will search a section for a specific option. If the
// section or option doesn't exist, a default value (passed into the routine)
// will be returned instead.
func (config *Configuration) GetStringWithDefault(sectionName, optionName, defaultValue string) ( string ) {
	value, error := config.GetString(sectionName, optionName)
	if error != nil {
		config.SetString("_default", sectionName, optionName)
		if config.OnDefaultAddToSection {
			config.SetString( sectionName , optionName , defaultValue)
		}
		return defaultValue
	}
	return value
}

// GetInt will return an int64 value of the number, convertered
func (config *Configuration) GetInt(sectionName , optionName string ) (int64 , error ){
	mm , err := config.GetString( sectionName , optionName )
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt( mm , 0 , 64 )
}

// GetIntWithDefault will return an int64 or the default value if nothing is available
func (config *Configuration) GetIntWithDefault( sectionName, optionName string , defaultValue int64) (int64 , error){
	i := strconv.FormatInt( defaultValue , 10 )
	mm  := config.GetStringWithDefault( sectionName , optionName , i )
	return strconv.ParseInt( mm , 0 , 64 )
}



// GetBool will return an boolean value of the string, converted
// Boolean values must conform to the strconv.ParseBool values (1/0, true/false,etc)
// The values 'yes' and 'no' do not work
func (config *Configuration) GetBool(sectionName , optionName string ) (bool , error ){
	mm , err := config.GetString( sectionName , optionName )
	if err != nil {
		return false, err
	}
	return strconv.ParseBool( mm  )
}

// GetBoolWithDefault will return an bool or the default value if nothing is available
func (config *Configuration) GetBoolWithDefault( sectionName, optionName string , defaultValue bool) (bool , error){
	mm , err := config.GetString(sectionName, optionName)
	if err != nil {
		return defaultValue,nil
	}
	return strconv.ParseBool( mm  )
}

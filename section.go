// Copyright 2014 Charles Gentry All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package gofig

import (
	"strings"
	// "fmt"
)




func conformSectionName( sectionName string ) string {
	return strings.TrimSpace( sectionName )
}
// IsSection will return true if a section name exists in the map
func (config *Configuration ) IsSection( sectionName string ) bool {
	_ , found := config.ConfigMap[sectionName]
	return found
}
// GetSection will return a map for a given section and true  or return nil and false
func (config *Configuration) GetSection(sectionName string) (ConfigOption, bool) {
	mm, found := config.ConfigMap[conformSectionName( sectionName )]
	return mm, found
}

// AddSection will return a map for a given section. If the section
// doesn't exist, it will be created
func (config *Configuration) AddSection(sectionName string) ConfigOption {
	sectionName  = conformSectionName( sectionName )
	mm, ok := config.ConfigMap[sectionName]
	if !ok {
		mm = make(ConfigOption, 10)
		config.ConfigMap[sectionName] = mm
		config.Sections++
	}
	return mm
}

// DeleteSection Delete all of the entries in a section. If the map doesn't exist, ignore it
func (config *Configuration) DeleteSection( sectionName string ) *Configuration {
	sectionName  = conformSectionName( sectionName )
	if _ , found := config.ConfigMap[sectionName] ; found {
		for opt,_ := range config.ConfigMap[sectionName] {
			delete( config.ConfigMap[sectionName] , opt )
		}
		delete( config.ConfigMap , sectionName)
		config.Sections--
	}
	return config
}

// GetSectionNames will return a complete list of all of the sections defined
func (config *Configuration) GetSectionNames() []string {

	names := make( []string , len(config.ConfigMap ) )

	i := 0
	for name := range config.ConfigMap {
		if name != "" {
			names[i] = name
			i++
		}
	}
	return names
}

// Configuration package.
// allows you to load a section-based configuration into a simple map.
// all acceses should occur through the routines within this package
package gofig

import (
	"testing"
	//"fmt"
)

var testdata_set1 = `
; comment 1
       # another comment
[begin]
a=1
b=  2
c="  3"

[  middle   ]
  key1 = value 1

# Begining of the end
[end]
key2 = value 2
`

var testdata_cascade = `
[one]
a=1
b=2

[two]
c=3
d=4

# Test an inherited value
[three:two:one]
e=5`

var testdata_types =`
[types]
bool=true
int=101
string=hello

badbool=yep
badint=10.654
`

var testdata_badsection=`
[section
int=10`

var testdata_badsection2=`
]section]
int=10`

func checkSection( t *testing.T , c * Configuration , section , option , shouldBe string ) bool {
	if ! c.IsSection( section ){
		t.Errorf("Section '%s' does not exist" , section )
		return false
	}
	opt , err := c.GetString( section , option )
	if err == nil {
		if opt != shouldBe {
			t.Errorf( "[%s] %s should be '%s' but is '%s'\n" , section, option , shouldBe , opt )
			return false
		}
		return true
	}
	t.Errorf("Section [%s]  %s\n", section , err )
	return false
}


func TestIni( t *testing.T ) {
	config, err := NewConfigurationFromIniFileWithCache("tst.ini", "/tmp/tst.gob")
	if err != nil {
		t.Errorf("Error from ini:")
	}

	sections := config.GetSectionNames()
	if len( sections )  !=4 {
		t.Errorf("Not enough sections: %d" , len( sections ) )
	}
}

func TestIniWithString( t *testing.T ){
	config,err := NewConfigurationFromIniString( testdata_set1 )

	if err != nil {
		t.Errorf( "Error: could not ini from string")
	}

	sections := config.GetSectionNames()
	if len( sections ) != 4 {
		t.Errorf("Sections should be 4 but is %d\n" , len( sections ))
	}
}
func TestIni_CheckStringValues( t *testing.T ){
	config,err := NewConfigurationFromIniString( testdata_set1 )

	if err != nil {
		t.Errorf( "Error: could not ini from string")
	}


	checkSection( t , config , "begin" , "a" , "1")
	checkSection( t , config , "begin" , "b" , "2")
	checkSection( t , config , "begin" , "c" , "  3")
	checkSection( t , config , "middle", "key1" , "value 1")
	checkSection( t , config , "end"   , "key2" , "value 2")
}
func TestIni_CheckConvertedValues( t *testing.T ){
	config,err := NewConfigurationFromIniString( testdata_set1 )

	if err != nil {
		t.Errorf( "Error: could not ini from string")
	}


	i,err := config.GetInt( "begin" , "a")
	if i != 1 || err != nil {
		t.Errorf("Expected an integer '1'")
	}

	defaultString := config.GetStringWithDefault("begin" , "nota" , "20")
	if defaultString != "20" {
		t.Errorf( "Expected default of 20 but got %s" , defaultString )
	}

	defaultInt,_ := config.GetIntWithDefault("begin" , "nota" , 22)
	if defaultInt != 22 {
		t.Errorf( "Expected default of 20 but got %d" , defaultInt )
	}

	// Check to see if we got the default defined
	defaultSet, err := config.GetString( "_default" , "begin" )
	if defaultSet != "nota"{
		t.Errorf("default was not set for nota: %s" , err)
	}

	config.SetAddOnDefault(true)
	defaultInt,_ = config.GetIntWithDefault( "end" , "nota" , 22 )
	defaultSet, err = config.GetString( "end" , "nota" )
	if err != nil {
		t.Errorf( "Did not save value: %s" , err)
	}
	if defaultSet != "22" {
		t.Errorf("Did not set default value %s" , defaultSet)
	}
}

func TestIsSection( t *testing.T ){
	config,_ := NewConfigurationFromIniString( testdata_set1 )
	if config.IsSection( "not-there"){
		t.Errorf( "Section 'not-there' was found?")
	}
	if ! config.IsSection( "begin" ){
		t.Errorf( "Section 'begin' should be there but wasn't Found")
	}
}

func TestCascadeSections( t *testing.T){
	config,err := NewConfigurationFromIniString( testdata_cascade )

	if err != nil {
		t.Errorf( "Error: could not ini from string")
	}
	checkSection( t , config , "one"    , "a" , "1")
	checkSection( t , config , "one"    , "b" , "2")
	checkSection( t , config , "two"    , "c" , "3")
	checkSection( t , config , "two"    , "d" , "4")
	checkSection( t , config , "three"  , "e" , "5")
	checkSection( t , config , "three"    , "a" , "1")
	checkSection( t , config , "three"    , "b" , "2")
	checkSection( t , config , "three"    , "c" , "3")
	checkSection( t , config , "three"    , "d" , "4")

}

func TestTypes( t *testing.T){
	config,err := NewConfigurationFromIniString( testdata_types )

	if err != nil {
		t.Errorf( "Error: could not ini from string")
	}

	// --------- Booleans
	b,err := config.GetBool("types" , "bool")
	if b!= true {
		t.Errorf( "Boolean value was wrong: %b"  , b )
	}

	if err != nil {
		t.Errorf( "Boolean returned an error: %s" , err )
	}

	b,err = config.GetBool( "types" , "badbool")
	if err == nil {
		t.Errorf("Badboolean value did not trigger an error")
	}
	// -------- INTEGERS
	i,err2 := config.GetInt("types" , "int")
	if i != 101 {
		t.Errorf( "Integer value was wrong: %d"  , i )
	}

	if err2 != nil {
		t.Errorf( "Boolean returned an error: %s" , err )
	}

	i,err2 = config.GetInt( "types" , "badint")
	if err2 == nil {
		t.Errorf("Badint value did not trigger an error")
	}
}
func TestBadSection( t * testing.T ){
	_,err := NewConfigurationFromIniString( testdata_badsection )

	if err == nil {
		t.Errorf( "Invalid section did not trigger an error")
	}
}

func TestBadSectionMarker( t * testing.T ){
	_,err := NewConfigurationFromIniString( testdata_badsection2 )

	if err == nil {
		t.Errorf( "Invalid section did not trigger an error")
	}
}
func TestAddSection_And_GetSectionNames( t *testing.T ){
	config := NewConfiguration()

	config.AddSection("TEST")
	sections := config.GetSectionNames()
	if len( sections )  !=2 {
		t.Errorf("Not enough sections: %d" , len( sections ) )
	}
	found := 0
	for _,name := range sections {
		if name == "_default" || name == "TEST" {
		found++
	}
	}
	if found != 2 {
		t.Errorf("Not all sections returned")
	}
}

func TestDeleteSection( t *testing.T ){
	config := NewConfiguration()
	config.AddSection("TEST")
	config.DeleteSection("TEST")

	sections := config.GetSectionNames()
	if len( sections ) != 1 {
		t.Errorf("Too many sections - delete did not work")
	}
	if sections[0] != "_default" {
		t.Errorf("Default section does not exist")
	}
}

// ------------- BENCHMARK ---------------
func BenchmarkLoadFile( b *testing.B ){
	for i:=0; i<b.N ; i++ {
		 NewConfigurationFromIniFile("tst.ini")
	}
}

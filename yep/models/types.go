// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Date type that JSON marshal and unmarshals as "YYYY-MM-DD"
type Date time.Time

// IsNull returns true if the Date is the zero value
func (d Date) IsNull() bool {
	if time.Time(d).Format("2006-01-02") == "0001-01-01" {
		return true
	}
	return false
}

// MarshalJSON for Date type
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsNull() {
		return []byte("null"), nil
	}
	dateStr := time.Time(d).Format("2006-01-02")
	dateStr = fmt.Sprintf(`"%s"`, dateStr)
	return []byte(dateStr), nil
}

// Value formats our Date for storing in database
// Especially handles empty Date.
func (d Date) Value() (driver.Value, error) {
	if d.IsNull() {
		return driver.Value("0001-01-01"), nil
	}
	return driver.Value(d), nil
}

// Today() returns the current date
func Today() Date {
	return Date(time.Now())
}

// DateTime type that JSON marshals and unmarshals as "YYYY-MM-DD HH:MM:SS"
type DateTime time.Time

// IsNull returns true if the DateTime is the zero value
func (d DateTime) IsNull() bool {
	if time.Time(d).Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		return true
	}
	return false
}

// Now() returns the current date/time
func Now() DateTime {
	return DateTime(time.Now())
}

// MarshalJSON for DateTime type
func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.IsNull() {
		return []byte("null"), nil
	}
	dateStr := time.Time(d).Format("2006-01-02 15:04:05")
	dateStr = fmt.Sprintf(`"%s"`, dateStr)
	return []byte(dateStr), nil
}

// Value formats our DateTime for storing in database
// Especially handles empty DateTime.
func (d DateTime) Value() (driver.Value, error) {
	if d.IsNull() {
		return driver.Value("0001-01-01 00:00:00"), nil
	}
	return driver.Value(time.Time(d).Format("2006-01-02 15:04:05")), nil
}

// FieldMap is a map of interface{} specifically used for holding model
// fields values.
type FieldMap map[string]interface{}

// Keys returns the FieldMap keys as a slice of strings
func (fm FieldMap) Keys() (res []string) {
	for k := range fm {
		res = append(res, k)
	}
	return
}

// Values returns the FieldMap values as a slice of interface{}
func (fm FieldMap) Values() (res []interface{}) {
	for _, v := range fm {
		res = append(res, v)
	}
	return
}

// RemovePK removes the entries of our FieldMap which
// references the ID field.
func (fm *FieldMap) RemovePK() {
	delete(*fm, "id")
	delete(*fm, "ID")
}

// RemovePKIfZero removes the entries of our FieldMap which
// references the ID field if the referenced id is 0.
func (fm *FieldMap) RemovePKIfZero() {
	if idl, ok := (*fm)["id"]; ok && idl.(int64) == 0 {
		delete(*fm, "id")
	}
	if idu, ok := (*fm)["ID"]; ok && idu.(int64) == 0 {
		delete(*fm, "ID")
	}
}

// KeySubstitution defines a key substitution in a FieldMap
type KeySubstitution struct {
	Orig string
	New  string
	Keep bool
}

// SubstituteKeys changes the column names of the given field map with the
// given substitutions.
func (fm *FieldMap) SubstituteKeys(substs []KeySubstitution) {
	for _, subs := range substs {
		value, exists := (*fm)[subs.Orig]
		if exists {
			if !subs.Keep {
				delete(*fm, subs.Orig)
			}
			(*fm)[subs.New] = value
		}
	}
}

// RecordIDWithName is a tuple with an ID and the display name of a record
type RecordIDWithName struct {
	ID   int64
	Name string
}

// MarshalJSON for RecordIDWithName type
func (rf RecordIDWithName) MarshalJSON() ([]byte, error) {
	arr := [2]interface{}{
		0: rf.ID,
		1: rf.Name,
	}
	res, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

// UnmarshalJSON for RecordIDWithName type
func (rf *RecordIDWithName) UnmarshalJSON(data []byte) error {
	var arr [2]interface{}
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	rf.ID = arr[0].(int64)
	rf.Name = arr[1].(string)
	return nil
}

// A RecordRef uniquely identifies a Record by giving its model and ID.
type RecordRef struct {
	ModelName string
	ID        int64
}

// A Selection is a set of possible (key, label) values for a model
// "selection" field.
type Selection map[string]string

// RecordSet identifies a type that holds a set of records of
// a given model.
type RecordSet interface {
	// ModelName returns the name of the model of this RecordSet
	ModelName() string
	// Ids returns the ids in this set of Records
	Ids() []int64
	// Env returns the current Environment of this RecordSet
	Env() Environment
}

// A ModelName is a type representing model names.
type ModelName string

// A FieldName is a type representing field names in models.
type FieldName string

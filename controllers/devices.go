/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mainflux/mainflux-core-server/db"
	"github.com/mainflux/mainflux-core-server/models"

	"github.com/satori/go.uuid"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/mgo.v2/bson"
)

func validateJsonSchema(b map[string]interface{}) bool {
	pwd, _ := os.Getwd()
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + pwd +
		"/models/deviceSchema.json")
	bodyLoader := gojsonschema.NewGoLoader(b)
	result, err := gojsonschema.Validate(schemaLoader, bodyLoader)
	if err != nil {
		log.Print(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
		return true
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return false
	}
}

/** == Functions == */
/**
 * CreateDevice ()
 */
func CreateDevice(b map[string]interface{}) string {
	if validateJsonSchema(b) != true {
		println("Invalid schema")
	}

	// Init new Mongo session
	// and get the "devices" collection
	// from this new session
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()

	// Turn map into a JSON to put it in the Device struct later
	j, err := json.Marshal(&b)
	if err != nil {
		fmt.Println(err)
	}

	// Set up defaults and pick up new values from user-provided JSON
	d := models.Device{Id: "Some Id", Name: "Some Name"}
	json.Unmarshal(j, &d)

	// Creating UUID Version 4
	uuid := uuid.NewV4()
	fmt.Println(uuid.String())

	d.Id = uuid.String()

	// Insert Device
	erri := Db.C("devices").Insert(d)
	if erri != nil {
		println("CANNOT INSERT")
		panic(erri)
	}

	return "Created Device req.deviceId"
}

/**
 * GetDevices()
 */
func GetDevices() string {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()

	results := []models.Device{}
	err := Db.C("devices").Find(nil).All(&results)
	if err != nil {
		log.Print(err)
	}

	r, err := json.Marshal(results)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(r)
}

/**
 * GetDevice()
 */
func GetDevice(id string) string {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()

	result := models.Device{}
	err := Db.C("devices").Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Print(err)
	}

	r, err := json.Marshal(result)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(r)
}

/**
 * UpdateDevice()
 */
func UpdateDevice(id string, b map[string]interface{}) string {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()

	// Validate JSON schema user provided
	if validateJsonSchema(b) != true {
		println("Invalid schema")
	}

	// Check if someone is trying to change "id" key
	// and protect us from this
	if _, ok := b["id"]; ok {
		println("Error: can not change device ID")
	}

	colQuerier := bson.M{"id": id}
	change := bson.M{"$set": b}
	err := Db.C("devices").Update(colQuerier, change)
	if err != nil {
		log.Print(err)
	}

	return string(`{"status":"updated"}`)
}

/**
 * DeleteDevice()
 */
func DeleteDevice(id string) string {
	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()

	err := Db.C("devices").Remove(bson.M{"id": id})
	if err != nil {
		log.Print(err)
	}

	return string(`{"status":"deleted"}`)
}

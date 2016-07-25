package main

import(
    "encoding/json"
    "fmt"
    "log"

    "github.com/satori/go.uuid"
    "gopkg.in/mgo.v2/bson"
)

type Channel struct {
    Id string   `json: "id"`
    Device string `json: "device"`
    Created string `json: "created"`
    Updated string `json: "updated"`

    Bt int `json: "bt"`
    Bn string `json: "bn"`
    Bu string `json: "bu"`
    Ver int `json: "ver"`
    E []map[string]interface{} `json: "e"`

    Msg map[string]interface{} `json: "msg"`

    Metadata map[string]interface{} `json: "metadata"`
    Mfprivate map[string]interface{} `json: "mfprivate"`
}



/** == Functions == */

/**
func insertTs(id string, ts map[string]interface{}) string {
    // Insert in SML in Influx
    // SML can contain several datapoints
    // and can target different tags
}

func insertMsg(id string, msg map[string]interface{}) string {
    // Insert Msg in Influx
    // Check if we can insert one message blob as a single datapoint
}

func queryTs() {
    // Query given measurement (datapoint) in Influx
    // Use some limit - for example 1k results
    // We need time limmits also - FROM(t=X) and UNTIL(t=y)
}

func queryMsg() {
    // Query message blobs
    // Retrieve for example last 1k messages
}
*/


/**
 * createChannel ()
 */
func createChannel(b map[string]interface{}) string {
    if validateJsonSchema(b) != true {
        println("Invalid schema")
    }

    // Turn map into a JSON to put it in the Device struct later
    j, err := json.Marshal(&b)
    if err != nil {
        fmt.Println(err)
    }

    // Set up defaults and pick up new values from user-provided JSON
    c := Channel{Id: "Some Id"}
    json.Unmarshal(j, &c)

    // Creating UUID Version 4
    uuid := uuid.NewV4()
    fmt.Println(uuid.String())

    c.Id = uuid.String()

    // Insert Device
    erri := mc.dColl.Insert(c)
	  if erri != nil {
        println("CANNOT INSERT")
		panic(erri)
	}

    return "Created Device req.deviceId"
}

/**
 * getChannels()
 */
func getChannels(id string) string {
    results := []Channel{}
    err := mc.cColl.Find(nil).All(&results)
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
 * getChannel()
 */
func getChannel(id string) string {
    result := Channel{}
    err := mc.cColl.Find(bson.M{"Id": id}).One(&result)
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
 * updateChannel()
 */
func updateChannel(id string, b map[string]interface{}) string {
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
    err := mc.cColl.Update(colQuerier, change)
    if err != nil {
        log.Print(err)
    }

    return string(`{"status":"updated"}`)
}

/**
 * deleteChannel()
 */
func deleteChannel(id string) string {
    err := mc.cColl.Remove(bson.M{"id": id})
    if err != nil {
        log.Print(err)
    }

    return string(`{"status":"deleted"}`)
}


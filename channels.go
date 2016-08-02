package main

import(
    "encoding/json"
    "fmt"
    "log"

    "github.com/satori/go.uuid"
    "gopkg.in/mgo.v2/bson"
    "github.com/xeipuuv/gojsonschema"
)

type SenML struct {
    Bt int `json: "bt"`
    Bn string `json: "bn"`
    Bu string `json: "bu"`
    Ver int `json: "ver"`
    E []map[string]interface{} `json: "e"`
}

type Channel struct {
    Id string   `json: "id"`
    Device string `json: "device"`
    Created string `json: "created"`
    Updated string `json: "updated"`

    Ts SenML `json: "ts"`

    Msg map[string]interface{} `json: "msg"`

    Metadata map[string]interface{} `json: "metadata"`
    Mfprivate map[string]interface{} `json: "mfprivate"`
}


/** == Functions == */

func insertTs(id string, ts SenML) string {
    // Insert in SenML in Influx
    // SenML can contain several datapoints
    // and can target different tags

    // Loop here for each attribute
    for k, v := range ts.E {
        tags := map[string]string{
            "attribute" : ts.E[k]["n"],
        }

        // Examine if "v" exists, then "sv", then "bv"
        var field map[string]interface{}
        if (v["v"]){
            field["value"] = v["v"]
        } else if (v["sv"]) {
            field["value"] = v["sv"]
        } else if (v["bv"]) {
            field["value"] = v["bv"]
        }

        pt, err := ic.NewPoint(id + "-ts", tags, field, ts.Bt + v["t"])

        if err != nil {
            log.Fatalln("Error: ", err)
        }

        ibp.AddPoint(pt)
    }

    // Write the batch
    ic.Write(ibp)
}

func insertMsg(id string, msg map[string]interface{}) string {
    // Insert Msg in Influx
    // Check if we can insert one message blob as a single datapoint
    tags := map[string]string{
        "attribute": "custom",
    }
    ic.NewPoint(id + "-msg", tags, msg, time.Now())
}

// queryDB convenience function to query the database
func queryInfluxDb(clnt client.Client, cmd string) (res []client.Result, err error) {
    q := ic.Query{
        Command:  cmd,
        Database: "Mainflux",
    }
    if response, err := ic.Query(q); err == nil {
        if response.Error() != nil {
            return res, response.Error()
        }
        res = response.Results
    } else {
        return res, err
    }
    return res, nil
}


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
    c := Channel{Id: "Some Id", Name: "Some Name"}
    json.Unmarshal(j, &c)

    // Creating UUID Version 4
    uuid := uuid.NewV4()
    fmt.Println(uuid.String())

    d.Id = uuid.String()

    // Insert Device
    erri := mc.dColl.Insert(d)
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

/**
 * sendChannel()
 */
func sendChannel(b map[string]interface{}) string {
    err := mc.cColl.Remove(bson.M{"id": id})
    if err != nil {
        log.Print(err)
    }

    return string(`{"status":"deleted"}`)
}

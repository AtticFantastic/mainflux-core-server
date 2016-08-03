package main

import(
    "encoding/json"
    "fmt"
    "log"
    "time"
    "math"

    "github.com/satori/go.uuid"
    "gopkg.in/mgo.v2/bson"
    "github.com/influxdata/influxdb/client/v2"
)

type SenML struct {
    Bt float64 `json: "bt"`
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

func insertTs(id string, ts SenML) int {
    rc := 0

    // Insert in SenML in Influx
    // SenML can contain several datapoints
    // and can target different tags

    // Loop here for each attribute
    for k, v := range ts.E {
        tags := map[string]string{
            "attribute" : ts.E[k]["n"].(string),
        }

        // Examine if "v" exists, then "sv", then "bv"
        var field map[string]interface{}
        if vv, okv := v["v"]; okv{
            field["value"] = vv
        } else if vsv, oksv := v["sv"]; oksv {
            field["value"] = vsv
        } else if vbv, okbv := v["bv"]; okbv {
            field["value"] = vbv
        }

        /**
         * Handle time
         * 
         * If either the Base Time or Time value is missing, the missing
         * attribute is considered to have a value of zero.  The Base Time and
         * Time values are added together to get the time of measurement.  A
         * time of zero indicates that the sensor does not know the absolute
         * time and the measurement was made roughly "now".  A negative value is
         * used to indicate seconds in the past from roughly "now".  A positive
         * value is used to indicate the number of seconds, excluding leap
         * seconds, since the start of the year 1970 in UTC.
         */
        // Set time base
        var tb float64
        if bt := ts.Bt; bt != 0.0 {
            // If bt is sent and is different than zero
            // N.B. if bt was not sent, `ts.Bt` will still be zero, as this is init value
            tb = bt
        } else {
            // If not that means that sensor does not have RTC
            // and want us to use our NTP - "roughly now"
            tb = float64(time.Now().Unix())
        }

        // Set relative time
        var tr int64
        if vt, okvt := v["t"]; okvt {
            // If there is relative time, use it
            tr = vt.(int64)
        } else {
            // Otherwise it is considered as zero
            tr = 0
        }

        // Total time
        tt := tb + float64(tr)
        // Break into int and fractional nb
        ts, tsf := math.Modf(tt)
        // Find nanoseconds number from fractional part
        tns := tsf * 1000 * 1000

        // Get time in Unix format, based on s and ns
        t := time.Unix(int64(ts), int64(tns))
        pt, err := client.NewPoint(id, tags, field, t)

        if err != nil {
            log.Fatalln("Error: ", err)
        }

        ic.bp.AddPoint(pt)
    }

    // Write the batch
    ic.c.Write(ic.bp)

    return rc
}

func insertMsg(id string, msg map[string]interface{}) int {
    rc := 0

    // Insert Msg in Influx
    // Check if we can insert one message blob as a single datapoint
    tags := map[string]string{
        "attribute": "msg",
    }
    client.NewPoint(id, tags, msg, time.Now())
    return rc
}

// queryDB convenience function to query the database
func queryInfluxDb(clnt client.Client, cmd string) (res []client.Result, err error) {
    q := client.Query{
        Command:  cmd,
        Database: "Mainflux",
    }
    if response, err := ic.c.Query(q); err == nil {
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

/**
 * sendChannel()
 */
func sendChannel(id string, b map[string]interface{}) string {
    if m, okm := b["Msg"]; okm {
        insertMsg(id, m.(map[string]interface{}))
    }

    if t, okt := b["Ts"]; okt {
        insertTs(id, t.(SenML))
    }

    return string(`{"status":"inserted"}`)
}

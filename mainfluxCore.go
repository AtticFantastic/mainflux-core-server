/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import(
    "encoding/json"
    "fmt"
    "log"
    "runtime"
    "strconv"

    "github.com/nats-io/nats"
    "gopkg.in/mgo.v2"
    "github.com/influxdata/influxdb/client/v2"

    "github.com/fatih/color"
    "github.com/spf13/viper"
)


type MainfluxMessage struct {
    Method string `json: "method"`
    Id string `json: "id"`
    Body map[string]interface{} `json: "body"`
}

/**
 * MongoDB Globals
 */
type MongoConn struct {
    session *mgo.Session
    dColl *mgo.Collection
    cColl *mgo.Collection
}

var mc MongoConn

type InfluxConn struct {
    c client.Client
    bp client.BatchPoints
}

var ic InfluxConn

/**
 * main()
 */
func main() {

    /**
     * Config
     */
     // Viper setup
    viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
    viper.SetConfigName("config") // name of config file (without extension)
    viper.AddConfigPath(".")   // path to look for the config file in
    err := viper.ReadInConfig() // Find and read the config file
    if err != nil { // Handle errors reading the config file
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    mnghost := viper.GetString("mongo.host")
    mngport := viper.GetInt("mongo.port")
    mngdb := viper.GetString("mongo.db")
    ifxhost := viper.GetString("influx.host")
    ifxport := viper.GetInt("influx.port")
    ifxdb := viper.GetString("influx.db")
    ntshost := viper.GetString("nats.host")
    ntsport := viper.GetInt("nats.port")

    /**
     * MongoDB
     */
     mgoSession, err := mgo.Dial("mongodb://" + mnghost + ":" + strconv.Itoa(mngport))
    if err != nil {
            panic(err)
    }
    //defer mgoSession.Close()

    // Optional. Switch the session to a monotonic behavior.
    mgoSession.SetMode(mgo.Monotonic, true)

    deviceMongo := mgoSession.DB(mngdb).C("devices")
    channelMongo := mgoSession.DB(mngdb).C("channels")

    /** Set-up globals */
    mc.session = mgoSession
    mc.dColl = deviceMongo
    mc.cColl = channelMongo

    /**
     * InfluxDB
     */
    // Make client
    icc, err := client.NewHTTPClient(client.HTTPConfig{
        Addr: "http://" + ifxhost + ":" + strconv.Itoa(ifxport),
        //Username: username,
        //Password: password,
    })

    if err != nil {
        log.Fatalln("Error: ", err)
    }

    // Create a new point batch
    icbp, err := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  ifxdb,
        Precision: "s",
    })

    if err != nil {
        log.Fatalln("Error: ", err)
    }

    ic.c = icc
    ic.bp = icbp


    /**
     * NATS
     */
     nc, err := nats.Connect("nats://" + ntshost + ":" + strconv.Itoa(ntsport))
	  if err != nil {
        log.Fatalf("Can't connect: %v\n", err)
	  }

    // Req-Reply
    nc.Subscribe("core_in", func(msg *nats.Msg) {
        var mfMsg MainfluxMessage

        log.Println(msg.Subject, string(msg.Data))

        // Unmarshal the message recieved from NATS
        err := json.Unmarshal(msg.Data, &mfMsg)
        if err != nil {
		      fmt.Println("error:", err)
        }

        fmt.Println(mfMsg)

        var res string
        switch mfMsg.Method {
            // Status
            case "getStatus":
                res = getStatus()
            // Devices
            case "createDevice":
                res = createDevice(mfMsg.Body)
            case "getDevices":
                res = getDevices()
            case "getDevice":
                res = getDevice(mfMsg.Id)
            case "updateDevice":
                res = updateDevice(mfMsg.Id, mfMsg.Body)
            case "deleteDevice":
                res = deleteDevice(mfMsg.Id)
            // Channels
            case "createChannel":
                res = createChannel(mfMsg.Body)
            case "getChannels":
                res = getChannels()
            case "getChannel":
                res = getChannel(mfMsg.Id)
            case "updateChannel":
                res = updateChannel(mfMsg.Id, mfMsg.Body)
            case "deleteChannel":
                res = deleteChannel(mfMsg.Id)
            default:
                fmt.Println("error: Unknown method!")
        }

        fmt.Println(res)
        nc.Publish(msg.Reply, []byte(res))
    })

	  log.Println("Listening on 'core_in'")

    color.Magenta(banner)

    /** Keep mainf() runnig */
	  runtime.Goexit()
}

var banner = `
_|      _|            _|                _|_|  _|                      
_|_|  _|_|    _|_|_|      _|_|_|      _|      _|  _|    _|  _|    _|  
_|  _|  _|  _|    _|  _|  _|    _|  _|_|_|_|  _|  _|    _|    _|_|    
_|      _|  _|    _|  _|  _|    _|    _|      _|  _|    _|  _|    _|  
_|      _|    _|_|_|  _|  _|    _|    _|      _|    _|_|_|  _|    _|  
                                                                      
    
                == Industrial IoT System ==
       
                Made with <3 by Mainflux Team

[w] http://mainflux.io
[t] @mainflux

                     ** CORE SERVER **
`

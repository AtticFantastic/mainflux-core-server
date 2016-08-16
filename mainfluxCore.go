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
    "os"

    "github.com/mainflux/mainflux-core-server/config"
    "github.com/mainflux/mainflux-core-server/db"
    "github.com/mainflux/mainflux-core-server/controllers"
    "github.com/mainflux/mainflux-core-server/broker"

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
 * main()
 */
func main() {

    // Parse config
    var cfg config.Config
    cfg.Parse()

    // MongoDb
    db.initMongo(cfg.MongoHost, cfg.MongoPort, cfg.MongoDatabase)
    Mdb := db.MgoDb{}
	  Mdb.Init()

    // InfluxDb
    db.initInflux(cfg.InfluxHost, cfg.InfluxPort, cfg.InfluxDatabase)

    // NATS
    broker.init(cfg.NatsHost, cfg.NatsPort)

    // Req-Reply
    broker.NatsConn.Subscribe("core_in", func(msg *nats.Msg) {
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
        broker.NatsConn.Publish(msg.Reply, []byte(res))
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

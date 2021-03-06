/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package broker

import (
	"github.com/nats-io/nats"
	"log"
	"strconv"
)

var (
	NatsConn *nats.Conn
)

func InitNats(host string, port int) error {
	/** Connect to NATS broker */
	var err error
	NatsConn, err = nats.Connect("nats://" + host + ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("NATS: Can't connect: %v\n", err)
	}

	return err
}

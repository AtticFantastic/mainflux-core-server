/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package models

type (
	DeviceLocation struct {
		Name      string `json: "name"`
		Latitude  int    `json: "latitude"`
		Longitude int    `json: "longitude"`
		Elevation int    `json: "elevation"`
	}

	Device struct {
		Id   string `json: "id"`
		Name string `json: "name"`

		Description string         `json: "name"`
		Visibility  string         `json: "name"`
		Status      string         `json: "name"`
		Tags        []string       `json: "name"`
		Location    DeviceLocation `json: "location"`

		Created string `json: "created"`
		Updated string `json: "updated"`

		Metadata  map[string]interface{} `json: "metadata"`
		Mfprivate map[string]interface{} `json: "mfprivate"`
	}
)

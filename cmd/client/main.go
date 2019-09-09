/*
 * Copyright 2019, authproxy authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"fmt"
	"github.com/cbrgm/authproxy/client"
	"github.com/urfave/cli"
	"os"
)

const (
	FlagTLSCert     = "tls-cert"
	FlagTLSKey      = "tls-key"
	FlagTLSClientCA = "tls-ca-cert"
)

type clientConf struct {
	TLSCert     string
	TLSKey      string
	TLSClientCA string
}

var (
	clientConfig = clientConf{}

	clientFlags = []cli.Flag{
		cli.StringFlag{
			Name:        FlagTLSKey,
			Usage:       "The tls key file to be used",
			Destination: &clientConfig.TLSKey,
		},
		cli.StringFlag{
			Name:        FlagTLSCert,
			Usage:       "The tls cert file to be used",
			Destination: &clientConfig.TLSCert,
		},
		cli.StringFlag{
			Name:        FlagTLSClientCA,
			Usage:       "The tls client ca file to be used",
			Destination: &clientConfig.TLSClientCA,
		},
	}

	clientActions = []cli.Command{
		{
			Name:   "login",
			Usage:  "issues a new bearer token from the authproxy",
			Action: loginAction,
		},
		{
			Name:   "authenticate",
			Usage:  "authenticates against the auth proxy",
			Action: authAction,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "example cli client for authproxy interactions"
	app.Commands = clientActions
	app.Flags = clientFlags

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("failed to run cli: %s", err)
		os.Exit(1)
	}
}

func loginAction(c *cli.Context) error {

	if len(c.Args()) == 0 {
		fmt.Println("please enter username as password")
	}

	username, password := c.Args()[0], c.Args()[1]

	cfg := client.AuthClientConfig{
		TLSKey:      clientConfig.TLSKey,
		TLSCert:     clientConfig.TLSCert,
		TLSClientCA: clientConfig.TLSClientCA,
	}

	cl, err := client.NewForConfig(&cfg)
	if err != nil {
		return err
	}

	token, err := cl.Login(username, password)
	if err != nil {
		return err
	}

	fmt.Println("Received token for user: " + token)
	return nil
}

func authAction(c *cli.Context) error {

	if len(c.Args()) == 0 {
		fmt.Println("please enter a bearerToken")
	}

	token := c.Args()[0]

	cfg := client.AuthClientConfig{
		TLSKey:      clientConfig.TLSKey,
		TLSCert:     clientConfig.TLSCert,
		TLSClientCA: clientConfig.TLSClientCA,
	}

	cl, err := client.NewForConfig(&cfg)
	if err != nil {
		return err
	}

	rr, err := cl.Authenticate(token)
	if err != nil {
		return err
	}

	if !rr.Status.Authenticated {
		fmt.Println("client unauthenticated, token is invalid")
		return nil
	}

	fmt.Println("client successfully authenticated, token is valid")
	return nil
}

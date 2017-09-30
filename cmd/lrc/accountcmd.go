/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package main

import (
	"fmt"
	ethClient "github.com/Loopring/ringminer/chainclient/eth"
	"github.com/Loopring/ringminer/crypto"
	"github.com/Loopring/ringminer/types"
	"gopkg.in/urfave/cli.v1"
)

func accountCommands() cli.Command {
	c := cli.Command{
		Name:     "account",
		Usage:    "generate, encrypt and decrypt an account",
		Category: "account commands:",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "generate",
				Usage:  "generate a new account",
				Action: generatePrivateKey,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "display",
						Usage: "display the privatekey",
					},
					cli.StringFlag{
						Name:  "passphrase,p",
						Usage: "passphrase for encrypted the private",
					},
					cli.StringFlag{
						Name:  "private-key,pk",
						Usage: "generate account from this private-key when set it",
					},
				},
			},
			cli.Command{
				Name:   "encrypt",
				Usage:  "encrypt a private key using the passphrase",
				Action: encrypt,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "passphrase,p",
						Usage: "passphrase for encrypted the private",
					},
					cli.StringFlag{
						Name:  "private-key,pk",
						Usage: "the private key to be encrypted",
					},
				},
			},
			cli.Command{
				Name:   "decrypt",
				Usage:  "decrypt a encrepted private key using the passphrase",
				Action: decrypt,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "passphrase,p",
						Usage: "passphrase for encrypted the private",
					},
					cli.StringFlag{
						Name:  "encrypted,e",
						Usage: "the encrypted private key",
					},
				},
			},
		},
	}
	return c
}

func encrypt(c *cli.Context) {
	pk := c.String("private-key")
	if !types.IsHex(pk) {
		panic("the private-key must be hex")
	}
	passphrase := c.String("passphrase")
	p := &types.Passphrase{}
	p.SetBytes([]byte(passphrase))

	if encrypted, err := crypto.AesEncrypted(p.Bytes(), types.FromHex(pk)); nil != err {
		fmt.Fprintf(c.App.Writer, "%v \n", err.Error())
	} else {
		fmt.Fprintf(c.App.Writer, "encrypted private key:%v \n", types.ToHex(encrypted))
	}
}

func decrypt(c *cli.Context) {
	encrypted := c.String("encrypted")
	if !types.IsHex(encrypted) {
		panic("the encrypted must be hex")
	}
	passphrase := c.String("passphrase")
	p := &types.Passphrase{}
	p.SetBytes([]byte(passphrase))

	if pk, err := crypto.AesDecrypted(p.Bytes(), types.FromHex(encrypted)); nil != err {
		fmt.Fprintf(c.App.Writer, "%v \n", err.Error())
	} else {
		fmt.Fprintf(c.App.Writer, "private key:%v \n", types.ToHex(pk))
	}
}

func generatePrivateKey(c *cli.Context) {
	passphrase := c.String("passphrase")
	diaplay := c.Bool("display")
	pk := c.String("private-key")
	p := &types.Passphrase{}
	p.SetBytes([]byte(passphrase))

	generateEthPrivateKey(pk, p, diaplay, c)
}

func generateEthPrivateKey(pk string, passphrase *types.Passphrase, display bool, c *cli.Context) {
	if account, err := ethClient.NewAccount(pk); nil != err {
		fmt.Fprintf(c.App.Writer, "%v \n", err.Error())
	} else {
		if _, err := account.Encrypted(passphrase); nil != err {
			fmt.Fprintf(c.App.Writer, "%v \n", err.Error())
		} else {
			fmt.Fprintf(c.App.Writer, "address:%v encrypted private key:%v \n", account.Address.Hex(), types.ToHex(account.EncryptedPrivKey))
			if display {
				fmt.Fprintf(c.App.Writer, "private key:%v \n", types.ToHex(account.PrivKey.D.Bytes()))
			}
		}
	}
}

// Copyright 2023 The Tongsuo Project Authors. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://github.com/jeffrey1205/tongsuo-go-sdk/blob/main/LICENSE

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	ts "github.com/jeffrey1205/tongsuo-go-sdk"
	"github.com/jeffrey1205/tongsuo-go-sdk/crypto"
)

func main() {
	cipherSuite := ""
	signCertFile := ""
	signKeyFile := ""
	encCertFile := ""
	encKeyFile := ""
	caFile := ""
	connAddr := ""

	flag.StringVar(&connAddr, "conn", "127.0.0.1:443", "host:port")
	flag.StringVar(&cipherSuite, "cipher", "ECC-SM2-SM4-CBC-SM3", "cipher suite")
	flag.StringVar(&signCertFile, "sign_cert", "", "sign certificate file")
	flag.StringVar(&signKeyFile, "sign_key", "", "sign private key file")
	flag.StringVar(&encCertFile, "enc_cert", "", "encrypt certificate file")
	flag.StringVar(&encKeyFile, "enc_key", "", "encrypt private key file")
	flag.StringVar(&caFile, "CAfile", "", "CA certificate file")

	flag.Parse()

	ctx, err := ts.NewCtxWithVersion(ts.NTLS)
	if err != nil {
		panic(err)
	}

	if err := ctx.SetCipherList(cipherSuite); err != nil {
		panic(err)
	}

	if signCertFile != "" {
		signCertPEM, err := os.ReadFile(signCertFile)
		if err != nil {
			panic(err)
		}
		signCert, err := crypto.LoadCertificateFromPEM(signCertPEM)
		if err != nil {
			panic(err)
		}

		if err := ctx.UseSignCertificate(signCert); err != nil {
			panic(err)
		}
	}

	if signKeyFile != "" {
		signKeyPEM, err := os.ReadFile(signKeyFile)
		if err != nil {
			panic(err)
		}
		signKey, err := crypto.LoadPrivateKeyFromPEM(signKeyPEM)
		if err != nil {
			panic(err)
		}

		if err := ctx.UseSignPrivateKey(signKey); err != nil {
			panic(err)
		}
	}

	if encCertFile != "" {
		encCertPEM, err := os.ReadFile(encCertFile)
		if err != nil {
			panic(err)
		}
		encCert, err := crypto.LoadCertificateFromPEM(encCertPEM)
		if err != nil {
			panic(err)
		}

		if err := ctx.UseEncryptCertificate(encCert); err != nil {
			panic(err)
		}
	}

	if encKeyFile != "" {
		encKeyPEM, err := os.ReadFile(encKeyFile)
		if err != nil {
			panic(err)
		}

		encKey, err := crypto.LoadPrivateKeyFromPEM(encKeyPEM)
		if err != nil {
			panic(err)
		}

		if err := ctx.UseEncryptPrivateKey(encKey); err != nil {
			panic(err)
		}
	}

	if caFile != "" {
		if err := ctx.LoadVerifyLocations(caFile, ""); err != nil {
			panic(err)
		}
	}

	conn, err := ts.Dial("tcp", connAddr, ctx, ts.InsecureSkipHostVerification)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cipher, err := conn.CurrentCipher()
	if err != nil {
		panic(err)
	}

	ver, err := conn.GetVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println("New connection: " + ver + ", cipher=" + cipher)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	request := text + "\n"
	fmt.Println(">>>\n" + request)
	if _, err := conn.Write([]byte(request)); err != nil {
		panic(err)
	}

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}

	fmt.Println("<<<\n" + string(buffer[:n]))

	return
}

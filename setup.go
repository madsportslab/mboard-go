package main

import (
	"log"
	"net"
	"net/http"

	"github.com/eknkc/amber"
	"github.com/skip2/go-qrcode"
)

func getAddress(name string) string {

  ifs, err := net.Interfaces()

	if err != nil {
		
		log.Println(err)
	  return ""

	}

	for _, iface := range ifs {

    if iface.Name == name {
			
			addrs, err := iface.Addrs()

			if err != nil {
				log.Println(err)
				break
			}

			ipnet, ok := addrs[1].(*net.IPNet)

			if ok {
				return ipnet.IP.String()
			}

		}

	}

	return "hello"

} // getAddress

func setupHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:

		log.Println("generate QR code")
		
		ip := getAddress("en0")

		err2 := qrcode.WriteFile(ip, qrcode.Medium, 512, "www/qr.png")

		if err2 != nil {
			log.Println(err2)
		}

		compiler := amber.New()

	  parseErr := compiler.ParseFile("www/setup.amber")

		if parseErr != nil {
			
			log.Printf("[%s][Error] %s", version(), parseErr)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template, compileErr := compiler.Compile()

		if compileErr != nil {
			
			log.Printf("[%s][Error] %s", version(), compileErr)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		template.Execute(w, nil)

		
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // setupHandler

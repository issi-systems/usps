package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

// demonstrate usps api
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("Username")
	b, err := validate(&Address{
		User:     user,
		Address1: "5 Jonathan Morris Circle",
		Address2: "",
		City:     "Media",
		State:    "Pa",
		Zip:      "19063",
		Zip_ext:  "",
	})
	if err != nil {
		log.Printf("Error %v", err)
	} else {
		log.Printf("%s", b)
	}
}

var tmp = `https://secure.shippingapis.com/ShippingAPI.dll?API=Verify&XML=`
var xml = `<AddressValidateRequest USERID="{{ .User }}"><Revision>1</Revision><Address ID="0"><Address1>{{ .Address1  }}</Address1><Address2>{{ .Address2 }}</Address2><City>{{ .City }}</City><State>{{ .State }}</State><Zip5>{{ .Zip }}</Zip5><Zip4>{{ .Zip_ext }}</Zip4></Address></AddressValidateRequest>`

type Address struct {
	User     string
	Address1 string
	Address2 string
	City     string
	State    string
	Zip      string
	Zip_ext  string
}

func validate(a *Address) ([]byte, error) {
	t, err := template.New("todos").Parse(
		xml)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, a)
	if err != nil {
		return nil, err
	}

	//log.Printf("%s\n", buf.String())
	res, err := http.Get(tmp + url.QueryEscape(buf.String()))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<AddressValidateResponse><Address ID="0"><Address2>5 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4>1069</Zip4><DeliveryPoint>05</DeliveryPoint><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>Y</DPVConfirmation><DPVCMRA>N</DPVCMRA><DPVFootnotes>AABB</DPVFootnotes><Business>N</Business><CentralDeliveryPoint>N</CentralDeliveryPoint><Vacant>N</Vacant></Address></AddressValidateResponse>

<?xml version="1.0" encoding="UTF-8"?>
<AddressValidateResponse><Address ID="0"><Address2>99 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4/><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>N</DPVConfirmation><DPVFalse>N</DPVFalse><DPVFootnotes>AAM3</DPVFootnotes><CentralDeliveryPoint>N</CentralDeliveryPoint></Address></AddressValidateResponse>

*/

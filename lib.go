package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

// //https://www.onlinetool.io/xmltogo/
var ApiKey string

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ApiKey = os.Getenv("Username")
	if len(ApiKey) == 0 {
		panic("api key is required")
	}
}

var tmp = `https://secure.shippingapis.com/ShippingAPI.dll?API=Verify&XML=`
var xmlv = `<AddressValidateRequest USERID="{{ .User }}"><Revision>1</Revision><Address ID="0"><Address1>{{ .Address1  }}</Address1><Address2>{{ .Address2 }}</Address2><City>{{ .City }}</City><State>{{ .State }}</State><Zip5>{{ .Zip }}</Zip5><Zip4>{{ .Zip_ext }}</Zip4></Address></AddressValidateRequest>`

type AddressValidateResponse struct {
	XMLName xml.Name `xml:"AddressValidateResponse"`
	Text    string   `xml:",chardata"`
	Address struct {
		Text                 string `xml:",chardata"`
		ID                   string `xml:"ID,attr"`
		Address1             string `xml:"Address1"`
		Address2             string `xml:"Address2"`
		City                 string `xml:"City"`
		State                string `xml:"State"`
		Zip5                 string `xml:"Zip5"`
		Zip4                 string `xml:"Zip4"`
		DeliveryPoint        string `xml:"DeliveryPoint"`
		CarrierRoute         string `xml:"CarrierRoute"`
		Footnotes            string `xml:"Footnotes"`
		DPVConfirmation      string `xml:"DPVConfirmation"`
		DPVCMRA              string `xml:"DPVCMRA"`
		DPVFootnotes         string `xml:"DPVFootnotes"`
		Business             string `xml:"Business"`
		CentralDeliveryPoint string `xml:"CentralDeliveryPoint"`
		Vacant               string `xml:"Vacant"`
	} `xml:"Address"`
}

func (ax *AddressValidateResponse) FromXml(data []byte) error {
	return xml.Unmarshal(data, ax)
}
func (ax *AddressValidateResponse) ToJson() ([]byte, error) {
	a := &ax.Address
	return json.Marshal(&Address{
		Address1:        a.Address2,
		Address2:        a.Address1,
		City:            a.City,
		State:           a.State,
		Zip:             a.Zip5,
		Zip_ext:         a.Zip4,
		DPVConfirmation: a.DPVConfirmation,
	})
}

type Address struct {
	User            string
	Address1        string `json:"Address1,omitempty"`
	Address2        string `json:"Address2,omitempty"`
	City            string `json:"City,omitempty"`
	State           string `json:"State,omitempty"`
	Zip             string `json:"Zip5,omitempty"`
	Zip_ext         string `json:"Zip4,omitempty"`
	DPVConfirmation string `xml:"DPVConfirmation"`
}

func (a *Address) FromJson(data []byte) error {
	return json.Unmarshal(data, a)
}
func (a *Address) Validate() ([]byte, error) {
	t, err := template.New("todos").Parse(
		xmlv)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, a)
	if err != nil {
		return nil, err
	}
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
func ValidateJson(s []byte) ([]byte, error) {
	var ad Address
	var ar AddressValidateResponse
	e := ad.FromJson(s)
	if e != nil {
		return nil, e
	}

	r1, e := ad.Validate()
	if e != nil {
		return nil, e
	}

	e = ar.FromXml(r1)
	if e != nil {
		return nil, e
	}

	return ar.ToJson()
}

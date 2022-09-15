package usps

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
var xmlv = `<AddressValidateRequest USERID="{{ .User }}"><Revision>1</Revision>
{{ range .Address }}
<Address ID="{{.ID}}"><Address1>{{ .Address1  }}</Address1><Address2>{{ .Address2 }}</Address2><City>{{ .City }}</City><State>{{ .State }}</State><Zip5>{{ .Zip }}</Zip5><Zip4>{{ .Zip_ext }}</Zip4></Address>
{{end}}
</AddressValidateRequest>`

type AddressValidateResponse struct {
	XMLName xml.Name `xml:"AddressValidateResponse"`
	Text    string   `xml:",chardata"`
	Address []struct {
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
	al := []*Address{}
	for _, a := range ax.Address {
		al = append(al, &Address{
			ID:              a.ID,
			Address1:        a.Address2,
			Address2:        a.Address1,
			City:            a.City,
			State:           a.State,
			Zip:             a.Zip5,
			Zip_ext:         a.Zip4,
			DPVConfirmation: a.DPVConfirmation,
		})
	}
	return json.Marshal(al)
}

type Address struct {
	ID              string `json:"ID,omitempty"`
	Address1        string `json:"Address1,omitempty"`
	Address2        string `json:"Address2,omitempty"`
	City            string `json:"City,omitempty"`
	State           string `json:"State,omitempty"`
	Zip             string `json:"Zip5,omitempty"`
	Zip_ext         string `json:"Zip4,omitempty"`
	DPVConfirmation string `json:"DPVConfirmation,omitempty"`
}
type AddressList struct {
	User    string
	Address []*Address `json:"Address,omitempty"`
}

func (a *AddressList) ToXml() (string, error) {
	a.User = ApiKey
	t, err := template.New("todos").Parse(
		xmlv)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, a)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a *AddressList) Validate() ([]byte, error) {
	x, e := a.ToXml()
	if e != nil {
		return nil, e
	}
	return a.ValidateXml(x)
}
func (a *AddressList) ValidateXml(x string) ([]byte, error) {

	res, err := http.Get(tmp + url.QueryEscape(x))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("query: %s\n response: %s\n", x, string(b))
	return b, nil
}

func ValidateJson(b []byte) ([]byte, error) {
	var address = &AddressList{}
	e := json.Unmarshal(b, &address.Address)
	if e != nil || len(address.Address) == 0 {
		log.Printf("0 %v,%s", e, string(b))
		return nil, e
	}
	x, e := address.ToXml()
	if e != nil {
		log.Printf("1b %v", e)
	}
	r1, e := address.ValidateXml(string(x))
	if e != nil {
		log.Printf("2 %v,%s", e, x)
		return nil, e
	}

	var ar = &AddressValidateResponse{}
	e = ar.FromXml(r1)
	if e != nil {
		log.Printf("3 %v,%v,%s", e, r1, x)
		return nil, e
	}
	tj, e := ar.ToJson()
	if e != nil {
		log.Printf("4 %v", e)
		return nil, e
	}
	return tj, nil

}

/*

	var ad AddressList
	var ar AddressValidateResponse
	e := json.Unmarshal(s,&ad)
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
*/

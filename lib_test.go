package usps

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var js1 = `[{"ID": "PRT", "Address1":"5 Jonathan morris circle","Address2":"","City":"Media","State":"PA","Zip5":"19063","Zip4":""},{"ID": "PRV", "Address1":"5 Jonathan morris","Address2":"","City":"Broomall","State":"PA","Zip5":"19063","Zip4":""}]`

func Test_toxml(t *testing.T) {
	Init()
	al := &AddressList{
		User:    ApiKey,
		Address: []*Address{},
	}
	spew.Dump(json.Unmarshal([]byte(js1), &al.Address))
	x, e := al.ToXml()
	log.Printf("xml %v, %s", e, x)
}

func Test_all(t *testing.T) {
	Init()
	o, e := ValidateJson([]byte(js1))
	log.Printf("xml %v, %s", e, o)
}

func Test_two(t *testing.T) {
	result := `<AddressValidateResponse><Address ID="0"><Address2>5 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4>1069</Zip4><DeliveryPoint>05</DeliveryPoint><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>Y</DPVConfirmation><DPVCMRA>N</DPVCMRA><DPVFootnotes>AABB</DPVFootnotes><Business>N</Business><CentralDeliveryPoint>N</CentralDeliveryPoint><Vacant>N</Vacant></Address></AddressValidateResponse>`
	var a AddressValidateResponse
	e := a.FromXml([]byte(result))
	if e != nil {
		panic(e)
	}
	spew.Dump(a)
	spew.Dump(a.ToJson())
}

func Test_validate(t *testing.T) {
	Init()
	r1, e1 := sample.Validate()
	var ar AddressValidateResponse
	e2 := ar.FromXml(r1)
	tj, e3 := ar.ToJson()
	spew.Dump(ar, e1, e2)
	spew.Dump(tj, e3)
}

var sample = &AddressList{
	User: "",
	Address: []*Address{
		{
			Address1: "5 Jonathan Morris Circle",
			Address2: "",
			City:     "Media",
			State:    "Pa",
			Zip:      "19063",
			Zip_ext:  "",
		},
		{
			Address1: "109 Columbia",
			Address2: "",
			City:     "Newtown Square",
			State:    "Pa",
			Zip:      "19073",
			Zip_ext:  "",
		},
	},
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<AddressValidateResponse><Address ID="0"><Address2>5 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4>1069</Zip4><DeliveryPoint>05</DeliveryPoint><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>Y</DPVConfirmation><DPVCMRA>N</DPVCMRA><DPVFootnotes>AABB</DPVFootnotes><Business>N</Business><CentralDeliveryPoint>N</CentralDeliveryPoint><Vacant>N</Vacant></Address></AddressValidateResponse>

<?xml version="1.0" encoding="UTF-8"?>
<AddressValidateResponse><Address ID="0"><Address2>99 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4/><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>N</DPVConfirmation><DPVFalse>N</DPVFalse><DPVFootnotes>AAM3</DPVFootnotes><CentralDeliveryPoint>N</CentralDeliveryPoint></Address></AddressValidateResponse>

*/

package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_main(t *testing.T) {
	Init()
	s := sample()
	spew.Dump(ValidateJson(s))
}
func Test_one(t *testing.T) {

	log.Printf(ApiKey)

	result := `<AddressValidateResponse><Address ID="0"><Address2>5 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4>1069</Zip4><DeliveryPoint>05</DeliveryPoint><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>Y</DPVConfirmation><DPVCMRA>N</DPVCMRA><DPVFootnotes>AABB</DPVFootnotes><Business>N</Business><CentralDeliveryPoint>N</CentralDeliveryPoint><Vacant>N</Vacant></Address></AddressValidateResponse>`
	var a AddressValidateResponse
	e := a.FromXml([]byte(result))
	if e != nil {
		panic(e)
	}
	spew.Dump(a)
	spew.Dump(a.ToJson())

	var ad Address
	var ar AddressValidateResponse

	s := sample()
	e = ad.FromJson(s)
	r1, e1 := ad.Validate()
	e2 := ar.FromXml(r1)
	tj, e3 := ar.ToJson()
	spew.Dump(ar, e, e1, e2)
	spew.Dump(tj, e3)
}

func sample() []byte {
	address := &Address{
		User:     ApiKey,
		Address1: "5 Jonathan Morris Circle",
		Address2: "",
		City:     "Media",
		State:    "Pa",
		Zip:      "19063",
		Zip_ext:  "",
	}

	bf, _ := json.Marshal(address)
	return bf
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<AddressValidateResponse><Address ID="0"><Address2>5 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4>1069</Zip4><DeliveryPoint>05</DeliveryPoint><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>Y</DPVConfirmation><DPVCMRA>N</DPVCMRA><DPVFootnotes>AABB</DPVFootnotes><Business>N</Business><CentralDeliveryPoint>N</CentralDeliveryPoint><Vacant>N</Vacant></Address></AddressValidateResponse>

<?xml version="1.0" encoding="UTF-8"?>
<AddressValidateResponse><Address ID="0"><Address2>99 JONATHAN MORRIS CIR</Address2><City>MEDIA</City><State>PA</State><Zip5>19063</Zip5><Zip4/><CarrierRoute>C069</CarrierRoute><Footnotes>N</Footnotes><DPVConfirmation>N</DPVConfirmation><DPVFalse>N</DPVFalse><DPVFootnotes>AAM3</DPVFootnotes><CentralDeliveryPoint>N</CentralDeliveryPoint></Address></AddressValidateResponse>

*/

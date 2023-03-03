package main

import (
	"Protocol_Buffers/pb"
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

func main() {
	/*
		-------------------------------------------------------
		* データのシリアライズ
		-------------------------------------------------------
	*/
	employee := &pb.Employee{
		Id:         1,
		Name:       "suzuki",
		Email:      "test@test.com",
		Occupation: pb.Occupation_ENGINEER,
		PhoneNumber: []string{
			"080-1234-5578",
			"090-1212-4343",
		},
		Project: map[string]*pb.Company_Project{
			"ProjectX": &pb.Company_Project{},
		},
		Profile: &pb.Employee_Text{
			Text: "my name is suzuki",
		},
		Birthday: &pb.Date{
			Year:  2000,
			Month: 1,
			Day:   2,
		},
	}

	binData, err := proto.Marshal(employee)
	if err != nil {
		log.Fatalln("Cant serialize", err)
	}

	if err := os.WriteFile("test.bin", binData, 0666); err != nil {
		log.Fatalln("Cant write", err)
	}

	// deserialize
	in, err := os.ReadFile("test.bin")
	if err != nil {
		log.Fatalln("Cant read file", err)
	}

	readEmployee := &pb.Employee{}
	err = proto.Unmarshal(in, readEmployee)
	if err != nil {
		log.Fatalln("Cant deserialize", err)
	}

	fmt.Println(readEmployee)

	/*
		-------------------------------------------------------
		* JSONマッピング
		-------------------------------------------------------
	*/

	m := jsonpb.Marshaler{}
	out, err := m.MarshalToString(employee)
	if err != nil {
		log.Fatalln("Cant marshal to JSON", err)
	}
	fmt.Println(out)

	readEmployee2 := &pb.Employee{}
	if err := jsonpb.UnmarshalString(out, readEmployee2); err != nil {
		log.Fatalln("Cant unmarshal from JSON", err)
	}

	fmt.Println(readEmployee2)
}

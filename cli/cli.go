package main

import (
	"awesomeProject/shipment/rpc"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func ParseFile(file string) (*rpc.Consignment, error) {
	var consignment *rpc.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, nil
}

func main() {
	con, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	client := rpc.NewShippingServiceClient(con)

	file := defaultFilename

	consignmet, err := ParseFile(file)

	if err != nil {
		log.Fatal(err)
	}

	r, err := client.CreateConsignment(context.Background(), consignmet)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(r.Created)

	getAll, err := client.GetConsignments(context.Background(), &rpc.GetRequest{})

	if err != nil {
		log.Fatal(err)
	}

	for _, cns := range getAll.Consignments {
		fmt.Printf("Id: %v\n", cns.GetId())
		fmt.Printf("Description: %v\n", cns.GetDescription())
		fmt.Printf("Weight: %d\n", cns.GetWeight())
		fmt.Printf("VesselId: %v\n", cns.GetVesselId())
		for _, cnt := range cns.GetContainers() {
			fmt.Printf("\tId: %v\n", cnt.GetId())
			fmt.Printf("\tUserId: %v\n", cnt.GetUserId())
			fmt.Printf("\tCustomerId: %v\n", cnt.GetCustomerId())
			fmt.Printf("\tOrigin: %v\n", cnt.GetOrigin())
		}
	}
}

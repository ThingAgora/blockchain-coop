// Copyright Luc Yriarte <luc.yriarte@thingagora.org> 2018 
// License: Apache-2.0

package main

import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SSMChaincode struct {
}


func (t *SSMChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect arg count. Expecting 1")
	}

	var err error	
	var admins []Agent
	err = json.Unmarshal([]byte(args[0]), &admins)
	if (err != nil) {
		return shim.Error(err.Error())
	}
	for i := 0; i < len(admins); i++ {
		admin_info, err := json.Marshal(admins[i])
		if (err != nil) {
			return shim.Error(err.Error())
		}
		err = stub.PutState("ADMIN_" + admins[i].Name, admin_info)
		if (err != nil) {
			return shim.Error(err.Error())
		}
	}
	return shim.Success(nil)
}

func (t *SSMChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	
	var err error	
	var admin Agent
	var user Agent
	if function == "register" {
		if len(args) != 3 {
			return shim.Error("Incorrect arg count.")
		}
		err = json.Unmarshal([]byte(args[0]), &user)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		dat, err := stub.GetState("ADMIN_" + args[1]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(dat, &admin)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		// TODO validate admin signature
		err = stub.PutState("USER_" + user.Name, []byte(args[0]))
		if (err != nil) {
			return shim.Error(err.Error())
		}
	}
	
	if function == "create" {
		if len(args) != 3 {
			return shim.Error("Incorrect arg count.")
		}
		var ssm SigningStateMachine
		err = json.Unmarshal([]byte(args[0]), &ssm)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		dat, err := stub.GetState("ADMIN_" + args[1]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(dat, &admin)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		// TODO validate admin signature
		err = stub.PutState("SSM_" + ssm.Name, []byte(args[0]))
		if (err != nil) {
			return shim.Error(err.Error())
		}
	}
	
	if function == "start" {
		if len(args) != 3 {
			return shim.Error("Incorrect arg count.")
		}
		var state State
		err = json.Unmarshal([]byte(args[0]), &state)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		dat, err := stub.GetState("ADMIN_" + args[1]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(dat, &admin)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		// TODO validate admin signature
		err = stub.PutState("STATE_" + state.Session, []byte(args[0]))
		if (err != nil) {
			return shim.Error(err.Error())
		}
	}
	
	if function == "perform" {
		if len(args) != 4 {
			return shim.Error("Incorrect arg count.")
		}
		var state State
		dat, err := stub.GetState("STATE_" + args[1]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(dat, &state)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		dat, err = stub.GetState("USER_" + args[2]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(dat, &user)
		if (err != nil) {
			return shim.Error(err.Error())
		}
		// TODO validate user signature
	}
	
	if function == "session" {
		if len(args) != 1 {
			return shim.Error("Incorrect arg count.")
		}
		dat, err := stub.GetState("STATE_" + args[0]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		return shim.Success(dat)
	}
	
	if function == "ssm" {
		if len(args) != 1 {
			return shim.Error("Incorrect arg count.")
		}
		dat, err := stub.GetState("SSM_" + args[0]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		return shim.Success(dat)
	}
	
	if function == "user" {
		if len(args) != 1 {
			return shim.Error("Incorrect arg count.")
		}
		dat, err := stub.GetState("USER_" + args[0]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		return shim.Success(dat)
	}
	
	if function == "admin" {
		if len(args) != 1 {
			return shim.Error("Incorrect arg count.")
		}
		dat, err := stub.GetState("ADMIN_" + args[0]);
		if (err != nil) {
			return shim.Error(err.Error())
		}
		return shim.Success(dat)
	}
	
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SSMChaincode))
	if err != nil {
		fmt.Println("Signing State Machines Chaincode init error:", err)
	}
}

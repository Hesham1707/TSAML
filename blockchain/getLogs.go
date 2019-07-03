
package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"encoding/json"
	"fmt"
)

func (setup *FabricSetup)GetLogs(clientID string) (string, error){
	

	
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")


	fname := "queryClientLogs"
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{[]byte(clientID)}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	c := client{}
	fmt.Println("clientID : ",clientID)
	err =json.Unmarshal(response.Payload,&c)
	if err != nil{
		fmt.Println("can't Unmarshal logs")
		return "", err
	}
	fmt.Println(c.Name)
	logs :=c.UserLogs
	i:=0
	for i=0;i<len(logs);i++{
		fmt.Println(logs[i].Type)
	}

	logstring:= setlogstring(logs,c.Name)

	return logstring,nil
}

func setlogstring(logs []log, name string) (string){
	logstring:="Showing logs for user: " + name +" \n "
	i:=0
	for i=0;i<len(logs);i++{
		if logs[i].Type != "transfer"{
			if i != len(logs)-1 {
				logstring += logs[i].toString()+" \n "
			}else {
				logstring += logs[i].toString()
			}
		}else{
			if i != len(logs)-1 {
				logstring += logs[i].toStringTransfer()+" \n "
			}else {
				logstring += logs[i].toStringTransfer()
			}
		}
	}
	return logstring
}

func (setup *FabricSetup)GetClientNames() (string, error){
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")


	fname := "queryuserList"
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	list := names{}
	err =json.Unmarshal(response.Payload,&list)
	if err != nil{
		fmt.Println("can't Unmarshal logs")
		return "", err
	}
	ClintInfo :=list.Usernames
	i:=0
	for i=0;i<len(ClintInfo);i++{
		fmt.Println(ClintInfo[i])
	}

	logstring:=""

	for i=0;i<len(ClintInfo);i++{
		logstring += ClintInfo[i]+" \n "
		logstring +="---------------\n"
	}
	return logstring,nil
}
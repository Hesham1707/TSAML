package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// InvokeHello
func (setup *FabricSetup) InvokeHello(fname string , args []string) (string, error) {

	// Prepare arguments
	//var args []string
	//args = append(args, "invoke")
	//args = append(args, "invoke")
	//args = append(args, "hello")
	//args = append(args, value)

	//eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	/*reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)*/


	response := channel.Response{}
	var err error

	switch len(args) {
		case 1: 
		   	// Create a request (proposal) and send it
			response, err = setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{[]byte(args[0])}, TransientMap: transientDataMap})
			if err != nil {
			return "", fmt.Errorf("failed %v", fname)
			}
		case 2:
		   	// Create a request (proposal) and send it
			   response, err = setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{[]byte(args[0]), []byte(args[1])}, TransientMap: transientDataMap})
			   if err != nil {
			   return "", fmt.Errorf("failed %v", fname)
			   }
		case 3:
			// Create a request (proposal) and send it
			response, err = setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2])}, TransientMap: transientDataMap})
			if err != nil {
			return "", fmt.Errorf("failed %v", fname)
			}
		case 4:
			// Create a request (proposal) and send it
			response, err = setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
			if err != nil {
			return "", fmt.Errorf("failed %v", fname)
			}
		case 5:
			// Create a request (proposal) and send it
			response, err = setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])}, TransientMap: transientDataMap})
			if err != nil {
			return "", fmt.Errorf("failed %v", fname)
			}
		}

	// Wait for the result of the submission
	/*select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}*/

	return string(response.TransactionID), nil
}
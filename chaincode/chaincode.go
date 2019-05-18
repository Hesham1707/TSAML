package main


import (
    "fmt"
    "encoding/json"
    "strconv"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {

}
// Asset client 
type client struct{
	Name string  `json:"name"`
	Id string 	`json:"id"`
    Balance float64 `json:"balance"`
}
// Asset of stocks
type stock struct{
    Name string `json:"name"`
    Number int `json:"number"`
    Price float64 `json:"price"`
    Owner string `json:"owner"`
    Offeredstocks int `json:"offeredstocks"`
}

// Init is called during chaincode instantiation to initialize any data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
    // Put in the ledger the key/value hello/world
    err := stub.PutState("hello", []byte("world"))
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The 'set'
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// Extract the function and args from the transaction proposal
    fn, args := stub.GetFunctionAndParameters()
    switch fn {
		case "transfer": 
            return t.transfer(stub, args)

		case "withdraw": 
            return t.withdraw(stub, args)

		case "deposit": 
            return t.deposit(stub, args)

		case "getBalance": 
            return t.getBalance(stub, args)

        case "createClient": 
            return t.createClient(stub, args)

        case "addStocks": 
            return t.AddStocks(stub, args)

        case "removeStocks": 
            return t.RemoveStocks(stub, args)

        case "buyStocks": 
            return t.buyStocks(stub, args)

        case "sellStocks": 
            return t.sellStocks(stub, args)
            
	}

	//return shim.Error("Unknown action, check the first argument")
	return shim.Error(fn)
}

//args[0] => nameOfstock, args[1] => numberOfstock, args[2] => ownerOfstock, args[3] => priceOfstock
func (t *SimpleAsset) AddStocks(stub shim.ChaincodeStubInterface, args []string) peer.Response{
    if len(args) != 5 {
        return shim.Error("Incorrect arguments. Expecting a key and a value")
    }
    
    sname := args[0]
    snumber,err := strconv.Atoi(args[1])
    sowner := args[2]
    sprice,err := strconv.ParseFloat(args[3],64)
    //check if price not float or number not integer
    if err != nil {
        return shim.Error(err.Error())
    }
    sofferedstocks := snumber
    //create compsite key (nameOfStock~ownerid) and pass it to stub
    indexName := "nameOfStock~ownerid"
	stockKey, err := stub.CreateCompositeKey(indexName, []string{sname, sowner})
    if err != nil {
        return shim.Error(err.Error())
    }
    //make object of struct stock and translate it to bytes
    stock1 := &stock{sname,snumber,sprice,sowner,sofferedstocks}
	fmt.Println("stock name = %s" , stock1.Name)
	fmt.Println(" Number of stcks = %s" , stock1.Number)
	fmt.Println("stock owner = " , stock1.Owner)
    fmt.Println("price of stock = " , stock1.Price)
    fmt.Println("offered stocks = " , stock1.Offeredstocks)
	stockBytes,err := json.Marshal(stock1) 
	if err != nil {
    	return shim.Error(err.Error())
	}

    err = stub.PutState(stockKey, stockBytes)
    if err != nil {
            return shim.Error("Failed to add stock")
    }
    log:=fmt.Sprintf("Sucessfully added stock with owner %s :) ", sowner)
	fmt.Println(log)
	
	return shim.Success([]byte(log))
}

//args[0] => nameOfstock, args[1] =>owner
func (t *SimpleAsset) RemoveStocks(stub shim.ChaincodeStubInterface, args []string) peer.Response{
    if len(args) != 2 {
        return shim.Error("Incorrect arguments. Expecting owner and name of stock")
    }
    sname := args[0]
    sowner := args[1]
    //create compsite key (nameOfStock~ownerid) and pass it to stub
    indexName := "nameOfStock~ownerid"
	stockKey, err := stub.CreateCompositeKey(indexName, []string{sname, sowner})
    if err != nil {
        return shim.Error(err.Error())
    }
    //remove stock
    err = stub.PutState(stockKey, nil)
    if err != nil {
        return shim.Error("Failed to remove stock")
    }

    log:=fmt.Sprintf("Sucessfully removed stock with owner %s :) ", sowner)
    fmt.Println(log)
    return shim.Success([]byte(log))
}
//args[0] => nameOfstock, args[1] => owner args[2] => priceOfstock, args[3] => offeredstocks
func (t *SimpleAsset) sellStocks(stub shim.ChaincodeStubInterface, args []string) peer.Response{
    if len(args) != 4 {
        return shim.Error("Incorrect arguments. Expecting no., name of stock and price")
    }
    sname :=  args[0]
    sowner := args[1]
    sprice,err := strconv.ParseFloat(args[2],64)
    sofferedstocks,err := strconv.Atoi(args[3])
    //check if price not float or number not integer
    if err != nil {
        return shim.Error(err.Error())
    }
    //create compsite key (nameOfStock~ownerid) and pass it to stub
    indexName := "nameOfStock~ownerid"
    stockKey, err := stub.CreateCompositeKey(indexName, []string{sname, sowner})
    if err != nil {
         return shim.Error(err.Error())
    }
    //check if client has this stocks
    exist,stockBytes := checkStocksExist(stub,stockKey)
    if exist!=true{
        fmt.Println("client don't have this stocks")
        return shim.Error("client don't have this stocks")
    }
    
    s := stock{}
    err=unMarshalStock(stockBytes,&s)
    if err != nil {
    	return shim.Error(err.Error())
    }
    //check if client have enough stocks to offer
    if sofferedstocks > s.Number{
        fmt.Println("client don't have enough stocks")
        return shim.Error("client don't have enough stocks")
    }
    s.Price=sprice
    s.Offeredstocks=sofferedstocks

    stockBytes,err = json.Marshal(s) 
	if err != nil {
    	return shim.Error(err.Error())
	}

    err = stub.PutState(stockKey, stockBytes)
    if err != nil {
            return shim.Error("Failed to edit stock")
    }
    //create composite key "offered~nameOfStock~ownerid" to make table of owners who offered stocks
    indexName = "offered~nameOfStock~ownerid"
    stockKey, err = stub.CreateCompositeKey(indexName, []string{"offered",sname, sowner})
    if err != nil {
         return shim.Error(err.Error())
    }
    //add owner who offer stocks
    value := []byte{0x00}
    stub.PutState(stockKey,value)
    
    return shim.Success([]byte("stocks added sucessfully"))
}

//args[0] => OwnerID,args[1] => nameOfstock, args[2] => priceOfstock, args[3] => offeredstocks
func (t *SimpleAsset) buyStocks(stub shim.ChaincodeStubInterface, args []string) peer.Response{
    if len(args) != 3 {
        return shim.Error("Incorrect arguments. Expecting no., name of stock and price")
    }
    BuyerId:=args[0]
    sname :=  args[1]
    sprice,err := strconv.ParseFloat(args[2],64)
    sofferedstocks,err := strconv.Atoi(args[3])
    //check if price not float or number not integer
    if err != nil {
        return shim.Error(err.Error())
    }
    //get client with his id
    clientBytes, err := stub.GetState(BuyerId)
    BuyerAccount := client{}
    c := client{}
    err = unMarshalClient(clientBytes,&c) 
	if err != nil {
    	return shim.Error(err.Error())
    }
    //check if client have enough blanance to offer
    temp2 := float64(sofferedstocks)*sprice
    if BuyerAccount.Balance < temp2{
        fmt.Println("client don't have enough money")
        return shim.Error("client don't have enough money")
    }
    //get buyer stock
    indexName := "nameOfStock~ownerid"
    stockKey, err := stub.CreateCompositeKey(indexName, []string{sname, BuyerId})
    stockBytes,err :=stub.GetState(stockKey)
    BuyerStock := stock{}
    err=unMarshalStock(stockBytes,&BuyerStock)
    if err != nil {
    	return shim.Error(err.Error());
    }
    
    indexName = "offered~nameOfStock~ownerid"
    owners,err :=stub.GetStateByPartialCompositeKey(indexName,[]string{"offered",sname})
    if err !=nil{
        return shim.Error(err.Error());
    }
    i:=0
    for i=0;owners.HasNext();i++{
        skey,err := owners.Next()
        if err !=nil{
            return shim.Error(err.Error());
        }
        _,parts,err :=stub.SplitCompositeKey(skey.Key)
        if err !=nil{
            return shim.Error(err.Error());
        }
        ownerId:=parts[2]

        fmt.Println("Owner ID : ", ownerId)
        indexName = "nameOfStock~ownerid"
        stockKey, err = stub.CreateCompositeKey(indexName, []string{sname, ownerId})
        if err != nil {
            return shim.Error(err.Error());
        }
        stockBytes,err = stub.GetState(stockKey)
	    if err != nil {
    	    return shim.Error(err.Error());
        }
        //get stock of owner
        sellerStock := stock{}
        err=unMarshalStock(stockBytes,&sellerStock)
        if err != nil {
    	    return shim.Error(err.Error());
        }   
        //getting lowest price of stock // momken n3mlo optimize bas b3dyn :D 3a4an ehna gamdy f45
        if sellerStock.Price <= sprice{
            //check if client have enough stocks to offer
            if sellerStock.Offeredstocks <= sofferedstocks{
                stockKey, err := stub.CreateCompositeKey(indexName, []string{"offered",sname, ownerId})
                if err != nil {
                    return shim.Error(err.Error());
                }
                stub.PutState(stockKey,nil)
                BuyerStock.Offeredstocks += sellerStock.Offeredstocks
                BuyerAccount.Balance -= float64(sellerStock.Offeredstocks)*sellerStock.Price
            }else {
                BuyerStock.Offeredstocks += sofferedstocks
                BuyerAccount.Balance -= float64(sofferedstocks)*sellerStock.Price
            }
            sellerStock.Offeredstocks -= sofferedstocks 

            //save new Buyer stock
            indexName = "nameOfStock~ownerid"
            stockKey, err = stub.CreateCompositeKey(indexName, []string{sname, BuyerId})
            if err != nil {
                return shim.Error(err.Error());
            }
            BuyerBytes,err := json.Marshal(BuyerStock)
            stub.PutState(stockKey,BuyerBytes)
            //save new Buyer Account
            accountBytes,err := json.Marshal(BuyerAccount)
            stub.PutState(BuyerId,accountBytes)
            //save seller Stock
            indexName = "nameOfStock~ownerid"
            stockKey, err = stub.CreateCompositeKey(indexName, []string{sname, ownerId})
            if err != nil {
                return shim.Error(err.Error());
            }
            SellerBytes,err := json.Marshal(sellerStock)
            stub.PutState(stockKey,SellerBytes)
        }
    }
    return shim.Success([]byte("stocks sold sucessfully"))
}

//args[0]=> source id ,args[1]=> destination id,args[2] => amount 
func (t *SimpleAsset) transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    if len(args) != 3 {
        return shim.Error("Incorrect arguments. Expecting a key and a value")
    }
    source := [...]string{args[0],args[2]}
    dest := [...]string{args[1],args[2]}

    //check if client 1 exists
    exist,sClientbytes := checkClientExist(stub,source[0])
    if exist==false{
        fmt.Println("Source client does not exist")
        return shim.Error("Source client does not exist")
    }

    //check if client 2 exists
    exist,dClientbytes := checkClientExist(stub,dest[0])
    if exist==true{
        fmt.Println("Destination client does not exist")
        return shim.Error("Destination client does not exist")
    }
    sClient:=client{}
    dClient:=client{}
    unMarshalClient(sClientbytes,&sClient)
    unMarshalClient(dClientbytes,&dClient)
    if sClient == (client{}) || dClient == (client{}){
        return shim.Error("there was an error in UnMarshal bytes")
    }
    amount,err :=strconv.ParseFloat(args[2], 64)
    if err !=nil{
        return shim.Error("args3 not float")
    }
    if sClient.Balance >= amount {
        sClient.Balance -= amount
        dClient.Balance +=amount
    }else{
        return shim.Error("sourceClient don't have engouh money")
    }
    
	fmt.Println("Transfer sucessfully")
	
    return shim.Success([]byte("Transfer success"))
}

//args[0]=> client id ,args[1]=> amount
func (t *SimpleAsset) withdraw(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    
    //Check el id wl amount
    if len(args) != 2 {
            return shim.Error("Incorrect arguments. Expecting client ID and amount to be withdrawn")
    }
    //client
    clientBytes, err := stub.GetState(args[0])
    //Amonut
    withdrawamount, err := strconv.ParseFloat(args[1], 64)
    c := client{}
    err = json.Unmarshal(clientBytes,&c) 
	if err != nil {
    	return shim.Error(err.Error())
	}
    if withdrawamount > c.Balance {
		return shim.Error("Sorry You don't have enough money")
    }else{
        c.Balance -= withdrawamount
    }
    clientBytes,err = json.Marshal(c)
    if err != nil {
    	return shim.Error(err.Error())
    }
    err = stub.PutState(args[0], clientBytes)
    if err != nil {
    	return shim.Error(err.Error())
    }

	fmt.Println("Amount withdrawen= ", withdrawamount)
	log:=fmt.Sprintf(" Your new Balance" , c.Balance)
    fmt.Println(log)
    return shim.Success([]byte(log))
}

//args[0]=> client id ,args[1]=> amount
func (t *SimpleAsset) deposit(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    
    //Check el id wl amount
    if len(args) != 2 {
            return shim.Error("Incorrect arguments. Expecting client ID and amount to be deposit")
    }
    //client
    clientBytes, err := stub.GetState(args[0])
    //Amonut
    depositamount, err := strconv.ParseFloat(args[1], 64)
    c := client{}
    err = json.Unmarshal(clientBytes,&c) 
	if err != nil {
    	return shim.Error(err.Error())
	}
    c.Balance += depositamount

    clientBytes,err = json.Marshal(c)
    if err != nil {
    	return shim.Error(err.Error())
    }
    err = stub.PutState(args[0], clientBytes)
    if err != nil {
    	return shim.Error(err.Error())
    }

    fmt.Println("Your Deposit ",depositamount)
    log := fmt.Sprintf(" Your new Balance " , c.Balance)
    fmt.Println(log)
	return shim.Success([]byte(log))
}

// Return the balance of the specified client args[0]=> client id
func (t *SimpleAsset) getBalance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    if len(args) != 1 {
            return shim.Error("Incorrect arguments. Expecting a key (client ID)")
    }
    //client
    clientBytes, err := stub.GetState(args[0])
    if err != nil {
    	return shim.Error(err.Error())
    }
    c := client{}
    err = json.Unmarshal(clientBytes,&c) 
	if err != nil {
    	return shim.Error(err.Error())
    }
    value := c.Balance
    log := fmt.Sprintf("Your Balance= " , value)
	fmt.Println(log)
    return shim.Success([]byte(log))
}

func (t *SimpleAsset) createClient(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    if len(args) != 3 {
            return shim.Error("Incorrect arguments. Expecting a key and a value")
    }

    // arg0 = id  ///  arg1 = name /// arg2 = balance  
    cid := args[0]
    cname := args[1]
    cbalance,err := strconv.ParseFloat(args[2], 64)
    // check if the balance not float
    if err != nil {
    	return shim.Error(err.Error())
    }
    //check if client exist
    exist,_ := checkClientExist(stub,cid)
    if exist==true{
        fmt.Println("there is a client with same id ")
        return shim.Error("there is a client with same id ")
    }
    

    client1 := &client{cname,cid,cbalance}
	fmt.Println("Client ID = %s" , client1.Id)
	fmt.Println(" Client Name = %s" , client1.Name)
	fmt.Println(" Client Balance = " , client1.Balance)
	clientBytes,err := json.Marshal(client1) 
	if err != nil {
    	return shim.Error(err.Error())
	}
	fmt.Println("bytes" , clientBytes)

    err = stub.PutState(args[0], clientBytes)
    if err != nil {
            return shim.Error("Failed to set asset")
    }
    log := fmt.Sprintf("Sucessfully created client with id %s :) ",cid)
	fmt.Println(log)
	
	return shim.Success([]byte(log))
}

func checkClientExist(stub shim.ChaincodeStubInterface,id string) (bool,[]byte){
    clientBytes, err := stub.GetState(id)
    if err!=nil{
        return false,nil
    }
    if clientBytes !=nil{
        fmt.Println("There is a client with the same id!")
        return true,clientBytes
    }
    return false,nil
}

func unMarshalClient(clientBytes []byte , c *client) error{
    err := json.Unmarshal(clientBytes,&c) 
	if err != nil {
    	fmt.Println("There was an error in UnMarshal")
    }
    return err
}

func checkStocksExist(stub shim.ChaincodeStubInterface,key string) (bool,[]byte){
    stockBytes, err := stub.GetState(key)
    if err!=nil{
        return false,nil
    }
    if stockBytes !=nil{
        fmt.Println("There is a stock with the same owner!")
        return true,stockBytes
    }
    return false,nil
}

func unMarshalStock(stockBytes []byte , s *stock) error{
    err := json.Unmarshal(stockBytes,&s) 
	if err != nil {
    	fmt.Println("There was an error in UnMarshal")
    }
    return err
}

// main function starts up the chaincode in the container during instantiate
func main() {
    if err := shim.Start(new(SimpleAsset)); err != nil {
            fmt.Println("Error starting SimpleAsset chaincode: %s", err)
    }
}

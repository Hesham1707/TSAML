package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"encoding/json"
	"fmt"
	"math"
)

var logstring string


func (setup *FabricSetup) CheckAnomaly(clientID string) (string, error){
	
	logstring:= ""
	
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")


	fname := "queryClientLogs"
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: fname, Args: [][]byte{[]byte(clientID)}, TransientMap: transientDataMap})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	c := client{}

	var logs []log
	var depositLogs []log
	var withdrawLogs []log
	fmt.Println("clientID : ",clientID)
	err =json.Unmarshal(response.Payload,&c)
	if err != nil{
		fmt.Println("Can't Unmarshal logs")
		return "", err
	}
	fmt.Println(c.Name)
	logstring+="Showing logs for user: " + c.Name +" \n "
	logs =c.UserLogs
	i:=0
	for i=0;i<len(logs);i++{
		fmt.Println(logs[i].Type)
	}

	if len(logs) <= 8{
		logstring=setlogstring(logs,c.Name)
		resultstring:="Client did not make enough transactions to check anomaly, Required atleast 9 transactions \n Showing all transactions... \n" +logstring
		
		return resultstring,nil
	}
	
	for i=0;i<len(logs);i++{
		fmt.Println("%s \n",logs[i].toString)
		if logs[i].Type=="withdraw" {
			withdrawLogs=append(withdrawLogs,logs[i])
		}else if logs[i].Type=="deposit" {
			depositLogs=append(depositLogs,logs[i])
		}else if logs[i].Type=="transfer" {
			if logs[i].Owner==clientID{
				withdrawLogs=append(withdrawLogs,logs[i])
			}else if logs[i].NewOwner == clientID{
				depositLogs = append(depositLogs,logs[i])
			}
		}
	}
	if len(depositLogs)> 3 {
		fmt.Println("check deposit logs")
		dCheck,lstring:=checkTrend(depositLogs)
		
		if dCheck == 1{
			logstring+=lstring
			resultstring:="Suspicious deposit transaction detected \n Showing depoist transactions... \n" + logstring
			return resultstring,nil
		}else{
			kCheck,kstring:=checkTrendWithKalman(depositLogs,c.Balance)
			if kCheck == 1{
				logstring+=kstring
				fmt.Println("detect anomaly from kalman filter")
				resultstring:="Suspicious deposit transaction detected with kalman filter \n Showing depoist transactions... \n" + logstring
				return resultstring,nil
			} 
		}
	}


	if len(withdrawLogs)> 3{
		fmt.Println("check withdraw logs")
		wCheck,lstring:=checkTrend(withdrawLogs)
		
		if wCheck == 1{
			logstring+=lstring
			resultstring:="Suspicious withdraw transaction detected \n  Showing withdraw transactions... \n" + logstring
			return resultstring,nil
		}else{
			kCheck,kstring:=checkTrendWithKalman(withdrawLogs,c.Balance)
			if kCheck == 1{
				logstring+=kstring
				fmt.Println("detect anomaly from kalman filter")
				resultstring:="Suspicious withdraw transaction detected with Kalman filter \n  Showing withdraw transactions... \n" + logstring
				return resultstring,nil
			}
		}
	}
	logstring=setlogstring(logs,c.Name)
	resultstring:="No Suspicious transaction detected \n  Showing all transactions... \n" + logstring
	return resultstring,nil
}

func checkTrend(Transactions []log) (int,string) {
	logtext:=""
	sum := 0;
	length :=len(Transactions)
	count := 0
	found :=0
	if length > 15 {
		count = int(float64(length) * 0.2)
	}else {
		count =3
	}
	i:=0
	for  i=0 ; i< count ; i++{ 
		sum += int(Transactions[i].Amount)
		logtext += Transactions[i].toString()+" \n "
	}

	var expected int = sum/count 
	threshold := 10 * expected
	
	fmt.Println("\n Initial expected = ",expected)

	fmt.Println("\n Initial Threshold =  ",threshold)


	for i=count ; i<length ; i++ {
		fmt.Println("\n amount= ",Transactions[i].Amount)
		fmt.Println("\n expected= ",expected)
		fmt.Println("\n Threshold= ",threshold)
	
		if Abs(int(Transactions[i].Amount)-expected)<threshold{
			sum+=int(Transactions[i].Amount)
			expected =sum/(i+1)
			threshold = 10 * expected
			logtext += Transactions[i].toString()+" \n "
		}else {
			logtext +="******" +Transactions[i].toString()+"******"+" \n "
			found=1
		}
	}
	if  found == 1{
		return 1,logtext
	}else{
	return 0,logtext
	}
}

func checkTrendWithKalman(Transactions []log,balance float64) (int,string) {
	
	var sum float64 = 0;
	logtext:=""
	length :=len(Transactions)
	count := 0
	found := 0
	if length > 15 {
		count = int(float64(length) * 0.2)
	}else {
		count =3
	}
	Emea,Eest:=calcDeviation(Transactions)
	if Emea == 0{
		return 0,"no anomly"
	}
	Eest =math.Pow(Eest,2)
	expected := Transactions[0].Amount
	i:=0
	for  i=0 ; i< count ; i++{ 
		sum += Transactions[i].Amount
	}
	threshold := 10.0 * expected

	for i=0 ; i<length ; i++ {
		KG:=Eest/(Eest+Emea)
		expected=expected+KG*(Transactions[i].Amount-expected)
		Eest=(1-KG)*Eest
		log:=fmt.Sprintf("Eest: %0.2f ,Emea:  %0.2f,KG : %0.2f, expected : %0.2f",Eest,Emea,KG,expected)
		fmt.Println(log)
		if math.Abs(Transactions[i].Amount-expected)>threshold && i>=count{
			logtext +="******" +Transactions[i].toString()+"******"+" \n "
			found=1
		}else{
		logtext += Transactions[i].toString()+" \n "
		}
	}	
	if  found == 1{
		return 1,logtext
	}else{
	return 0,logtext
	}
}

func calcDeviation(Transactions []log) (float64,float64){
	var sum float64 = 0;
	length :=len(Transactions)
	var stDeviation float64 =0
	i:=0
	for i=0 ; i<length ; i++ {
		sum += Transactions[i].Amount
	}
	mean :=sum/float64(length)
	for i=0 ; i<length ; i++ {
		temp:=Transactions[i].Amount-mean
		stDeviation+=math.Pow(temp,2)
	}
	stDeviation= stDeviation/float64(length)
	return stDeviation,mean
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
	
type client struct{
	Name string  `json:"name"`
	Id string 	`json:"id"`
    Balance float64 `json:"balance"`
    UserLogNo int `json:"userlogno"`
    UserLogs []log `json:"userlog"`
}
// Asset of logss
type log struct{
    Type string `json:"type"`
    Owner string `json:"owner"`
    Amount float64 `json:"amount"`
    NewOwner string `json:"newOwner"`
    Quantity int `json:"quantity"`
    NameOfStock string `json:"nameOfStock"`

}
func (l *log) toString() string {
    logString:=fmt.Sprintf("%s amount = %0.2f" ,l.Type, l.Amount)
    return logString
}
func (l *log) toStringTransfer() string {
    logString:=fmt.Sprintf(" %s amount = %0.2f to ID = %s" ,l.Type, l.Amount,l.NewOwner)
    return logString
}
type names struct{
    Usernames []string `json:"usernames"`
}
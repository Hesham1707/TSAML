package blockchain

import (
	"math/rand"
	"strconv"
	
)

// Returns an int >= min, < max
func randomInt(min, max int) int {
    return min + rand.Intn(max-min)
}

func (setup *FabricSetup) makeTransactions (id string ,initBalance int , no int , start int , end int ){
	i := 0 ;
		for i = 0 ; i < no ; i++ {
			randTransaction := randomInt(0,2)
			randBalance := randomInt(start , end)
			if randTransaction == 0 {
	
				setup.InvokeHello("deposit" , []string{id ,strconv.Itoa(randBalance)})
				initBalance = initBalance + randBalance
				
			}else{
				if initBalance < randBalance {
				continue
				}else{
				setup.InvokeHello("withdraw" , []string{id ,strconv.Itoa(randBalance)})
				initBalance = initBalance - randBalance
				}
				
			}
	
	
		}
	
	}
	

	

func (setup *FabricSetup) GenerateClients (){
	noOfUsers := 1 
	noOftranse := randomInt(100 , 500)  
	
	i := 0
	for i = 0 ; i < noOfUsers ; i++ {
		name := "Abdo (do2do2)" //"User "+strconv.Itoa(i)
		randBalance := randomInt(1000 , 10000)
		setup.InvokeHello("createClient" , []string{strconv.Itoa(i) ,name ,strconv.Itoa(randBalance) })
		setup.makeTransactions("0" ,randBalance ,noOftranse , 1000 , 10000)
	}


	for i = 0 ; i < noOfUsers ; i++ {
		name := "Hesham" //"User "+strconv.Itoa(i)
		randBalance := randomInt(100 , 1000000)
		setup.InvokeHello("createClient" , []string{strconv.Itoa(i) ,name ,strconv.Itoa(randBalance) })
		setup.makeTransactions("1" ,randBalance ,noOftranse , 100 , 1000000)
	}

}


package blockchain

import (
	"math/rand"
	"strconv"
	
)

// Returns an int >= min, < max
func randomInt(min, max int) int {
    return min + rand.Intn(max-min)
}

func (setup *FabricSetup) makeTransactions (id string ,initBalance int , no int , start int , end int  , anom int ){
	
		makeanom := randomInt(0,no)

	i := 0 ;
		for i = 0 ; i < no ; i++ {
			if anom == 1 && makeanom == i {
				setup.InvokeHello("deposit" , []string{id ,strconv.Itoa(90000000)})
				continue
			}

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
	
		noOftranse := randomInt(50 , 100)  
	
	
	
		name := "Abdo (el gamd fash5)" //"User "+strconv.Itoa(i)
		randBalance := randomInt(1000 , 10000)
		setup.InvokeHello("createClient" , []string{strconv.Itoa(0) ,name ,strconv.Itoa(randBalance) })
		setup.makeTransactions(strconv.Itoa(0) ,randBalance ,noOftranse , 1000 , 10000 ,0)
	

	
		name = "Hesham (حلو صغير جميل مثير)" //"User "+strconv.Itoa(i)
		randBalance = randomInt(100 , 100000)
		setup.InvokeHello("createClient" , []string{strconv.Itoa(1) ,name ,strconv.Itoa(randBalance) })
		setup.makeTransactions(strconv.Itoa(1) ,randBalance ,noOftranse , 100 , 10000 ,1)
	

}


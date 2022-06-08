package main


import (
	"time"
	"fmt"
	"strconv"
)

//Basic list node to store channels
type Node struct{
	next *Node
	channel chan string
	id int
}

func main() {
	etcdServer := startEtcd()
	defer etcdServer.Close()

	etcdClient := getClient()
	defer etcdClient.Close()

	startWatchers(etcdClient)
  
	head := &Node{
		next: nil,
		channel: nil,
		id: 0,
	}
	
	head.channel = createNewVM(opts)
	if(head.channel == nil){
		fmt.Errorf("Error creating machine")
	}
	var cmd string
	for{
		fmt.Print("embr>")
		fmt.Scanln(&cmd)
		if(cmd == "exit"){
			head.channel <- "shutdown"
			for(head.next != nil){
				head = head.next
				head.channel<-"shutdown"
			}
			break
		}
		if(cmd == "startVM"){
			if(head.next == nil){
				temp := &Node{
					next: nil,
					channel: nil,
					id: head.id + 1,
				}
				temp.channel = createNewVM(opts)
				if(temp.channel == nil){
					fmt.Errorf("Error creating machine")
				}
				head.next = temp
			} else{
				temp := head
				for temp.next != nil{
					temp = temp.next
				}
				t := &Node{
					next: nil,
					channel: nil,
					id: temp.id + 1,
				}
				t.channel = createNewVM(opts)
				if(t.channel == nil){
					fmt.Errorf("Error creating machine")
				}
				temp.next = t
			}
		}
		if(cmd == "status"){
			if(head == nil){
				fmt.Println("No VM's currently running")
			}else{
				fmt.Println("VM with ID " + strconv.Itoa(head.id) + "running")
				if(head.next != nil){
					temp := head
					for temp.next != nil{
						temp = temp.next
						fmt.Println("VM with ID " + strconv.Itoa(temp.id) + "running")
					}
				}	
			}
		}
	}
	time.Sleep(2* time.Second)
  log.Fatal(<-etcdServer.Err()) //Blocking statement
}	


package tx

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"testing"
	"github.com/GrFrHuang/gox/log"
	"net/rpc"
	"net"
	"time"
	"strconv"
	"fmt"
)

func TestNewUpstream(t *testing.T) {
	db, err := sql.Open("mysql", "huang:baiying@tcp(127.0.0.1:3306)/app_server_test")
	if err != nil {
		log.Panic(err)
	}
	sqls := []string{"show create table news"}
	u := NewUpstream(db, "send_queue", "confirm_queue", "amqp://guest:guest@127.0.0.1:5672", 20, sqls)
	for i := 0; i < 5; i++ {
		err = u.StartTX()
		if err != nil {
			log.Error(err)
		}
	}
	u.ListenQuery()
	ch := make(chan int)
	<-ch
}

func TestNewDownstream(t *testing.T) {
	db, err := sql.Open("mysql", "huang:baiying@tcp(127.0.0.1:3306)/work_together")
	if err != nil {
		log.Panic(err)
	}
	sqls := []string{`insert into test_table (user_id,update_time) values (1,1525725181);`}
	d := NewDownstream(db, "send_queue", "confirm_queue", "amqp://guest:guest@127.0.0.1:5672", "", 1, 20, sqls)
	err = d.EndTX()
	if err != nil {
		log.Error(err)
	}
	d.ListenTimeout("127.0.0.1:1234")
	ch := make(chan int)
	<-ch
}

func TestRpc(t *testing.T) {
	//rpcServer := rpc.NewServer()
	//err := rpcServer.Register(new(Query))
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//listener, err := net.Listen("tcp", "127.0.0.1:1234")
	//if err != nil {
	//	log.Printf("net.Listen tcp :0: %v", err)
	//}
	//go rpc.Accept(listener)
	////ch := make(chan int)
	////<-ch
	//time.Sleep(5 * time.Second)
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Error(err)
		return
	}
	conn, _ := net.DialTCP("tcp", nil, addr)
	defer conn.Close()

	rpcClient := rpc.NewClient(conn)
	defer rpcClient.Close()

	arg := &Args{7, 8}
	result := make([]string, 10)
	err = rpcClient.Call("Query.CheckTxState", arg, &result)
	if err != nil {
		log.Error(err)
	}
	time.Sleep(time.Second * 3)
}

//RPC test
type Args struct {
	A, B int
}

type Bean int

func (t *Bean) CheckTxState(args *Args, reply *([]string)) error {
	*reply = append(*reply, strconv.Itoa(args.B), "GrFrHuang")
	return nil
}

func TestRpc2(t *testing.T) {
	newServer := rpc.NewServer()
	newServer.Register(new(Bean))

	lst, e := net.Listen("tcp", "127.0.0.1:5255") // any available address
	if e != nil {
		log.Printf("net.Listen tcp :0: %v", e)
	}

	//Listen tcp port for rpc request.
	go newServer.Accept(lst)
	//newServer.HandleHTTP("/foo", "/bar")

	time.Sleep(5 * time.Second)

	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5255")
	if err != nil {
		panic(err)
	}
	conn, _ := net.DialTCP("tcp", nil, address)
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	args := &Args{7, 8}
	result := make([]string, 10)
	err = client.Call("Bean.CheckTxState", args, &result)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println(result)
}

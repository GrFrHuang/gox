// Use rabbitmq and mysql local transaction to complete transaction's eventually.
// Downstream node need a eliminating redundancy table.
// Record place which may be general problem.
//
//BA：Basic Available 基本可用
//整个系统在某些不可抗力的情况下，仍然能够保证“可用性”，即一定时间内仍然能够返回一个明确的结果。只不过“基本可用”和“高可用”的区别是：
//“一定时间”可以适当延长
//当举行大促时，响应时间可以适当延长
//给部分用户返回一个降级页面
//给部分用户直接返回一个降级页面，从而缓解服务器压力。但要注意，返回降级页面仍然是返回明确结果。
//S：Soft State：柔性状态
//同一数据的不同副本的状态，可以不需要实时一致。
//E：Eventual Consisstency：最终一致性
//同一数据的不同副本的状态，可以不需要实时一致，但一定要保证经过一定时间后仍然是一致的。

package tx

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/GrFrHuang/gox/tools/rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
	"time"
	"strconv"
	"net/rpc"
	"net"
	"github.com/GrFrHuang/gox/log"
)

const (
	LOGTABLE = "idempotent_log"
	READY    = "ready"
	SUCCESS  = "success"
	FAIL     = "fail"
)

type Upstream struct {
	sql          []string // Mysql sql
	sendQueue    string   // Upstream member commit data to downstream member.
	confirmQueue string   // The queue is that can notify downstream's operation result.
	rabbitMQ     *rabbitmq.RabbitMQ
	db           *sql.DB
}

type Downstream struct {
	sql           []string
	sendQueue     string
	confirmQueue  string
	idempotentLog string
	retry         int64
	rabbitMQ      *rabbitmq.RabbitMQ
	confirmClient *rabbitmq.RabbitMQ
	timeoutClient *rabbitmq.RabbitMQ
	db            *sql.DB
}

// New a transaction's upstream member.
func NewUpstream(db *sql.DB, sendQueue, confirmQueue, mqUrl string, mqTimeout int64, sqls []string) (*Upstream) {
	upstream := &Upstream{
		sql:          sqls,
		sendQueue:    sendQueue,
		confirmQueue: confirmQueue,
		db:           db,
	}
	rabbitClient, err := rabbitmq.NewRabbit(mqUrl, mqTimeout)
	if err != nil {
		panic(err)
	}
	upstream.rabbitMQ = rabbitClient
	return upstream
}

// New a transaction's downstream member.
func NewDownstream(db *sql.DB, sendQueue, confirmQueue, mqUrl, IdempotentLog string, retry, mqTimeout int64, sqls []string) (*Downstream) {
	downstream := &Downstream{
		sql:          sqls,
		sendQueue:    sendQueue,
		confirmQueue: confirmQueue,
		retry:        retry,
		db:           db,
	}
	if IdempotentLog == "" {
		downstream.idempotentLog = LOGTABLE
	}
	if retry < 1 {
		downstream.retry = 1
	}
	rabbitClient, err := rabbitmq.NewRabbit(mqUrl, mqTimeout)
	if err != nil {
		panic(err)
	}
	downstream.rabbitMQ = rabbitClient
	// Confirm queue's mesaages never be overdue.
	confirmClient, err := rabbitmq.NewRabbit(mqUrl, 0)
	downstream.confirmClient = confirmClient
	if err != nil {
		panic(err)
	}
	timeoutClient, err := rabbitmq.NewRabbit(mqUrl, 0)
	if err != nil {
		panic(err)
	}
	downstream.timeoutClient = timeoutClient
	return downstream
}

func (u *Upstream) StartTX() (err error) {
	if len(u.sql) == 0 {
		err = fmt.Errorf("sql collection error: %v", u.sql)
		return
	}
	tx, err := u.db.Begin()
	if err != nil {
		return
	}
	for _, v := range u.sql {
		if v == "" {
			err = fmt.Errorf("sql error: %v", v)
			return
		}
		_, err = tx.Exec(v)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	// todo 如果本地事务提交了但是消息还没到rabbitMQ就断电断网?
	// Publish to downstream.
	body := []byte("")
	header := make(amqp.Table)
	header["id"] = strconv.Itoa(int(time.Now().UTC().UnixNano()))
	// isConfirm = true, ensure the message arrive rabbitMQ broker, then commit mysql local transaction.
	err = u.rabbitMQ.Publish("ex1", u.sendQueue, body, header, true)
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

type Query struct{}

// Downstream check rabbitmq's dead queue and judge upstream transaction's state.
func (q *Query) QueryTxState(arg string, reply *string) error {
	// Alter pointer's point this value, because reply has been initialized.
	*reply = SUCCESS
	if 0 == 0 {
		// Upstream operation successful, remove the timeout message.
		success()
	} else {
		// Upstream operation fail, take the message requeue.
		fail()
	}
	return nil
}

// Start listen downstream rpc client's query request.
// Default tcp port is 12345.
func (u *Upstream) ListenQuery() (error) {
	rpcServer := rpc.NewServer()
	rpcServer.Register(new(Query))
	lst, err := net.Listen("tcp", ":1234") // any available address
	if err != nil {
		log.Printf("net.Listen tcp :0: %v", err)
		return err
	}
	//Listen tcp port for rpc request.
	go rpcServer.Accept(lst)
	return nil
}

func (u *Upstream) GetSql() []string {
	return u.sql
}

func (u *Upstream) GetDB() *sql.DB {
	return u.db
}

func (d *Downstream) msgHandler(msg amqp.Delivery) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// If downstream's transaction execute fail, requeue it.
			msg.Reject(true)
		}
	}()
	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (upstream_id, create_time) VALUES (%s,%d);", d.idempotentLog, msg.Headers["id"].(string), time.Now().Unix()))
	if err != nil {
		return err
	}
	for _, v := range d.sql {
		if v == "" {
			err = fmt.Errorf("sql error: %v", d.sql)
			return err
		}
		_, err = tx.Exec(v)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// Publish to upstream.
	header := make(amqp.Table)
	header["id"] = msg.Headers["id"]
	body := []byte("ok")
	err = d.confirmClient.Publish("ex1", d.confirmQueue, body, header, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (d *Downstream) EndTX() (err error) {
	if len(d.sql) == 0 {
		err = fmt.Errorf("sql collection error: %v", d.sql)
		return
	}
	// Check whether idempotent log table exists or not.
	_, err = d.db.Exec("SHOW CREATE TABLE " + d.idempotentLog)
	if err != nil {
		// Create a new idempotent log table at the same database.
		_, err = d.db.Exec(fmt.Sprintf("CREATE TABLE %s (id int(20) NOT NULL AUTO_INCREMENT, upstream_id varchar(30) DEFAULT NULL, create_time int(11) DEFAULT NULL, PRIMARY KEY (id), UNIQUE KEY upstream_id_index (upstream_id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;", d.idempotentLog))
		if err != nil {
			return
		}
	}
	err = d.rabbitMQ.Receive("ex1", d.sendQueue, d.msgHandler)
	if err != nil {
		return
	}
	return
}

func (d *Downstream) GetSql() []string {
	return d.sql
}

func (d *Downstream) GetDB() *sql.DB {
	return d.db
}

// Set max retry numbers for deal fail.
func (d *Downstream) SetRetry(retry int64) {
	if retry < 1 {
		retry = 1
	}
	d.retry = retry
}

// Listen timeout queue always, if find, query state with upstream.
func (d *Downstream) ListenTimeout(upstreamRPCAddr string) error {
	var QueryState = func(msg amqp.Delivery) error {
		addr, err := net.ResolveTCPAddr("tcp", upstreamRPCAddr)
		if err != nil {
			return err
		}
		conn, _ := net.DialTCP("tcp", nil, addr)
		defer conn.Close()

		client := rpc.NewClient(conn)
		defer client.Close()

		id := msg.Headers["id"].(string)
		reply := ""
		err = client.Call("Query.QueryTxState", id, &reply)
		if err != nil {
			return err
		}
		if reply == SUCCESS {
			// Upstream operation successful, remove the timeout message.
			err = fmt.Errorf("Continue consume message: %v ", msg.Headers["id"])
			log.Info(err)
		} else {
			// Upstream operation fail, take the message requeue.
			log.Info("Upstream node execute transaction fail: ", msg.Headers["id"])
		}
		return err
	}
	err := d.timeoutClient.ListenTimeOut("ex1", QueryState)
	return err
}

func success() {
	log.Info(SUCCESS)
}

func fail() {
	log.Info(FAIL)
}

package main

type ServerPool struct {
	accessClients map[string]string
	maxClient     int // maximum web socket client number.
}

func NewServerPool(address string, maxClient int) {

}

//func (pool *ServerPool) RegisterClient(address, authFeild string) {
//	ws.accessClients = make(map[string]string, ws.maxClient)
//	ws.accessClients[address] = authFeild
//}

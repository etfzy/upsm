package connspool

// 将conns pool的创建交给用户，可自由创建https pool、tcp pool 或tcps pool
type ConnsBuilder func(addr string) (interface{}, error)

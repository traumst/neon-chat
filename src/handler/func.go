package handler

type FuncPerConn func(conn *Conn) error

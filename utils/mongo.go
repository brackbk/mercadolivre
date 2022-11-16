package utils

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

//Context -> struct
type Context struct {
	MongoSession *mgo.Session
}

var sessionDb *mgo.Session

//CreateSessionDb -> Create de session db
func CreateSessionDb() {

	var err error

	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dialInfo := mgo.DialInfo{
		Addrs: []string{
			os.Getenv("CLUSTER0"),
			os.Getenv("CLUSTER1"),
			os.Getenv("CLUSTER2"),
		},
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig) // add TLS config
		return conn, err
	}

	sessionDb, err = mgo.DialWithInfo(&dialInfo)

	if err != nil {
		log.Fatalf("[Session]: %s\n", err)
	}
}

//GetSessionDb -> get de Session Db if exist
func GetSessionDb() *mgo.Session {
	if sessionDb == nil {
		CreateSessionDb()
	}
	return sessionDb
}

//Close -> close de session Mongo
func (c *Context) Close() {
	c.MongoSession.Close()
}

//Collection -> get by collection name
func (c *Context) Collection(name string) *mgo.Collection {
	return c.MongoSession.DB(os.Getenv("DATABASE")).C(name)
}

//NewContext -> create a new Context
func NewContext() *Context {
	session := GetSessionDb().Copy()
	context := &Context{
		MongoSession: session,
	}
	return context
}

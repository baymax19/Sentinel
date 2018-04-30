package dbo

import (
	"log"

	. "github.com/than-os/socks-info-server/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Socks5DBO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "nodeinfo"
)

func (m *Socks5DBO) Connect() {

	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}

	db = session.DB(m.Database)
}

func (m *Socks5DBO) FindAll() ([]Node, error) {

	var nodes []Node
	err := db.C(COLLECTION).Find(bson.M{}).All(&nodes)

	return nodes, err
}

func (m *Socks5DBO) FindById(id string) (Node, error) {

	var node Node
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&node)

	return node, err
}

func (m *Socks5DBO) Insert(node Node) error {

	err := db.C(COLLECTION).Insert(&node)

	return err
}

func (m *Socks5DBO) Delete(node Node) error {

	err := db.C(COLLECTION).Remove(&node)
	return err

}

func (m *Socks5DBO) Update(node Node) error {

	err := db.C(COLLECTION).UpdateId(node.ID, &node)
	return err
}

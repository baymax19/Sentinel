package models

import "gopkg.in/mgo.v2/bson"

// Node is a struct defined for our socks5 node info
type Node struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	IPAddr        string        `bson:"ipAddr" json:"ipAddr"`
	PortPasswords PortPass      `bson:"portPasswords" json:"portPasswords"`
	Location      string        `bson:"location" json:"location"`
	WalletAddress string        `bson:"walletAddress" json:"walletAddress"`
	Method        string        `bson:"maethod" json:"method"`
}

//PortPass is Shadowsocks node Passwords, static right now
// but can be dynamic
type PortPass struct {
	PortPassword0 string `json:"4200"`
	PortPassword1 string `json:"4201"`
	PortPassword2 string `json:"4202"`
	PortPassword3 string `json:"4203"`
}

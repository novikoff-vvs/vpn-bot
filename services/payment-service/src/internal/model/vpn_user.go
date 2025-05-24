package model

type VPNUser struct {
	UUID         string       `yaml:"uuid" json:"uuid"`
	Subscription Subscription `yaml:"subscription" json:"subscription"`
}

package model

type Role int8

const (
	Host Role = iota
	Guest
)

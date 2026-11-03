package config

import (
	"github.com/celsiainternet/elvis/utility"
)

/**
* Load
* @param et.Json config
* @return error
**/
func Load(stage, packageName string) error {
	return nil
}

/**
* PasswordHash
* @param string password
* @return string
**/
func PasswordHash(password string) string {
	return utility.ToBase64(password)
}

/**
* PasswordUnhash
* @param string password
* @return string
**/
func PasswordUnhash(password string) string {
	return utility.FromBase64(password)
}

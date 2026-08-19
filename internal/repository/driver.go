package repository

import "strings"

type RepositoryDriver string

const (
	MySQLDriver = "mysql"
)

func GetDriver(driver string) RepositoryDriver {
	switch strings.ToLower(driver) {
	case "mysql":
		return MySQLDriver
	default:
		return MySQLDriver
	}
}

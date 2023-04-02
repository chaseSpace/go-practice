package adapayCore

import (
	"fmt"
	"log"
)

func Println(v ...interface{}) {
	log.Println("<AdapayLog> " + fmt.Sprint(v...))
}

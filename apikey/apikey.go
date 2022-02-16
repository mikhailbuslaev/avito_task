package apikey
import (	
	"time"
)
func Generate() time.Time{
	t := time.Now()
	return t
}
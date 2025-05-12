package jsondelete

import (
	"fmt"
	"syscall"
)


func DeleteFile(fn string)error{
	if err:=syscall.Unlink(fn);err!=nil{
		return fmt.Errorf("error during deleting file %v with value %v",fn,err)
	}
	return nil
}

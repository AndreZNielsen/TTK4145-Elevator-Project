package utility

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
	
)

func BoolToIntArrayConverter (boolArray [][]bool) [][]int {
	rows := len(boolArray)
	cols := len(boolArray[0])
	intArray := make([][]int,rows)
	for i:= range boolArray{
		intArray[i]=make([]int,cols)
		for j := range boolArray[i]{
			if boolArray[i][j] {
				intArray[i][j] = 1 
			}else{
				intArray[i][j] = 0
			}}
	}

	return intArray 
}
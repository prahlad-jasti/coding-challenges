package command

import (
	"bytes"
	"fmt"
	"redis-lite/resp"
	"redis-lite/util"
)

type Operation struct {
}

func (operation *Operation) Handle(buffer []byte) string {
	command := util.ClearAllZeroBytes(buffer)
	fmt.Println(string(command))
	checkByteHasTerminatorAtTheEndOfTheArray(command)
	splitCommand := bytes.Split(command, util.SeparatorBytes)
	val, _ := resp.ParseResp(splitCommand)
	//not sure of this
	strArr := util.ConvertInterfaceToStringArr(val)
	result := CommandFactory(strArr).execute(strArr)
	return result.(string)
}

func checkByteHasTerminatorAtTheEndOfTheArray(command []byte) {
	if command[len(command)-2] != util.CarriageReturn || command[len(command)-1] != util.NewLine {
		panic("wrong terminator operator for basic string")
	}
}
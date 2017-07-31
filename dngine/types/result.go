// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package types

import (
	"fmt"
)

// CONTRACT: a zero Result is OK.
type Result struct {
	Code CodeType
	Data []byte
	Log  string // Can be non-deterministic
}

type NewRoundResult struct {
}

type CommitResult struct {
	AppHash      []byte
	ReceiptsHash []byte
}

type ExecuteInvalidTx struct {
	Bytes []byte
	Error error
}

type ExecuteResult struct {
	ValidTxs   [][]byte
	InvalidTxs []ExecuteInvalidTx
	Error      error
}

func NewResult(code CodeType, data []byte, log string) Result {
	return Result{
		Code: code,
		Data: data,
		Log:  log,
	}
}

func (res Result) IsOK() bool {
	return res.Code == CodeType_OK
}

func (res Result) IsErr() bool {
	return res.Code != CodeType_OK
}

func (res Result) Error() string {
	return fmt.Sprintf("{code:%v, data:%X, log:%v}", res.Code, res.Data, res.Log)
}

func (res Result) String() string {
	return fmt.Sprintf("{code:%v, data:%X, log:%v}", res.Code, res.Data, res.Log)
}

func (res Result) PrependLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  log + ";" + res.Log,
	}
}

func (res Result) AppendLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  res.Log + ";" + log,
	}
}

func (res Result) SetLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  log,
	}
}

func (res Result) SetData(data []byte) Result {
	return Result{
		Code: res.Code,
		Data: data,
		Log:  res.Log,
	}
}

//----------------------------------------

// NOTE: if data == nil and log == "", same as zero Result.
func NewResultOK(data []byte, log string) Result {
	return Result{
		Code: CodeType_OK,
		Data: data,
		Log:  log,
	}
}

func NewError(code CodeType, log string) Result {
	return Result{
		Code: code,
		Log:  log,
	}
}

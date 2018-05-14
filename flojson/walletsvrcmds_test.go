// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2018 The Flo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package flojson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/bitspill/flod/flojson"
)

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return flojson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &flojson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return flojson.NewAddMultisigAddressCmd(2, keys, flojson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &flojson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   flojson.String("test"),
			},
		},
		{
			name: "addwitnessaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("addwitnessaddress", "1address")
			},
			staticCmd: func() interface{} {
				return flojson.NewAddWitnessAddressCmd("1address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"addwitnessaddress","params":["1address"],"id":1}`,
			unmarshalled: &flojson.AddWitnessAddressCmd{
				Address: "1address",
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return flojson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &flojson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return flojson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &flojson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "encryptwallet",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("encryptwallet", "pass")
			},
			staticCmd: func() interface{} {
				return flojson.NewEncryptWalletCmd("pass")
			},
			marshalled: `{"jsonrpc":"1.0","method":"encryptwallet","params":["pass"],"id":1}`,
			unmarshalled: &flojson.EncryptWalletCmd{
				Passphrase: "pass",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &flojson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &flojson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &flojson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &flojson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &flojson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &flojson.GetBalanceCmd{
				Account: nil,
				MinConf: flojson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBalanceCmd(flojson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &flojson.GetBalanceCmd{
				Account: flojson.String("acct"),
				MinConf: flojson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBalanceCmd(flojson.String("acct"), flojson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &flojson.GetBalanceCmd{
				Account: flojson.String("acct"),
				MinConf: flojson.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetNewAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &flojson.GetNewAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getnewaddress", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetNewAddressCmd(flojson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct"],"id":1}`,
			unmarshalled: &flojson.GetNewAddressCmd{
				Account: flojson.String("acct"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &flojson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetRawChangeAddressCmd(flojson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &flojson.GetRawChangeAddressCmd{
				Account: flojson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &flojson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: flojson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetReceivedByAccountCmd("acct", flojson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &flojson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: flojson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &flojson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: flojson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetReceivedByAddressCmd("1Address", flojson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &flojson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: flojson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &flojson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetTransactionCmd("123", flojson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &flojson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: flojson.Bool(true),
			},
		},
		{
			name: "getwalletinfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getwalletinfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetWalletInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getwalletinfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetWalletInfoCmd{},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return flojson.NewImportPrivKeyCmd("abc", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &flojson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  flojson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return flojson.NewImportPrivKeyCmd("abc", flojson.String("label"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &flojson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   flojson.String("label"),
				Rescan:  flojson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return flojson.NewImportPrivKeyCmd("abc", flojson.String("label"), flojson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &flojson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   flojson.String("label"),
				Rescan:  flojson.Bool(false),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return flojson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &flojson.KeyPoolRefillCmd{
				NewSize: flojson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return flojson.NewKeyPoolRefillCmd(flojson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &flojson.KeyPoolRefillCmd{
				NewSize: flojson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return flojson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &flojson.ListAccountsCmd{
				MinConf: flojson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewListAccountsCmd(flojson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &flojson.ListAccountsCmd{
				MinConf: flojson.Int(6),
			},
		},
		{
			name: "listaddressgroupings",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listaddressgroupings")
			},
			staticCmd: func() interface{} {
				return flojson.NewListAddressGroupingsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listaddressgroupings","params":[],"id":1}`,
			unmarshalled: &flojson.ListAddressGroupingsCmd{},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return flojson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &flojson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAccountCmd{
				MinConf:          flojson.Int(1),
				IncludeEmpty:     flojson.Bool(false),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAccountCmd(flojson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAccountCmd{
				MinConf:          flojson.Int(6),
				IncludeEmpty:     flojson.Bool(false),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAccountCmd(flojson.Int(6), flojson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAccountCmd{
				MinConf:          flojson.Int(6),
				IncludeEmpty:     flojson.Bool(true),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAccountCmd(flojson.Int(6), flojson.Bool(true), flojson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAccountCmd{
				MinConf:          flojson.Int(6),
				IncludeEmpty:     flojson.Bool(true),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAddressCmd{
				MinConf:          flojson.Int(1),
				IncludeEmpty:     flojson.Bool(false),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAddressCmd(flojson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAddressCmd{
				MinConf:          flojson.Int(6),
				IncludeEmpty:     flojson.Bool(false),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAddressCmd(flojson.Int(6), flojson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAddressCmd{
				MinConf:          flojson.Int(6),
				IncludeEmpty:     flojson.Bool(true),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return flojson.NewListReceivedByAddressCmd(flojson.Int(6), flojson.Bool(true), flojson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &flojson.ListReceivedByAddressCmd{
				MinConf:          flojson.Int(6),
				IncludeEmpty:     flojson.Bool(true),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return flojson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &flojson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: flojson.Int(1),
				IncludeWatchOnly:    flojson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewListSinceBlockCmd(flojson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &flojson.ListSinceBlockCmd{
				BlockHash:           flojson.String("123"),
				TargetConfirmations: flojson.Int(1),
				IncludeWatchOnly:    flojson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewListSinceBlockCmd(flojson.String("123"), flojson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &flojson.ListSinceBlockCmd{
				BlockHash:           flojson.String("123"),
				TargetConfirmations: flojson.Int(6),
				IncludeWatchOnly:    flojson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return flojson.NewListSinceBlockCmd(flojson.String("123"), flojson.Int(6), flojson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &flojson.ListSinceBlockCmd{
				BlockHash:           flojson.String("123"),
				TargetConfirmations: flojson.Int(6),
				IncludeWatchOnly:    flojson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return flojson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &flojson.ListTransactionsCmd{
				Account:          nil,
				Count:            flojson.Int(10),
				From:             flojson.Int(0),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewListTransactionsCmd(flojson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &flojson.ListTransactionsCmd{
				Account:          flojson.String("acct"),
				Count:            flojson.Int(10),
				From:             flojson.Int(0),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return flojson.NewListTransactionsCmd(flojson.String("acct"), flojson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &flojson.ListTransactionsCmd{
				Account:          flojson.String("acct"),
				Count:            flojson.Int(20),
				From:             flojson.Int(0),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return flojson.NewListTransactionsCmd(flojson.String("acct"), flojson.Int(20),
					flojson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &flojson.ListTransactionsCmd{
				Account:          flojson.String("acct"),
				Count:            flojson.Int(20),
				From:             flojson.Int(1),
				IncludeWatchOnly: flojson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return flojson.NewListTransactionsCmd(flojson.String("acct"), flojson.Int(20),
					flojson.Int(1), flojson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &flojson.ListTransactionsCmd{
				Account:          flojson.String("acct"),
				Count:            flojson.Int(20),
				From:             flojson.Int(1),
				IncludeWatchOnly: flojson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return flojson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &flojson.ListUnspentCmd{
				MinConf:   flojson.Int(1),
				MaxConf:   flojson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewListUnspentCmd(flojson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &flojson.ListUnspentCmd{
				MinConf:   flojson.Int(6),
				MaxConf:   flojson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return flojson.NewListUnspentCmd(flojson.Int(6), flojson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &flojson.ListUnspentCmd{
				MinConf:   flojson.Int(6),
				MaxConf:   flojson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return flojson.NewListUnspentCmd(flojson.Int(6), flojson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &flojson.ListUnspentCmd{
				MinConf:   flojson.Int(6),
				MaxConf:   flojson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []flojson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return flojson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1}]],"id":1}`,
			unmarshalled: &flojson.LockUnspentCmd{
				Unlock: true,
				Transactions: []flojson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "move",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("move", "from", "to", 0.5)
			},
			staticCmd: func() interface{} {
				return flojson.NewMoveCmd("from", "to", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5],"id":1}`,
			unmarshalled: &flojson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     flojson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "move optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("move", "from", "to", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewMoveCmd("from", "to", 0.5, flojson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6],"id":1}`,
			unmarshalled: &flojson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     flojson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "move optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("move", "from", "to", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return flojson.NewMoveCmd("from", "to", 0.5, flojson.Int(6), flojson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6,"comment"],"id":1}`,
			unmarshalled: &flojson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     flojson.Int(6),
				Comment:     flojson.String("comment"),
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return flojson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &flojson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     flojson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewSendFromCmd("from", "1Address", 0.5, flojson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &flojson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     flojson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return flojson.NewSendFromCmd("from", "1Address", 0.5, flojson.Int(6),
					flojson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &flojson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     flojson.Int(6),
				Comment:     flojson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return flojson.NewSendFromCmd("from", "1Address", 0.5, flojson.Int(6),
					flojson.String("comment"), flojson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &flojson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     flojson.Int(6),
				Comment:     flojson.String("comment"),
				CommentTo:   flojson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return flojson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &flojson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     flojson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return flojson.NewSendManyCmd("from", amounts, flojson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &flojson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     flojson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return flojson.NewSendManyCmd("from", amounts, flojson.Int(6), flojson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &flojson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     flojson.Int(6),
				Comment:     flojson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return flojson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &flojson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return flojson.NewSendToAddressCmd("1Address", 0.5, flojson.String("comment"),
					flojson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &flojson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   flojson.String("comment"),
				CommentTo: flojson.String("commentto"),
			},
		},
		{
			name: "setaccount",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("setaccount", "1Address", "acct")
			},
			staticCmd: func() interface{} {
				return flojson.NewSetAccountCmd("1Address", "acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"setaccount","params":["1Address","acct"],"id":1}`,
			unmarshalled: &flojson.SetAccountCmd{
				Address: "1Address",
				Account: "acct",
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return flojson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &flojson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return flojson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &flojson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return flojson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &flojson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    flojson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []flojson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return flojson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &flojson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]flojson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    flojson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []flojson.RawTxInput{}
				privKeys := []string{"abc"}
				return flojson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &flojson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]flojson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    flojson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []flojson.RawTxInput{}
				privKeys := []string{}
				return flojson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					flojson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &flojson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]flojson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    flojson.String("ALL"),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return flojson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &flojson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return flojson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &flojson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return flojson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &flojson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := flojson.MarshalCmd(testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = flojson.MarshalCmd(testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request flojson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = flojson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}

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

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
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
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("addnode", "127.0.0.1", flojson.ANRemove)
			},
			staticCmd: func() interface{} {
				return flojson.NewAddNodeCmd("127.0.0.1", flojson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &flojson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: flojson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []flojson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return flojson.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123}],"id":1}`,
			unmarshalled: &flojson.CreateRawTransactionCmd{
				Inputs:  []flojson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []flojson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return flojson.NewCreateRawTransactionCmd(txInputs, amounts, flojson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &flojson.CreateRawTransactionCmd{
				Inputs:   []flojson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: flojson.Int64(12312333333),
			},
		},

		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &flojson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return flojson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &flojson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &flojson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetAddedNodeInfoCmd(true, flojson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &flojson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: flojson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &flojson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockCmd("123", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &flojson.GetBlockCmd{
				Hash:      "123",
				Verbose:   flojson.Bool(true),
				VerboseTx: flojson.Bool(false),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				// Intentionally use a source param that is
				// more pointers than the destination to
				// exercise that path.
				verbosePtr := flojson.Bool(true)
				return flojson.NewCmd("getblock", "123", &verbosePtr)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockCmd("123", flojson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true],"id":1}`,
			unmarshalled: &flojson.GetBlockCmd{
				Hash:      "123",
				Verbose:   flojson.Bool(true),
				VerboseTx: flojson.Bool(false),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblock", "123", true, true)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockCmd("123", flojson.Bool(true), flojson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true,true],"id":1}`,
			unmarshalled: &flojson.GetBlockCmd{
				Hash:      "123",
				Verbose:   flojson.Bool(true),
				VerboseTx: flojson.Bool(true),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &flojson.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &flojson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &flojson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: flojson.Bool(true),
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &flojson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := flojson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return flojson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &flojson.GetBlockTemplateCmd{
				Request: &flojson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := flojson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return flojson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &flojson.GetBlockTemplateCmd{
				Request: &flojson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := flojson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return flojson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &flojson.GetBlockTemplateCmd{
				Request: &flojson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &flojson.GetChainTipsCmd{},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &flojson.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &flojson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &flojson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &flojson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetInfoCmd{},
		},
		{
			name: "getmempoolentry",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getmempoolentry", "txhash")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetMempoolEntryCmd("txhash")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getmempoolentry","params":["txhash"],"id":1}`,
			unmarshalled: &flojson.GetMempoolEntryCmd{
				TxID: "txhash",
			},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &flojson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &flojson.GetNetworkHashPSCmd{
				Blocks: flojson.Int(120),
				Height: flojson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetNetworkHashPSCmd(flojson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &flojson.GetNetworkHashPSCmd{
				Blocks: flojson.Int(200),
				Height: flojson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetNetworkHashPSCmd(flojson.Int(200), flojson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &flojson.GetNetworkHashPSCmd{
				Blocks: flojson.Int(200),
				Height: flojson.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetRawMempoolCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &flojson.GetRawMempoolCmd{
				Verbose: flojson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetRawMempoolCmd(flojson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &flojson.GetRawMempoolCmd{
				Verbose: flojson.Bool(false),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &flojson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: flojson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetRawTransactionCmd("123", flojson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &flojson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: flojson.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &flojson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: flojson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return flojson.NewGetTxOutCmd("123", 1, flojson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &flojson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: flojson.Bool(true),
			},
		},
		{
			name: "gettxoutproof",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gettxoutproof", []string{"123", "456"})
			},
			staticCmd: func() interface{} {
				return flojson.NewGetTxOutProofCmd([]string{"123", "456"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"]],"id":1}`,
			unmarshalled: &flojson.GetTxOutProofCmd{
				TxIDs: []string{"123", "456"},
			},
		},
		{
			name: "gettxoutproof optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gettxoutproof", []string{"123", "456"},
					flojson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			staticCmd: func() interface{} {
				return flojson.NewGetTxOutProofCmd([]string{"123", "456"},
					flojson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"],` +
				`"000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"],"id":1}`,
			unmarshalled: &flojson.GetTxOutProofCmd{
				TxIDs:     []string{"123", "456"},
				BlockHash: flojson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &flojson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &flojson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return flojson.NewGetWorkCmd(flojson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &flojson.GetWorkCmd{
				Data: flojson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return flojson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &flojson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return flojson.NewHelpCmd(flojson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &flojson.HelpCmd{
				Command: flojson.String("getblock"),
			},
		},
		{
			name: "invalidateblock",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("invalidateblock", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewInvalidateBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"invalidateblock","params":["123"],"id":1}`,
			unmarshalled: &flojson.InvalidateBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return flojson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &flojson.PingCmd{},
		},
		{
			name: "preciousblock",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("preciousblock", "0123")
			},
			staticCmd: func() interface{} {
				return flojson.NewPreciousBlockCmd("0123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"preciousblock","params":["0123"],"id":1}`,
			unmarshalled: &flojson.PreciousBlockCmd{
				BlockHash: "0123",
			},
		},
		{
			name: "reconsiderblock",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("reconsiderblock", "123")
			},
			staticCmd: func() interface{} {
				return flojson.NewReconsiderBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"reconsiderblock","params":["123"],"id":1}`,
			unmarshalled: &flojson.ReconsiderBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return flojson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &flojson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     flojson.Int(1),
				Skip:        flojson.Int(0),
				Count:       flojson.Int(100),
				VinExtra:    flojson.Int(0),
				Reverse:     flojson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return flojson.NewSearchRawTransactionsCmd("1Address",
					flojson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &flojson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     flojson.Int(0),
				Skip:        flojson.Int(0),
				Count:       flojson.Int(100),
				VinExtra:    flojson.Int(0),
				Reverse:     flojson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return flojson.NewSearchRawTransactionsCmd("1Address",
					flojson.Int(0), flojson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &flojson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     flojson.Int(0),
				Skip:        flojson.Int(5),
				Count:       flojson.Int(100),
				VinExtra:    flojson.Int(0),
				Reverse:     flojson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return flojson.NewSearchRawTransactionsCmd("1Address",
					flojson.Int(0), flojson.Int(5), flojson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &flojson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     flojson.Int(0),
				Skip:        flojson.Int(5),
				Count:       flojson.Int(10),
				VinExtra:    flojson.Int(0),
				Reverse:     flojson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return flojson.NewSearchRawTransactionsCmd("1Address",
					flojson.Int(0), flojson.Int(5), flojson.Int(10), flojson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &flojson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     flojson.Int(0),
				Skip:        flojson.Int(5),
				Count:       flojson.Int(10),
				VinExtra:    flojson.Int(1),
				Reverse:     flojson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return flojson.NewSearchRawTransactionsCmd("1Address",
					flojson.Int(0), flojson.Int(5), flojson.Int(10), flojson.Int(1), flojson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &flojson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     flojson.Int(0),
				Skip:        flojson.Int(5),
				Count:       flojson.Int(10),
				VinExtra:    flojson.Int(1),
				Reverse:     flojson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return flojson.NewSearchRawTransactionsCmd("1Address",
					flojson.Int(0), flojson.Int(5), flojson.Int(10), flojson.Int(1), flojson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &flojson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     flojson.Int(0),
				Skip:        flojson.Int(5),
				Count:       flojson.Int(10),
				VinExtra:    flojson.Int(1),
				Reverse:     flojson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return flojson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &flojson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: flojson.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return flojson.NewSendRawTransactionCmd("1122", flojson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &flojson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: flojson.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return flojson.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &flojson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: flojson.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return flojson.NewSetGenerateCmd(true, flojson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &flojson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: flojson.Int(6),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return flojson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &flojson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return flojson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &flojson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := flojson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return flojson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &flojson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &flojson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "uptime",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("uptime")
			},
			staticCmd: func() interface{} {
				return flojson.NewUptimeCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"uptime","params":[],"id":1}`,
			unmarshalled: &flojson.UptimeCmd{},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return flojson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &flojson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return flojson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &flojson.VerifyChainCmd{
				CheckLevel: flojson.Int32(3),
				CheckDepth: flojson.Int32(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return flojson.NewVerifyChainCmd(flojson.Int32(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &flojson.VerifyChainCmd{
				CheckLevel: flojson.Int32(2),
				CheckDepth: flojson.Int32(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return flojson.NewVerifyChainCmd(flojson.Int32(2), flojson.Int32(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &flojson.VerifyChainCmd{
				CheckLevel: flojson.Int32(2),
				CheckDepth: flojson.Int32(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return flojson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &flojson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
		},
		{
			name: "verifytxoutproof",
			newCmd: func() (interface{}, error) {
				return flojson.NewCmd("verifytxoutproof", "test")
			},
			staticCmd: func() interface{} {
				return flojson.NewVerifyTxOutProofCmd("test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifytxoutproof","params":["test"],"id":1}`,
			unmarshalled: &flojson.VerifyTxOutProofCmd{
				Proof: "test",
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
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
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

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &flojson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &flojson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        flojson.Error{ErrorCode: flojson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &flojson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        flojson.Error{ErrorCode: flojson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%[3]v), "+
				"want %T", i, test.name, err, test.err)
			continue
		}

		if terr, ok := test.err.(flojson.Error); ok {
			gotErrorCode := err.(flojson.Error).ErrorCode
			if gotErrorCode != terr.ErrorCode {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.ErrorCode)
				continue
			}
		}
	}
}

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/pkg/kv"
	"github.com/pingcap/tidb/pkg/tablecodec"
	"github.com/pingcap/tidb/pkg/types"
	"github.com/pingcap/tidb/pkg/util/codec"
)

func main() {
	// init flag

	tableId := flag.Int64("table_id", 0, "select tidb_table_id from information_schema.tables where table_schema=<schema_name> and table_name=<table_name>;")
	col1Str := flag.String("binary_column", "0x0", "binary column value")
	col2 := flag.Int64("int_column", 0, "int column value")
	flag.Parse()

	// trim 0x or 0X for col1Str
	if len(*col1Str) > 2 && (*col1Str)[0] == '0' && ((*col1Str)[1] == 'x' || (*col1Str)[1] == 'X') {
		*col1Str = (*col1Str)[2:]
	}
	col1, err := hex.DecodeString(*col1Str)
	if err != nil {
		panic(errors.Trace(err))
	}
	tablePrefix := tablecodec.GenTableRecordPrefix(*tableId)
	val, err := codec.EncodeKey(time.UTC, nil, types.NewBytesDatum(col1), types.NewIntDatum(*col2))
	if err != nil {
		panic(errors.Trace(err))
	}
	//ch, err := kv.NewCommonHandle(encoded)
	handle, err := kv.NewCommonHandle(val)
	if err != nil {
		panic(errors.Trace(err))
	}
	key := tablecodec.EncodeRecordKey(tablePrefix, handle)

	cmd := exec.Command("tiup", "ctl:v7.5.1", "tikv", "--to-escaped", key.String())
	stdout, err := cmd.Output()
	if err != nil {
		panic(errors.Trace(err))
	}

	cmd = exec.Command("tiup", "ctl:v7.5.1", "tikv", "--encode", strings.TrimSuffix(string(stdout), "\n"))
	stdout, err = cmd.Output()
	if err != nil {
		panic(errors.Trace(err))
	}

	cmd = exec.Command("tiup", "ctl:v7.5.1", "tikv", "--to-escaped", strings.TrimSuffix(string(stdout), "\n"))
	stdout, err = cmd.Output()
	if err != nil {
		panic(errors.Trace(err))
	}
	fmt.Printf("z%s\n", strings.TrimSuffix(string(stdout), "\n"))
}

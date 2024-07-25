# ONCALL

This is a oncall tool.

# Build
```
go build main.go
```

# Dependency
```
tiup install ctl:v7.5.1
```

# Step
2. get table id
```
select tidb_table_id from information_schema.tables where table_schema=<schema_name> and table_name=<table_name>;
```

1. get mvcc key
```
./main --table_id=<table_id> --binary_column=<binary_column_val> --int_column=<int_column_val>
```

1. check mvcc info
```
tiup ctl:nightly tikv --host 127.0.0.1:20160 mvcc -k <mvcc_key> --show-cf=lock,write,default
```

1. get timestamp
```
select TIDB_PARSE_TSO(<tso>);
```
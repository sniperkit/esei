# esei
Elasticsearch batch export/import tool

## How to use it?

### Export

**If you have username/password, so you should type these.**

```
./esei -esurl=http://[ES domain / ES ip:port] -index=[Index Name] -user=[User Name] -passwd=[Password] -type=[Type Name] 
```

Then ESEI will save receive data to out.json. If you want save other file, type -out specify output file.

Also ESEI default get 10 records from ElasticSearch. If you want to get more records, type -size=[]. 
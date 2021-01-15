# bitsy - Bitset Data Indexer

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/bitsy)](https://goreportcard.com/report/github.com/viant/bitsy)
[![GoDoc](https://godoc.org/github.com/viant/bitsy?status.svg)](https://godoc.org/github.com/viant/bitsy)


This library is compatible with Go 1.13+

- [Motivation](#Motivation)
- [Introduction](#Introduction)
- [Indexing column types](#Indexing-column-types)
- [Rules](#Rules)
- [CLI usage](#Bitsy-CLI)
- [Use in SQL](#SQL-example)
- [Deployment](#Deployment)
    - [Deployment Configuration](#Deployment-Configuration)
    - [Deploy with endly](#Deploy-with-endly)
- [License](#license)


## Motivation

In Big Query, using repeated field values as part of conditional expressions in SQL queries may lead to expensive execution
plans both in terms of money and run time.

This project is a tool to improve Big Query performance and reduce runtime costs up to 100x, depending on the query
nature, by generating bitmap indexes for repeated and scalar columns/fields.


## Introduction

Bitsy improves query performance by building bitmap indexes on the columns that are part of 
an SQL query conditional expression.

The implementation uses 64 bit bitsets to indicate a specific value presence in one (or several) records 
out of 64 based on the record sequential number.  To allow indexing more than 64 rows/records, the source 
dataset has to be split into as many batches as necessary to represent the entire dataset. In addition 
to the indexing columns and the record sequential number, the source dataset should include the batch number, 
therefore.  Indexing is implemented as a Google cloud function and uses rules to describe index build requests.

A typical workflow happens as follows:  

- export the source data indexing columns along with the batch id column and a timestamp in the json format.

- index data source (bitsy) and save the output in the json format.

- import index data into a table (one per each index) for subsequent use in SQL conditional 
expressions.


## Indexing column types

The implementation supports the following data types for indexing columns:

- Integer
- Floating
- Text
- Boolean

Each column can be either a scalar or an array of multiple values.


## Rules

An indexing rule is specified in the yaml format.  A typical rule may look like this:

```yaml
when:
  prefix: /Users/
  suffix: json

dest:
  URL: mem://localhost/index/case01/$fragment/data.json
  TableRoot: myTable_

timeField: ts
batchField: batch_id
sequenceField: seq
partitionField: part_id
allowQuotedNumbers: true
recordsField: records
valueField: value
indexingFields:
  - Name: city_id
    Type: int
  - Name: name
    Type: string
```

- batchField: the source data field name containing batch number/id
- sequenceField: the source data field name containing sequence number/id
- recordsField: the output data bitset field name
- valueField: the output data indexed value field name


## Bitsy CLI

Provided CLI allows to:

```bash
# Show usage:
./bitsy -h

# Validate a rule
./bitsy -V -r valid.yaml

# Generate a rule
./bitsy -s test_data/data.json -d /tmp/bitsy/$fragment/data.json    -b batchId  -q seq  -f x:string  -f y:int
```

## SQL example

```sql
SELECT * FROM t1
WHERE EXISTS (SELECT 1 FROM t1_idx WHERE value IN ('idx1', 'idx2') AND t1.batch_id = t1_idx.batch_id and  AND t1_idx.events & (1 << t1.seq) != 0)
```

- table t1_idx contains the generated index
- column t1_idx.value contains index values
- column t1_idx.events contains bitsets


## Deployment

The following are used by storage mirror services:

**Prerequisites**

- config.json : main deployment configuration file 
- _$configBucket_: bucket storing rules
- _$triggerBucket_: bucket storing data that needs to be indexed, event triggered by GCP
- _$opsBucket_: bucket storing corrupted and failed input data for subsequent troubleshooting




##### Deployment Configuration

The main configuration file contains cloud function settings and the $configBucket, $triggerBucket,
$opsBucket locations. E.g.:
```json
{
  "DeadlineReductionMs": 70000,
  "LoaderDeadlineLagMs": 5000,
  "MaxRetries": 5,
  "Concurrency": 100,
  "RetryURL": "gs://${triggerBucket}/",
  "FailedURL": "gs://${operationBucket}/${appName}/failed",
  "CorruptionURL": "gs://${operationBucket}/${appName}/corrupted",
  "MaxExecTimeMs": 540000,
  "BaseURL" : "gs://${configBucket}/${appName}/Rules",
  "CheckFrequencyMs" : 100
}
```


##### Deploy with endly
To deploy with endly automation runner use the following workflow:

```bash
git checkout https://github.com/viant/bitsy.git
cd bitsy/deployment
endly authWith=myGoogleSecrets.json
```

## License

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.

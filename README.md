# Bitset Data Indexer (bitsy)


## Motivation

The project attempts to improve Big Query performance and reduce runtime costs.  


## Introduction

Bitsy improves query performance by building bitmap indexes on the columns that are part of 
an SQL query conditional expression.

The implementation uses 64 bit bitsets to indicate a specific value presence in one (or several) records 
out of 64 based on the record sequential number.  To allow indexing more than 64 rows/records, the source 
dataset has to be split into as many batches as necessary to represent the entire dataset. In addition 
to the indexing columns and the record sequential number, the source dataset should include the batch number, 
therefore.  Indexing is implemented as a Google cloud function and uses rules to describe index building requests.

A typical workflow happens as follows:  

- export the source data indexing columns along with the batch id column and a timestamp in the json format.

- index data source (bitsy) and save the output in the json format.

- import index data into a table (one per each index) for subsequent use in SQL conditional 
expressions.


## Indexing column types

The implementation supports the following data types for indexing columns:

- Numeric
- Floating
- Text
- Boolean

Each column can be either a scalar or an array of multiple value.


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
indexingFields:
  - Name: city_id
    Type: int
  - Name: name
    Type: string
```

## Bitsy CLI

Provided CLI allows to:

```bash
# Show usage:
./bitsy -h

# Validate a rule
./bitsy -V -r valid.yaml

# Genarate a rule
./bitsy -s test_data/data.json -d /tmp/bitsy/$fragment/data.json    -b batchId  -q seq  -f x:string  -f y:int
```

## SQL example

```sql
select * from t1
where EXISTS (SELECT 1 FROM t1_idx WHERE value IN ('idx1', 'idx2') AND t1.batch_id = t1_idx.batch_id and  AND t1_idx.events & (1 << t1.seq) != 0)
```

- table t1_idx contains the generated index
- column t1_idx.value contains index values
- column t1_idx.events contains bitsets

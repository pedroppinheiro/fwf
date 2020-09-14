# fwf

![CircleCI](https://img.shields.io/circleci/build/github/pedroppinheiro/fwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/pedroppinheiro/fwf)](https://goreportcard.com/report/github.com/pedroppinheiro/fwf)
[![GoDoc](https://godoc.org/github.com/pedroppinheiro/fwf?status.svg)](https://godoc.org/github.com/pedroppinheiro/fwf)

fwf is a little command line tool to make it easier to work with fixed-width files. It helps you visualize the fields on each line by generating an html which highlights and adds tooltips for each field. The records and fields of a fixed-width file must be given to fwf in a yaml file

## Usage
```
Usage fwf:
  -file string
        the full path for the file to generate the visualization
  -yaml string
        the full path for the yaml configuration
  -o string
        the path to where the exported file should be created (default "./")
```

Let's use the following fixed-width file "people.txt" as an example:

```
John Smith         40
Homer Simpson      30
Foo Bar            20
```

It has field "name" with size 19 (position 1 to 19), and a field "age" with size 2 (position 20 to 21). Let's create a yaml file "configuration.yaml" to represent the layout of this fixed-width file:

```
records:
  - name: "People"
    regex: .*
    fields:
      - name: "Person Name"
        initial: 1
        end: 19
      - name: "Age"
        initial: 20
        end: 21
```

All we have to to now is feed these files to our fwf command line with the following command:

```
./fwf -yaml="configuration.yaml" -file="people.txt"
```

The fwf tool will generate an index.html file which highlights fields. If you hover your mouse over the fields a tooltip will show up with the name of the fields.

## Building
A good command to certify that everything is working and building is the following:

```
go get -v -t -d ./... && go test ./... && go build && ./fwf -yaml="test_assets/test.yaml" -file="test_assets/test.txt" && cat index.html
```
or
```
go get -v -t -d ./... && go test ./... && go build && ./fwf -yaml="test_assets/test.yaml" -file="test_assets/test.txt" -o="./test_assets/" && cat ./test_assets/index.html
```

# fwf

fwf is a little command line tool to make it easier to work with fixed-width files. It helps you visualize the fields on each line by generating an html with tooltips for each defined field. The records and fields of a fixed-width file must be given to fwf in a yaml file

A good command to certify that everything is working is the following:

```
go test ./... && go build && ./fwf -yaml="test_assets/test.yaml" -file="test_assets/test.txt" && cat index.html
```
or
```
go test ./... && go build && ./fwf -yaml="test_assets/test.yaml" -file="test_assets/test.txt" -o="./test_assets/" && cat ./test_assets/index.html
```

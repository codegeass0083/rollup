### About the App
Those rows which have less than 4 columns or invalid year/month/day will be skipped.  

### Steps to run this program

##### Install
Follow the instructions on the link below and install GO Version 1.9 on your Linux system 
https://golang.org/doc/install

To verify it is successfully installed, run command `go version` on your terminal. 

##### Compile 
`go build rollup.go`

##### Run this app
`./rollup ./input.tsv y m d` will aggregate over all the prefixes of [y, m, d]  
`./rollup ./input.tsv y m` will aggregate over all the prefixes of [y, m]  
`./rollup ./input.tsv y` will aggregate over all the prefixes of [y]  
`./rollup ./input.tsv` same as `./rollup ./input.tsv y m d`  


# use official Golang image
FROM golang:1.16.3-alpine3.13

# set working directory
WORKDIR /app

# copy the source code
COPY . . 

#download and install all dependencies 
RUN go get -d -v ./...

# build the Go app
RUN go build -o avito .

# expose the port
EXPOSE 8000

# run the executable
CMD [ "./avito" ]


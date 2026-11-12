# syntax=docker/dockerfile:1
FROM golang:1.26.8-alpine

#Set destination for COPY
WORKDIR /app

#Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

#Build the docker image
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/credential-service


#Bind to port 50051
EXPOSE 50051

#RUN 
CMD ["/app/server"]
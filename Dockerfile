FROM golang:latest

# Install packages
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get install -y git nodejs

# Copy in source and install deps
WORKDIR /app-monitoring-archiver
ADD . .
RUN npm install -g serverless && npm install

WORKDIR /app-monitoring-archiver/lambda
RUN go get ./..


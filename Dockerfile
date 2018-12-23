
# golang:onbuild automatically copies the package source, 
# fetches the application dependencies, builds the program, 
# and configures it to run on startup 
FROM golang:onbuild
LABEL Name=keyauth Version=0.0.1 
EXPOSE 8080

# For more control, you can copy and build manually
# FROM golang:latest 
# LABEL Name=keyauth Version=0.0.1 
# RUN mkdir /app 
# ADD . /app/ 
# WORKDIR /app 
# RUN go build -o main .
# EXPOSE 8080 
# CMD ["/app/main"]

# console-chat

### HowTo run project
Before starting the authorization server, you must start all docker containers
```
cd auth-server/
docker-compose up -d
go run cmd/main.go 
```
Start chat server
```
cd chat-server/
go run cmd/main.go 
```
Build chat client
```
cd chat-client
make build
```
Create users with USER and ADMIN privileges. 
Note: only ADMIN can create new chat
```
./chat-client create user -u alex -p 123 -c 123 -r USER
./chat-client create user -u oleg -p 123 -c 123 -r ADMIN
```
Admin creates a new chat and connects to it. The chat ID is displayed in the terminal
```
./chat-client create chat -u alex -p 123
chat was created with id: c075bb8f-bebd-11ef-84c0-4c52622b5cbd
```
Another user can connect to an existing chat knowing this ID
```
./chat-client connect chat -u alex -p 123 -c c075bb8f-bebd-11ef-84c0-4c52622b5cbd
```
After this, you can send messages to each other from both clients. 
Press Ctrl-C to exit

**Aqua Runtime Assignment**

I have created a simple REST API that includes basic functionality such as Get Policy, Update Policy and more.
I have tried to keep it simple, reading from a configuration file and also used an interface in case we want to replace PostgreSQL tomorrow for something else.

To check the functionality please:
1. Run on terminal the following command:
   docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
   This will create a container running a postgres instance that you can see in your favorite GUI client (I used TablePlus).

2. build (`go build -o bin/RuntimeAssignment2`) and run (`go run .`) the server.

3. Use the Postman Collection file I have provided.

Thanks

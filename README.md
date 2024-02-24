**Aqua Runtime Assignment**

I have created a simple REST API that includes basic functionality such as Get Policy, Update Policy and more.

I have tried to keep it simple, reading from a configuration file and also used an interface in case we want to replace PostgreSQL tomorrow for something else.

I also added a WebSocket endpoint (running on ws://localhost:8088/ws) to notify connected clients (enforcers),
which can be tested by adding "WebSocket Test Client" from the Chrome extension store, and running the relevant postman calls (either Create or Update).

Regarding the Bonus Challenge, I would have either implemented a solution that can consist of another (indexed) column that inclueds a version of that policy: 
could be an md5 consisting of the whole json body, or a Unix Timestamp. Then keep a map of an id and version, and querying the db myself in some interval.
If on request, the client sends me his id and version for the policy, and it is the same on both the client and my map, there is no need to query. Otherwise, I query the db on the sent Id.
Or the use of a flag, if the value has changed since it was set, because NotifyCliens has been called.
Upon request, before querying the db, I will check the flag and if not, I send the cached value instead.

To check the functionality please:
1. Run on terminal the following command:
   
   `docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres`
   
   This will create a container running a postgres instance that you can see in your favorite GUI client (I used TablePlus).

3. build (`go build -o bin/RuntimeAssignment2`) and run (`go run .`) the server.

4. Use the Postman Collection file I have provided.

Thanks;

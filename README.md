# url-shortener

System Design URL Shortener in Golang

Client: Users accessing the web application.
Load Balancer: Distributes incoming traffic across multiple web servers.
Web Servers: Serve the web content and handle client requests.
Analytics MicroService: Collects and analyzes data from the web servers.
Database: Stores the application's data.
Cache Database: Stores frequently accessed data for faster retrieval.

```mermaid
graph TD
A[Client] --> B[Load Balancer]
B --> C[Web Servers]
C --> D[Database]
C --> E[Cache Database]
C --> F[Analytics MicroService]
B --> A
```

The Client sends a request to the Load Balancer.
The Load Balancer distributes the request to one of the Web Servers.
The Web Servers process the request, which may involve reading from or writing to the Database.
The Web Servers can also read from or write to the Cache Database to improve performance.
The Web Servers send relevant data to the Analytics MicroService for analysis.
The Web Servers return the response to the Client via the Load Balancer.

Client ⟶ Load Balancer ⟶ Web Servers
Web Servers ⟶ Database
Web Servers ⟶ Cache Database
Web Servers ⟶ Analytics MicroService
Load Balancer ⟶ Client (response)

## CQRS (Command Query Responsibility Segregation)

CQRS is a software architecture pattern that separates the responsabilities of reading (queries) and writing (commands) data in an application.

Benefits:
- You can improve performance and scalability since each model can be optimized independently (e.g., use a fast NoSQL DB for reads, a transactional SQL DB for writes).
- Separation enables clearer design, especially in complex or large-scale systems.
- Flexibility due to its natural supports for eventual consistency, event sourcing, and asynchronous processing.

#### When to use this pattern?

High performance systems, microservices, and domains with divergent read/write patterns.

e.g: Imagine an e-commerce system with 100 writes per second vs 100.000 reads per second. You may want to separate the reads from the writes.

#### Example

Imagine **a real time sales dashboard** for an e-commerce website.

In a No-CQRS scenario you would potentially have an orders table with millions of rows. Every time a user makes a purchase a row is added to the orders table. Now, imagine that the customer success team wants to have a dashboard with the total amount of sales per day, as well as the average spendings per user.

So if you perform a **SUM**, **GROUP BY**, and **COUNT** over the table every time someone refreshes the dashboard it could easily turn into a problem.

Here's where CQRS comes to the rescue! 

For the **command side**, you can use a relational database (PostgreSQL) to store every order. Optimized for fast transactions.

For the **query side**, ou can use a NoSQL database, or document based database like Redis, or MongoDB.




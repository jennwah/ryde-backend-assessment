# Ryde Backend

## Overview
Service is built with <b>Go</b> with efficient, minimalistic web framework <b>Gin</b>. It features user APIs with common CRUD actions, establish users
friendship relation and ability to retrieve nearby friends within a specified radius meter provided.

## Prerequisites
1. Go (uses `1.23.2` version for current repo)
2. Docker (optional, but it makes bootstrapping easier)
3. PostgreSQL (optional, only if you do not have docker setup locally)

## Instructions
The service is packaged with Docker build and as such, you can simply run `make run` to bring up the API service and our database service together. Once up, access the APIs via http://localhost:8080/api/v1

If you have no docker installed locally, you can still run the service with `go run cmd/service.go`, but you must have `PostgreSQL` installed (either via homebrew or other means) and running locally. 

After getting services up, we can run initial database migrations with `make migrate_up`.

At this point, you can start exploring with the APIs to create users or friendship. We also have `seed` utility for quicker exploration;

Run `make seed_data` to get a few users and friendship inserted into your database.

## API
The APIs are designed with Restful principels and documented with swagger tools. Here're the summary of APIs; 

1. `POST   /api/v1/users`            
2. `GET    /api/v1/users/:id`
3. `PATCH  /api/v1/users/:id`
4. `DELETE /api/v1/users/:id`
5. `POST   /api/v1/friends`
6. `GET    /api/v1/friends/nearby` 

Though, a more detailed API specifications can be viewed locally with the help of swagger tools, simply run `make run` (if you haven't started your services) and access via `http://localhost:8080/swagger/index.html`

## Codebase structure
I thought it's good to mention about how our codebase structure is setup. It follows Clean Architecture principle, where our service is divided into few core layers;

1. <b>Controller layer</b>

This is mostly your API controllers where it received request first-hand and responsible for common first-hand requests operations like request validation, auth middleware strategy etc. It then pass on its responsibility over to the next year.

2. <b>UseCase layer</b>

All business logic should reside on this layer for good encapsulation purposes. It should deal with most of every business / application logic to great extent.

3. <b>Repository layer</b>

A layer where service interacts with storage like your SQL / NoSQL databases. It deals with data saving, updating, deletion, retrieval etc making its responsibility clear to maintainers of repository.

4. <b>Others</b>

There're others useful mini packages that maintainers can define, as long as it make sense to have it on codebase. For example, the `pkg` directory hosts all implementation to third-party useful packages like postgresql client, in our case.

Additionally, each layer interacts with each other through the use of `interface` to ensure good <b>dependency injection practice</b> and ease of test mocks during unit tests.

## Persistence
SQL database (PostgreSQL) is chosen in this case. 

## Tests
Unit tests are covered with each application's layers. Maintainers are to be advised to write unit tests for good engineering practice. Run `make test` to run all unit tests and the test coverage will be output.

As of current writing, total test coverage is 75%.

## CI
Github CI actions are included as checks for two stages;

1. Build stage

Simply build our Go application as an containerized image and ready to be uploaded to any container register or shipped to anywhere, where a container orchestration platform would support; eg `k8s`

2. Test stage

Runs our unit tests with `make test` and output coverage reports.

## About the challenge (advanced part)
This service covers the basic requirements of challenge which covers User CRUD operations and other technical details. Specifically on the <b>advanced requirement parts</b>, this service was design to also handle

1. Complete application logging

We use Go's standard library `log/slog` logger package with JSON output format to `os.stdout`. As such, it would be compatible with every other Observability tools like ELK, or cloud platform native logs aggregator tools for full observability purposes. Simply collect
the logs from `os.stdout` and developers gain access to application logs tracing. 

2. User friendship feature

This feature is designed with relational join on our `users.friends` table with `user_id` and `friend_id` being the unique pair on table. Only one pair of `(user_id, friend_id)` could exist with the constraint
`user_id < friend_id` during inserts. This way, we ensure consistent friendship establishment between two users. 

It's good to note that, for a simple requirement like a friendship relation between two users, the current solution could scale reasonably well. However, a proper graph database like `neo4j` would scale better when a more complex
query like finding mutual friends, advanced graph traversal, friends of friends etc are to be implemented. But based on this assignment, the current solution would be sufficient for simple use cases.

3. Nearby friends feature

Once we have the friends feature, we can now extend on top of it with location (latitude, longtitude) parameters. The design is to use PostgreSQL's Postgis extension, specifically on this geospatial queries support.

Here's our `users` table design; 

```sql
CREATE TABLE IF NOT EXISTS users.users
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    date_of_birth DATE         NOT NULL,
    address       TEXT         NOT NULL,
    description   TEXT,
    location      GEOGRAPHY(POINT, 4326), -- Stores latitude and longitude in WGS 84
    created_at    TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_users_location ON users.users USING GIST (location);
```

We create GIST index on location column, which stores user's supplied latitude and longitude in WGS 84 (a standard used in Earth coordinates system). The `GIST` index internally uses an R-tree-like datastructure that would divide spaces 
into bounding boxes, which would reduce our search query as now query can be focused on, say specified radius. We will see this in our API get nearby friends soon.

Let's take a look on our Get Nearby Friends databse query; 

```sql
		WITH user_location AS (
			SELECT location
			FROM users.users
			WHERE id = $1
		)
		SELECT 
			u.id AS id,
			u.name,
			u.date_of_birth,
			u.address,
			u.description,
			public.ST_X(u.location::public.geometry) AS longitude,
			public.ST_Y(u.location::public.geometry) AS latitude
		FROM friends f
		JOIN users u ON 
			(f.user_id = $1 AND f.friend_id = u.id) OR
			(f.friend_id = $1 AND f.user_id = u.id)
		WHERE public.ST_DWithin(
			u.location,
			(SELECT location FROM user_location),
			$2 -- Radius in meters
		);
```

The part where we would get friends within certain radius meters would be on `ST_DWithin` function. It's a PostGis spatial function to determine
where two geometries are within a certain distance from each other (in our case, our users' geometry points). `ST_DWithin` would first check if both geometries intersect with given distance and because of 
`GIST` index we applied, this is fast. Later on, after filtering intersections, `ST_DWithin` would then calculate Euclidean distance between them.

Running `EXPLAIN ANALYZE` on above query would show indeed GIST index was used in filtering efficiently nearby coordinates;

```
  ->  Index Scan using idx_users_location on users u  (cost=0.26..20.78 rows=1 width=632) (actual time=0.598..0.609 rows=2 loops=1)
          Index Cond: (location && _st_expand((InitPlan 1).col1, '500000000'::double precision))
          Filter: st_dwithin(location, (InitPlan 1).col1, '500000000'::double precision, true)
```

For the purpose of this challenge, to scale our above solution for location filtering query when say our users get to a reasonably large scale, we can perform partitions based on users' locations with a simple partition on region-based,

```sql
    region_id     INT GENERATED ALWAYS AS (
        CASE
            WHEN ST_Within(location, ST_MakeEnvelope(-180, 0, 0, 90, 4326)::geography) THEN 1  -- North America
            WHEN ST_Within(location, ST_MakeEnvelope(0, 0, 180, 90, 4326)::geography) THEN 2  -- Europe/Asia
            WHEN ST_Within(location, ST_MakeEnvelope(-180, -90, 0, 0, 4326)::geography) THEN 3  -- South America
            ELSE 4 -- Others
        END
    ) STORED,
```
Create a region for each defined bounding boxes of geography in our world, depending on where our users are mostly residing. At global scale, it usefully make sense to parition by global regions, since nearby friends wouldn't make sense to retrieve friends from Asia, say you're from Europe.

However, in the case where users' location are always on update and dynamic (which is not requirement of this challenge, but let's explore it anyway), a in-memory database would scale better in such cases. For example, Redis's Geospatial data structure.
Since, users' locations are always dynamic (say at a frequency of new updates every 5 minutes), an in-memory datastore would not be an issue when it comes to persistency, as it could just mean users' location are lost for at most 5 minutes, after which new location data will be sent.
In any cases, Redis' Geospatial operations `GEOSEARCH` command provides a fairly logarithm operation on searching within radius of all users. 

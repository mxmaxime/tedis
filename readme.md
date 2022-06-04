## Todo
- implement a redis repository for Clusters.
    - take the same interface as the redis_repo and use it.
- list of connected clients

- detail of a hash:
    - json view
    - be able to edit the json view -> open it in your editor and when it quits I read the content to send it back to Redis.
        - it will be fun to implement :)
        - maybe add a step to validate changes before saving it in Redis? Using a diff?
    - a table with two columns: keys/values.

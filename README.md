# Hiring test. Golang users REST API

This repository contains Golang service implemented REST API with simple functional for manage users.


Using Golang implement REST API that provide functionality of users registration and authentication (using Bearer auth scheme).


API should provide next methods:


- POST /register
- POST /login
- GET /users
- GET /users/&amp;amp;lt;user-id&amp;amp;gt;
- PUT /users/&amp;amp;lt;user-id&amp;amp;gt;
- GET /users/me


- Anauthorized user can only /register or /login
- User could have such roles: Admin, User.
- User with role User can only perform users/me request.
- User with role Admin can perform all requests.
- User with role Admin could change other users roles calling PUT /users/&amp;amp;lt;user-id&amp;amp;gt; (you can make first admin manually in database)


Feel free to use any framework and database (we use mongo but you can choose what you are better in e.g. sqlite, postgres, tinydb etc)


Containerize your app with docker and include Dockerfile - and optionally docker-compose.yaml - to repository.
Push your code to Github/Bitbucket and provide us a link to repository.


extra tasks: deploy your solution to any cloud platform (aws, heroku etc); write unit-tests


estimated time (without extra tasks): 1-3 hours

## How to build and launch server
1) Make sure you have installed Docker and Docker Compose
2) Run `make serve.docker`
3) By default, server listen 6999 port and you can make requests like `http://127.0.0.1:6999/login`
4) Default superuser is `admin` with password `pass`:
```
curl -X POST http://127.0.0.1:6999/login/ --data-raw '{
    "identity": "admin",
    "secret": "pass"
}'
```
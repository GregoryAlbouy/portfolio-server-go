# portfolio2020-server-go

Backend-server for my portfolio https://gregoryalbouy.com written in Go + SQL (previously in Node.js + MongoDB)

It is functionnal and deployed on Google App Engine. I use it to access, store and manage project information using JWT authentication, and receive messages from the contact form – which is quite unlikely to happen.

## Features

- REST API providing CRUD operations on projects data
- JWT authentication
- Tests using Go's built-in tools

## Storage

I use SQLite for dynamic storage because an embedded solution suits well my usage, though SQL language is not my favorite to work with. I might switch to a more mongo-ish solution if I find a satisfying one.

## See also

* Portfolio website : https://gregoryalbouy.com
* Portfolio front repo : https://github.com/gregoryalbouy/portfolio-wcf-2020
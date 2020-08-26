# portfolio2020-server-go

Backend-server for my portfolio https://gregoryalbouy.com written in Go + SQL (previously in Node.js + MongoDB)

It is functionnal and deployed on Google App Engine. I use it to access, store and manage project information using JWT authentication, and receive messages from the contact form – which is quite unlikely to happen.

Portfolio : https://gregoryalbouy.com
Front repo : https://github.com/gregoryalbouy/portfolio-wcf-2020

## Storage

I use SQLite for dynamic storage because an embedded solution suits well my usage, though SQL language is not my favorite to work with. I wasn't aware until late of MongoDB equivalents such as [NeDB](https://github.com/louischatriot/nedb), so I might switch to it in the future.
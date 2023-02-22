<!-- PROJECT LOGO -->
<div align="center">
  <a href="https://kejith.de/">
    <img src="https://github.com/othneildrew/Best-README-Template/raw/master/images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Ginbar a custom made Imageboard with React & Redux</h3>
  <p>This is the backend part of Ginbar. It's fully written in Golang and Fiber is used as a web framework.</p>
</div>

<!-- ABOUT THE PROJECT -->
## About The Project
This is a custom made backend for the Ginbar Frontend. I exclusively used Golang for it's build-in concurrency system and performance.  

Here's why:
* Extremly Fast
* Scalable
* Learning Purposes 
* up to 2.000 concurrent User on a low-end Server

# ginbar-backend
## To-Dos 
- ~~rebase to Fiber (fasthttp)~~
- ~~make sessions permanent~~
- caching
  - [database caching](https://redis.io/documentation)
  - HTTP Response Caching
- refactor utilities to generalize them
- using generics for models and build a general API for models
- generics for posts
- Unit Testing
- decide whether SQLC is used in the future or not
  - look into OpenApi 3.0 and Swaggo

## Resources
- [gofiber/fiber](https://github.com/gofiber/fiber) on Github
- 

# Usefull Docker commands

- Build application: `docker build -t open-bar .`
- View image: `docker images | grep open-bar`
- Run application: `docker run -p 3000:3000 open-bar`
- Remove all stopped containers: `docker container prune`
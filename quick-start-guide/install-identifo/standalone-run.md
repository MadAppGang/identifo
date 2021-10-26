---
description: How to start with docker
---

# 🐳 Run with Docker

### Run in stateless mode

To run Idenitfo from the official Docker hub all you need is a simple command:

```
docker run -p 8081:8081 madappgangd/identifo:latest 
```

It runs identifo with all defaults settings and you can access:

* admin panel with the link [http://localhost:8081/adminpanel](http://localhost:8081/adminpanel)&#x20;
  * use default login and password: admin@admin.com/password to login to the admin panel
* login web app with the link [http://localhost:8081/web](http://localhost:8081/web/login)
* make API calls to http://localhost:8081/api

{% hint style="danger" %}
All the data is stored on a local docker temp file system, so all the changes to settings and data (users and apps) will be deleted after the docker image stops. You can mount volume or host fs to make the data persistent.
{% endhint %}

### Run and persist all the data on the host system

To run and persist data and config, we need to create a local folder on the host machine, download the default config file there and attach it to the docker image.

```bash
mkdir data #create the dir for indetifo data 
curl -o data/config.yaml https://raw.githubusercontent.com/MadAppGang/identifo/master/cmd/config-boltdb.yaml
```

This is all you need, now just run identifo with the command:

```
docker run \                                          
       -p 8081:8081 \
       --mount type=bind,source="$(pwd)"/data,target=/data \
       madappgangd/identifo:latest \
       --config=file://data/config.yaml
```

### Persist user and application

Identifo with default config keeps the boltdb database file in the `/data` folder. To make the data persistent between docker restarts you can mount this folder to the local filesystem or create docker volume. More information about Docker Volume [you can find in the docs.](https://docs.docker.com/storage/volumes/)

Use the following command to create and use docker named volume for that:

```bash
docker run \
       -p 8081:8081 \
       --mount type=volume,source=identifo-data,target=/data \
       madappgangd/identifo:latest
```

Where the docker will create and use volume with the name `identifo-data`. You can see the volume if you run `docker volume ls`

To get details about the docker volume use: `docker volume inspect identifo-data`

To delete the volume:` docker volume rm identifo-data`

Another option is to attach the host's folder to the docker image.&#x20;

```bash
docker run \
       -p 8081:8081 \
       --mount type=bind,source="$(pwd)"/data,target=/data \
       madappgangd/identifo:latest
```

This command will use `./data` folder from the current directory as `/data` folder in docker image. Just ensure the directory exists, docker will not create it for you.

### Persist the configuration file

To persist the data from the config file, you need to put it in persistent storage. You can grab the [default config file from github](https://raw.githubusercontent.com/MadAppGang/identifo/master/cmd/config-boltdb.yaml). And place it in persistent storage.

And then use the [`--config` flag ](../../settings/)to point to the config file you want to use.

Let's assume we want to use the same data folder as we are alredy using for data persistent.&#x20;

The options as persistent storage could be filesystem or s3. You can use a local file system or docker volume as FS persistent storage, the same way as we did it for data persistent storage.

```
docker run \                                          
       -p 8081:8081 \
       --mount type=bind,source="$(pwd)"/data,target=/data \
       madappgangd/identifo:latest \
       --config=file://data/config.yaml
```

Or use S3 storage:

```
docker run \                                          
       -p 8081:8081 \
       --mount type=bind,source="$(pwd)"/data,target=/data \
       madappgangd/identifo:latest \
       --config=s3://ap-southeast-2@my-bucket/identifo/config/custom-config.yaml
```


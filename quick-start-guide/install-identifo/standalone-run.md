---
description: How to use Identifo Docker image
---

# 🐳 Run with Docker

### Running in stateless mode

The simplest way to get acquainted with Identifo is to run the official Docker image:

```
docker run -p 8081:8081 madappgangd/identifo:latest 
```

This command launches Identifo with default settings and allows you to access:

* admin panel on [http://localhost:8081/adminpanel](http://localhost:8081/adminpanel)&#x20;
  * use default login `admin@admin.com`and default password `password`
* web app on [http://localhost:8081/web](http://localhost:8081/web/login)
* JSON API on [http://localhost:8081/api](http://localhost:8081/api)

{% hint style="danger" %}
All the configuration and application data is stored on a temporary filesystem in Docker container, so all your changes to settings and data (users and apps) will be lost after container stops. There are a couple of ways to address that (please see the sections below).
{% endhint %}

### Preserving state across container restarts

To introduce state to the above setup, [Docker bind mounts](https://docs.docker.com/storage/bind-mounts/) or [Docker volumes](https://docs.docker.com/storage/volumes/) can be used for either configuration or application data, or both; Identifo also allows for using AWS S3 a persistent configuration storage.

#### Preserving configuration changes

The simplest way to introduce state is to create a folder on the host machine, place a config file there and bind mount it to the Docker container. You can download the default config file from the Identifo repository and use it as is:

```bash
mkdir data #create a directory for Identifo data 
curl -o data/config.yaml https://raw.githubusercontent.com/MadAppGang/identifo/master/cmd/config-boltdb.yaml
```

Now run the image specifying a mount:

```
docker run \                                          
       -p 8081:8081 \
       --mount type=bind,source="$(pwd)"/data,target=/data \
       madappgangd/identifo:latest \
       --config=file://data/config.yaml
```

This way, all the updates made to configuration will be reflected in the `data/config.yaml` file on your host machine.

&#x20;[`--config` flag ](../../settings/)also accepts S3 prefixes as a valid config file location:&#x20;

```
docker run \                                          
       -p 8081:8081 \
       --mount type=bind,source="$(pwd)"/data,target=/data \
       madappgangd/identifo:latest \
       --config=s3://<my-region>@<my-bucket>/identifo/config/custom-config.yaml
```

#### Preserving users and applications

Default Identifo configuration uses BoltDB as a persistent storage for application data, with database file located under `/data` folder. Similar to the example with preserving configuration changes, to make application data survive container restarts, you can either mount the `/data` folder from the host filesystem, or create a Docker volume.

Use the following command to create and use a Docker-managed volume called `identifo-data`:

```bash
docker run \
       -p 8081:8081 \
       --mount type=volume,source=identifo-data,target=/data \
       madappgangd/identifo:latest
```

To get the details about this volume: `docker volume inspect identifo-data`.

To delete the volume:` docker volume rm identifo-data`.

Another option is to mount a host's folder when running the Docker image:

```bash
docker run \
       -p 8081:8081 \
       --mount type=bind,source="$(pwd)"/data,target=/data \
       madappgangd/identifo:latest
```

This command will use `./data` directory on your host filesystem as a `/data` folder in Docker container filesystem. Just ensure the directory exists on your host machine, Docker will not create it for you.


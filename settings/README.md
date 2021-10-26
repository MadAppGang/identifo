---
description: Project settings
---

# Settings

you can run Identifo with no settings set and it will run with `server-config.yml` file from the local folder.

Or you can specify a `--config` flag on start with three options.

### Local file config

To load data from a local file, you need to specify a file path with no scheme or with `file` scheme in it, something like:

```bash
./identifo --config=file://./custom-config.yml
./identifo --config=file:///root/user/custom-config.yml
./identifo --config=../custom-config.yml
./identifo --config=/root/user/custom-config.yml

```

### Load config from s3 folder

To load it form S3 folder you need to specify path with `s3` scheme in the following format: `s3://[region]@[bucket-name]/[file-key]`.

The region is optional but recommended.

You don't need to provide AWS credentials explicitly. You can use the default [AWS credentials chain](https://docs.aws.amazon.com/AWSJavaSDK/latest/javadoc/com/amazonaws/auth/DefaultAWSCredentialsProviderChain.html). Or you can use shared credentials (in `~/.aws/credentials` file).

The simplest way is to assign an AWS IAM role to your machine instance or export env variables.

Here is an example of running the identifo with a `s3` config.

```bash
export AWS_ACCESS_KEY_ID=ABDCS
export AWS_SECRET_ACCESS_KEY=FFDDDSSDDFFDSS
./identifo --config=s3://ap-southeast-2@my-bucket/identifo/config/custom-config.yaml
```

or using AWS shared credentials with profile (region will be read from identifo profile in the shared config file):

```bash
AWS_PROFILE=identifo ./identifo --config=s3://my-bucket/custom-config.yaml
```


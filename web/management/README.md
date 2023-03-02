# Management API

In addition to using the Admin panel(Dashboard), you can retrieve, create, update or delete users using the Management API. If you want to call the Management API directly, you will first need to generate the appropriate Access Token ID and Access Token Secret. 

There are couple of additional method to let you backend do some custom actions without compromising the private key. For example:
- create password reset token for user
- create invite token



The management API has `management` path prefix, so to make a call to it you have to user the following base URL:

```sh
http POST https://nativelogin.com/management/token/reset email=user@gmail.com
```



## Authenticating requests

All authentication methods are based on HMAC-SHA1 [ RFC 2104 - Keyed-Hashing for Message Authentication ](https://www.ietf.org/rfc/rfc2104.txt). 

The signature logic is based on AWS S3 REST request authentication:
- https://docs.aws.amazon.com/AmazonS3/latest/userguide/S3_Authentication2.html
- https://docs.aws.amazon.com/AmazonS3/latest/userguide/RESTAuthentication.html

The signature is constructed like that:

```
Signature = URL-Encode( Base64( HMAC-SHA1( YourSecretAccessKey, UTF-8-Encoding-Of( StringToSign ) ) ) );

StringToSign = HTTP-VERB + "\n" +
    Content-MD5 + "\n" +
    Content-Type + "\n" +
    Date + "\n" +
    Expires+ "\n" +
    HTTP-HOST 
```
The first few header elements of StringToSign (Content-Type, Date, and Content-MD5) are positional in nature. StringToSign does not include the names of these headers, only their values from the request. 

If a positional header called for in the definition of StringToSign is not present in your request (for example, Content-Type or Content-MD5 are optional for PUT requests and meaningless for GET requests), substitute the empty string ("") for that position.

Notice how the Signature is URL-Encoded to make it suitable for placement in the query string.

Expires - the time when the signature expires, specified as the number of seconds since the epoch (00:00:00 UTC on January 1, 1970). A request received after this time (according to the server) will be rejected.

There is no limits for Expires, it could be 1000 years long, but for security reason make it as short as possible. If your http request timeout is 30 seconds, there is no reason to make Expires more than `epoch_now_in_seconds()+30`.

HTTP-HOST - a string containing the domain (that is the hostname) followed by (if a port was specified) a ':' and the port of the URL. Not scheme included.

Here are some examples:

| URL | HTTP-HOST |
| --- | --- |
| http://google.com:443/email | google.com:443 |
| https://google.com:443/calendar | google.com:443 |
| https://google.com:8080/calendar | google.com:8080 |
| https://google.com/calendar | google.com |
| https://developer.mozilla.org/en-US/docs/Web/API/URL | developer.mozilla.org |

### Authentication examples

Get Request:
```http
GET /token/invite HTTP/1.1
Host: nativelogin.com
Date: Tue, 27 Mar 2022 19:36:42 +0000
Expires: 1175139620

```

StringToSign:
```http
GET\n
\n
\n
Tue, 27 Mar 2022 19:36:42 +0000\n
1175139620\n
nativelogin.com/token/invite
```

Post Request:
```http
POST /token/invite HTTP/1.1
Host: nativelogin.com
User-Agent: curl/7.15.5
Date: Tue, 27 Mar 2022 19:36:42 +0000
Content-MD5: 671d1a43130f6f9a041ab20ff3c8559f
content-type: application/json
Expires: 1175139620

{
    "email": "user@gmail.com"
}

```

StringToSign:
```http
POST\n
671d1a43130f6f9a041ab20ff3c8559f\n
application/json\n
Tue, 27 Mar 2022 19:36:42 +0000\n
1175139620\n
nativelogin.com/token/invite
```


## Architecture security

Beside request authentication the good idea would be to close management API to internal network by you firewall, reverse proxy or load balancer.

You can use architecture components to limit access to internal network only or to specific IP addresses.
# identifo.js
Identifo JS SDK using for:
- store session token
- generate URL's for auth flow
- provide API handlers for Identifo server http API

[Read more](./identifo.js/README.md)

# web-element
A simple web form that can be used in HTML as ES module or in your favorite framework. Example usage in [React](./demo/src/register-identifo.ts) or [HTML](../static/web/index.html)

[Web element API](web-element/readme.md)

# demo
Demo project that use *identifo.js* and *web-element*

# admin
Identifo admin interface

# scripts
Run [./update-admin.sh](./update-admin.sh) to build and update admin static folder.

To update static/web you need to run [./update-web.sh](./update-web.sh) when change *identifo.js* or *web-element* projects.
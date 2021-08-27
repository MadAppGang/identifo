---
description: main page and positioning
---

# Main page

## Positioning

Identifo is a cloud-native user authentication system, which provides a secure native user experience with zero development efforts and full customisation and extendability.

Key points:

* cloud-native
* ios native expierence
* android native expierence
* secure
* web server-side rendering
* web full client-side integration with API flow \(no iframes etc\)
* distributed support with JWT token
* OIDC support
* one line integration

{% tabs %}
{% tab title="React" %}
```jsx
import identifo from 'identifo.js';

identifo.init({
    url: "https://mydomain.com:123",
    app_id: "aabbccssddd",
}).login();


```
{% endtab %}

{% tab title="Swift iOS" %}
```swift
identifo.init()
identiof.login()
```
{% endtab %}

{% tab title="Kotlin Android" %}
```kotlin
identifo.init()
identiof.login()
```
{% endtab %}
{% endtabs %}

{% hint style="info" %}
Identifo is proudly created and supported by [MadAppGang](https://madappgang.com) and community. If you are missign any integration or customisation, we can do it for you as a consulting company.
{% endhint %}


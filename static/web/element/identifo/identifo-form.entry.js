import{r as t,c as i,h as e,g as s,H as o}from"./index-b085f682.js";var r,n,a,l;(n=r||(r={})).PleaseEnableTFA="error.api.request.2fa.please_enable",n.NetworkError="error.network",(l=a||(a={})).TFATypeApp="app",l.TFATypeSMS="sms",l.TFATypeEmail="email";class h extends Error{constructor(t){super((null==t?void 0:t.message)||"Unknown API error"),this.detailedMessage=null==t?void 0:t.detailed_message,this.id=null==t?void 0:t.id,this.status=null==t?void 0:t.status}}var c=Object.defineProperty,p=Object.getOwnPropertySymbols,d=Object.prototype.hasOwnProperty,u=Object.prototype.propertyIsEnumerable,f=(t,i,e)=>i in t?c(t,i,{enumerable:!0,configurable:!0,writable:!0,value:e}):t[i]=e,m=(t,i)=>{for(var e in i||(i={}))d.call(i,e)&&f(t,e,i[e]);if(p)for(var e of p(i))u.call(i,e)&&f(t,e,i[e]);return t},g=(t,i,e)=>new Promise(((s,o)=>{var r=t=>{try{a(e.next(t))}catch(t){o(t)}},n=t=>{try{a(e.throw(t))}catch(t){o(t)}},a=t=>t.done?s(t.value):Promise.resolve(t.value).then(r,n);a((e=e.apply(t,i)).next())}));class x{constructor(t,i){this.config=t,this.tokenService=i,this.defaultHeaders={"X-Identifo-Clientid":"",Accept:"application/json","Content-Type":"application/json"},this.catchNetworkErrorHandler=t=>{if("Network Error"===t.message||"Failed to fetch"===t.message||"Preflight response is not successful"===t.message||t.message.indexOf("is not allowed by Access-Control-Allow-Origin")>-1)throw console.error(t.message),new h({id:r.NetworkError,status:0,message:"Configuration error",detailed_message:`Please check Identifo URL and add "${window.location.protocol}//${window.location.host}" to "REDIRECT URLS" in Identifo app settings.`});throw t},this.checkStatusCodeAndGetJSON=t=>g(this,null,(function*(){if(!t.ok){const i=yield t.json();throw new h(null==i?void 0:i.error)}return t.json()})),this.baseUrl=t.url.replace(/\/$/,""),this.defaultHeaders["X-Identifo-Clientid"]=t.appId,this.appId=t.appId}get(t,i){return this.send(t,m({method:"GET"},i))}put(t,i,e){return this.send(t,m({method:"PUT",body:JSON.stringify(i)},e))}post(t,i,e){return this.send(t,m({method:"POST",body:JSON.stringify(i)},e))}send(t,i){const e=m({},i);return e.credentials="include",e.headers=m(m({},e.headers),this.defaultHeaders),fetch(`${this.baseUrl}${t}`,e).catch(this.catchNetworkErrorHandler).then(this.checkStatusCodeAndGetJSON).then((t=>t))}getUser(){return g(this,null,(function*(){var t,i;if(!(null==(t=this.tokenService.getToken())?void 0:t.token))throw new Error("No token in token service.");return this.get("/me",{headers:{Authorization:`Bearer ${null==(i=this.tokenService.getToken())?void 0:i.token}`}})}))}renewToken(){return g(this,null,(function*(){var t,i;if(!(null==(t=this.tokenService.getToken("refresh"))?void 0:t.token))throw new Error("No token in token service.");return this.post("/auth/token",{scopes:this.config.scopes},{headers:{Authorization:`Bearer ${null==(i=this.tokenService.getToken("refresh"))?void 0:i.token}`}}).then((t=>this.storeToken(t)))}))}updateUser(t){return g(this,null,(function*(){var i,e;if(!(null==(i=this.tokenService.getToken())?void 0:i.token))throw new Error("No token in token service.");return this.put("/me",t,{headers:{Authorization:`Bearer ${null==(e=this.tokenService.getToken("access"))?void 0:e.token}`}})}))}login(t,i,e,s){return g(this,null,(function*(){return this.post("/auth/login",{email:t,password:i,device_token:e,scopes:s}).then((t=>this.storeToken(t)))}))}federatedLogin(t,i,e,s){return g(this,arguments,(function*(t,i,e,s,o={width:600,height:800,popUp:!1}){var r=document.createElement("form");r.style.display="none",o.popUp&&(r.target="TargetWindow"),r.method="POST";const n=new URLSearchParams;if(n.set("appId",this.config.appId),n.set("provider",t),n.set("scopes",i.join(",")),n.set("redirectUrl",e),s&&n.set("callbackUrl",s),r.action=`${this.baseUrl}/auth/federated?${n.toString()}`,document.body.appendChild(r),o.popUp){const t=window.screenX+window.outerWidth/2-(o.width||600)/2,i=window.screenY+window.outerHeight/2-(o.height||800)/2;window.open("","TargetWindow",`status=0,title=0,height=${o.height},width=${o.width},top=${i},left=${t},scrollbars=1`)&&r.submit()}else window.location.assign(`${this.baseUrl}/auth/federated?${n.toString()}`)}))}federatedLoginComplete(t){return g(this,null,(function*(){return this.get(`/auth/federated/complete?${t.toString()}`).then((t=>this.storeToken(t)))}))}register(t,i,e){return g(this,null,(function*(){return this.post("/auth/register",{email:t,password:i,scopes:e}).then((t=>this.storeToken(t)))}))}requestResetPassword(t){return g(this,null,(function*(){return this.post("/auth/request_reset_password",{email:t})}))}resetPassword(t){return g(this,null,(function*(){var i,e;if(!(null==(i=this.tokenService.getToken())?void 0:i.token))throw new Error("No token in token service.");return this.post("/auth/reset_password",{password:t},{headers:{Authorization:`Bearer ${null==(e=this.tokenService.getToken())?void 0:e.token}`}})}))}getAppSettings(){return g(this,null,(function*(){return this.get("/auth/app_settings")}))}enableTFA(){return g(this,null,(function*(){var t,i;if(!(null==(t=this.tokenService.getToken())?void 0:t.token))throw new Error("No token in token service.");return this.put("/auth/tfa/enable",{},{headers:{Authorization:`BEARER ${null==(i=this.tokenService.getToken())?void 0:i.token}`}})}))}verifyTFA(t,i){return g(this,null,(function*(){var e,s;if(!(null==(e=this.tokenService.getToken())?void 0:e.token))throw new Error("No token in token service.");return this.post("/auth/tfa/login",{tfa_code:t,scopes:i},{headers:{Authorization:`BEARER ${null==(s=this.tokenService.getToken())?void 0:s.token}`}}).then((t=>this.storeToken(t)))}))}logout(){return g(this,null,(function*(){var t,i,e;if(!(null==(t=this.tokenService.getToken())?void 0:t.token))throw new Error("No token in token service.");return this.post("/me/logout",{refresh_token:null==(i=this.tokenService.getToken("refresh"))?void 0:i.token},{headers:{Authorization:`Bearer ${null==(e=this.tokenService.getToken())?void 0:e.token}`}})}))}storeToken(t){return t.access_token&&this.tokenService.saveToken(t.access_token,"access"),t.refresh_token&&this.tokenService.saveToken(t.refresh_token,"refresh"),t}}const w=/^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-=]*$)/;class b extends class{constructor(t,i,e){this.preffix="identifo_",this.storageType="localStorage",this.access=`${this.preffix}access_token`,this.refresh=`${this.preffix}refresh_token`,this.isAccessible=!0,this.access=i?this.preffix+i:this.access,this.refresh=e?this.preffix+e:this.refresh,this.storageType=t}saveToken(t,i){return!!t&&(window[this.storageType].setItem(this[i],t),!0)}getToken(t){var i;return null!=(i=window[this.storageType].getItem(this[t]))?i:""}deleteToken(t){window[this.storageType].removeItem(this[t])}}{constructor(t,i){super("localStorage",t,i)}}var _=(t,i,e)=>new Promise(((s,o)=>{var r=t=>{try{a(e.next(t))}catch(t){o(t)}},n=t=>{try{a(e.throw(t))}catch(t){o(t)}},a=t=>t.done?s(t.value):Promise.resolve(t.value).then(r,n);a((e=e.apply(t,i)).next())}));class v{constructor(t){this.tokenManager=t||new b}handleVerification(t,i,e){return _(this,null,(function*(){if(!this.tokenManager.isAccessible)return!0;try{return yield this.validateToken(t,i,e),this.saveToken(t),!0}catch(t){return this.removeToken(),Promise.reject(t)}}))}validateToken(t,i,e){return _(this,null,(function*(){var s;if(!t)throw new Error("Empty or invalid token");const o=this.parseJWT(t),r=this.isJWTExpired(o);if((null==(s=o.aud)?void 0:s.includes(i))&&(!e||o.iss===e)&&!r)return Promise.resolve(!0);throw new Error("Empty or invalid token")}))}parseJWT(t){const i=t.split(".")[1];if(!i)return{aud:[],iss:"",exp:10};const e=i.replace(/-/g,"+").replace(/_/g,"/"),s=decodeURIComponent(atob(e).split("").map((t=>`%${`00${t.charCodeAt(0).toString(16)}`.slice(-2)}`)).join(""));return JSON.parse(s)}isJWTExpired(t){const i=(new Date).getTime()/1e3;return!!(t.exp&&i>t.exp)}isAuthenticated(t,i){if(!this.tokenManager.isAccessible)return Promise.resolve(!0);const e=this.tokenManager.getToken("access");return this.validateToken(e,t,i)}saveToken(t,i="access"){return this.tokenManager.saveToken(t,i)}removeToken(t="access"){this.tokenManager.deleteToken(t)}getToken(t="access"){const i=this.tokenManager.getToken(t);return i?{token:i,payload:this.parseJWT(i)}:null}}class y{constructor(t){this.config=t}getUrl(t){var i,e;const s=(null==(i=this.config.scopes)?void 0:i.join())||"",o=encodeURIComponent(null!=(e=this.config.redirectUri)?e:window.location.href),r=`appId=${this.config.appId}&scopes=${s}`,n=`${r}&callbackUrl=${o}`,a=this.config.postLogoutRedirectUri?`&callbackUrl=${encodeURIComponent(this.config.postLogoutRedirectUri)}`:`&callbackUrl=${o}&redirectUri=${this.config.url}/web/login?${encodeURIComponent(r)}`,l={signup:`${this.config.url}/web/register?${n}`,signin:`${this.config.url}/web/login?${n}`,logout:`${this.config.url}/web/logout?${r}${a}`,renew:`${this.config.url}/web/token/renew?${r}&redirectUri=${o}`,default:"default"};return l[t]||l.default}createSignupUrl(){return this.getUrl("signup")}createSigninUrl(){return this.getUrl("signin")}createLogoutUrl(){return this.getUrl("logout")}createRenewSessionUrl(){return this.getUrl("renew")}}var k=Object.defineProperty,$=Object.defineProperties,P=Object.getOwnPropertyDescriptors,z=Object.getOwnPropertySymbols,E=Object.prototype.hasOwnProperty,U=Object.prototype.propertyIsEnumerable,C=(t,i,e)=>i in t?k(t,i,{enumerable:!0,configurable:!0,writable:!0,value:e}):t[i]=e,A=(t,i,e)=>new Promise(((s,o)=>{var r=t=>{try{a(e.next(t))}catch(t){o(t)}},n=t=>{try{a(e.throw(t))}catch(t){o(t)}},a=t=>t.done?s(t.value):Promise.resolve(t.value).then(r,n);a((e=e.apply(t,i)).next())}));class T{constructor(t){var i,e,s,o;this.token=null,this.isAuth=!1,this.config=(s=((t,i)=>{for(var e in i||(i={}))E.call(i,e)&&C(t,e,i[e]);if(z)for(var e of z(i))U.call(i,e)&&C(t,e,i[e]);return t})({},t),o={autoRenew:null==(i=t.autoRenew)||i},$(s,P(o))),this.tokenService=new v(t.tokenManager),this.urlBuilder=new y(this.config),this.api=new x(t,this.tokenService),this.handleToken((null==(e=this.tokenService.getToken())?void 0:e.token)||"","access")}handleToken(t,i){if(t)if("access"===i){const i=this.tokenService.parseJWT(t);this.token={token:t,payload:i},this.isAuth=!0,this.tokenService.saveToken(t)}else this.tokenService.saveToken(t,"refresh")}resetAuthValues(){this.token=null,this.isAuth=!1,this.tokenService.removeToken(),this.tokenService.removeToken("refresh")}signup(){window.location.href=this.urlBuilder.createSignupUrl()}signin(){window.location.href=this.urlBuilder.createSigninUrl()}logout(){this.resetAuthValues(),window.location.href=this.urlBuilder.createLogoutUrl()}handleAuthentication(){return A(this,null,(function*(){const{access:t,refresh:i}=this.getTokenFromUrl();if(!t)return this.resetAuthValues(),Promise.reject();try{return yield this.tokenService.handleVerification(t,this.config.appId,this.config.issuer),this.handleToken(t,"access"),i&&this.handleToken(i,"refresh"),yield Promise.resolve(!0)}catch(t){return this.resetAuthValues(),yield Promise.reject()}finally{window.history.pushState({},document.title,window.location.pathname)}}))}getTokenFromUrl(){const t=new URLSearchParams(window.location.search),i={access:"",refresh:""},e=t.get("token"),s=t.get("refresh_token");return s&&w.test(s)&&(i.refresh=s),e&&w.test(e)&&(i.access=e),i}getToken(){return A(this,null,(function*(){const t=this.tokenService.getToken(),i=this.tokenService.getToken("refresh");if(t){if(this.tokenService.isJWTExpired(t.payload)&&i)try{return yield this.renewSession(),yield Promise.resolve(this.token)}catch(t){throw this.resetAuthValues(),new Error("No token")}return Promise.resolve(t)}return Promise.resolve(null)}))}renewSession(){return A(this,null,(function*(){try{const{access:t,refresh:i}=yield this.renewSessionWithToken();return this.handleToken(t,"access"),this.handleToken(i,"refresh"),yield Promise.resolve(t)}catch(t){return Promise.reject()}}))}renewSessionWithToken(){return A(this,null,(function*(){try{return yield this.api.renewToken().then((t=>({access:t.access_token||"",refresh:t.refresh_token||""})))}catch(t){return Promise.reject(t)}}))}}const I=/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,R=class{constructor(e){t(this,e),this.complete=i(this,"complete",7),this.error=i(this,"error",7),this.route="login",this.theme="light",this.scopes="",this.afterLoginRedirect=t=>{if(this.phone=t.user.phone||"",this.email=t.user.email||"",this.lastResponse=t,t.require_2fa){if(!t.enabled_2fa)return"tfa/setup";if(t.enabled_2fa)return"tfa/verify"}return t.access_token&&t.refresh_token||t.access_token&&!t.refresh_token?"callback":void 0},this.loginCatchRedirect=t=>{if(t.id===r.PleaseEnableTFA)return"tfa/setup";throw t}}processError(t){this.lastError=t,this.error.emit(t)}async signIn(){await this.auth.api.login(this.email,this.password,"",this.scopes.split(",")).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((t=>this.openRoute(t))).catch((t=>this.processError(t)))}async loginWith(t){this.route="loading";const i=this.federatedRedirectUrl||window.location.origin+window.location.pathname;this.auth.api.federatedLogin(t,this.scopes.split(","),i,this.callbackUrl)}async signUp(){this.validateEmail(this.email)&&await this.auth.api.register(this.email,this.password,this.scopes.split(",")).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((t=>this.openRoute(t))).catch((t=>this.processError(t)))}async verifyTFA(){this.auth.api.verifyTFA(this.tfaCode,[]).then((()=>this.openRoute("callback"))).catch((t=>this.processError(t)))}async setupTFA(){if(this.tfaType==a.TFATypeSMS)try{await this.auth.api.updateUser({new_phone:this.phone})}catch(t){return void this.processError(t)}await this.auth.api.enableTFA().then((t=>{t.provisioning_uri||this.openRoute("tfa/verify"),t.provisioning_uri&&(this.provisioningURI=t.provisioning_uri,this.provisioningQR=t.provisioning_qr,this.openRoute("tfa/verify"))}))}restorePassword(){this.auth.api.requestResetPassword(this.email).then((()=>{this.success=!0,this.openRoute("password/forgot/success")})).catch((t=>this.processError(t)))}setNewPassword(){this.token&&this.auth.tokenService.saveToken(this.token,"access"),this.auth.api.resetPassword(this.password).then((()=>{this.success=!0,this.openRoute("login"),this.password=""})).catch((t=>this.processError(t)))}openRoute(t){this.lastError=void 0,this.route=t}usernameChange(t){this.username=t.target.value}passwordChange(t){this.password=t.target.value}emailChange(t){this.email=t.target.value}phoneChange(t){this.phone=t.target.value}tfaCodeChange(t){this.tfaCode=t.target.value}validateEmail(t){return!!I.test(t)||(this.processError({detailedMessage:"Email address is not valid",name:"Validation error",message:"Email address is not valid"}),!1)}renderRoute(t){var i,o,r,n,l,h;switch(t){case"login":return e("div",{class:"login-form"},!this.registrationForbidden&&e("p",{class:"login-form__register-text"},"Don't have an account?",e("a",{onClick:()=>this.openRoute("register"),class:"login-form__register-link"}," ","Sign Up")),e("input",{type:"text",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingInput",value:this.email,placeholder:"Email",onInput:t=>this.emailChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.email||!this.password)&&this.signIn()}),e("input",{type:"password",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingPassword",value:this.password,placeholder:"Password",onInput:t=>this.passwordChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.email||!this.password)&&this.signIn()}),!!this.lastError&&e("div",{class:"error",role:"alert"},null===(i=this.lastError)||void 0===i?void 0:i.detailedMessage),e("div",{class:"login-form__buttons "+(this.lastError?"login-form__buttons_mt-32":"")},e("button",{onClick:()=>this.signIn(),class:"primary-button",disabled:!this.email||!this.password},"Login"),e("a",{onClick:()=>this.openRoute("password/forgot"),class:"login-form__forgot-pass"},"Forgot password")),this.federatedProviders.length>0&&e("div",{class:"social-buttons"},e("p",{class:"social-buttons__text"},"or continue with"),e("div",{class:"social-buttons__social-medias"},this.federatedProviders.indexOf("apple")>-1&&e("div",{class:"social-buttons__media",onClick:()=>this.loginWith("apple")},e("img",{src:s("assets/images/apple.svg"),class:"social-buttons__image",alt:"login via apple"})),this.federatedProviders.indexOf("google")>-1&&e("div",{class:"social-buttons__media",onClick:()=>this.loginWith("google")},e("img",{src:s("assets/images/google.svg"),class:"social-buttons__image",alt:"login via google"})),this.federatedProviders.indexOf("facebook")>-1&&e("div",{class:"social-buttons__media",onClick:()=>this.loginWith("facebook")},e("img",{src:s("assets/images/fb.svg"),class:"social-buttons__image",alt:"login via facebook"})))));case"register":return e("div",{class:"register-form"},e("input",{type:"text",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingInput",value:this.email,placeholder:"Email",onInput:t=>this.emailChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.password||!this.email)&&this.signUp()}),e("input",{type:"password",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingPassword",value:this.password,placeholder:"Password",onInput:t=>this.passwordChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.password||!this.email)&&this.signUp()}),!!this.lastError&&e("div",{class:"error",role:"alert"},null===(o=this.lastError)||void 0===o?void 0:o.detailedMessage),e("div",{class:"register-form__buttons "+(this.lastError?"register-form__buttons_mt-32":"")},e("button",{onClick:()=>this.signUp(),class:"primary-button",disabled:!this.email||!this.password},"Continue"),e("a",{onClick:()=>this.openRoute("login"),class:"register-form__log-in"},"Go back to login")));case"otp/login":return e("div",{class:"otp-login"},"registrationForbidden",!1===this.registrationForbidden&&e("p",{class:"otp-login__register-text"},"Don't have an account?",e("a",{onClick:()=>this.openRoute("register"),class:"login-form__register-link"}," ","Sign Up")),e("input",{type:"phone",class:"form-control",id:"floatingInput",value:this.phone,placeholder:"Phone number",onInput:t=>this.phoneChange(t)}),e("button",{onClick:()=>this.openRoute("tfa/verify"),class:"primary-button",disabled:!this.phone},"Continue"),this.federatedProviders.length>0&&e("div",{class:"social-buttons"},e("p",{class:"social-buttons__text"},"or continue with"),e("div",{class:"social-buttons__social-medias"},this.federatedProviders.indexOf("apple")>-1&&e("div",{class:"social-buttons__media",onClick:()=>this.loginWith("apple")},e("img",{src:s("assets/images/apple.svg"),class:"social-buttons__image",alt:"login via apple"})),this.federatedProviders.indexOf("google")>-1&&e("div",{class:"social-buttons__media",onClick:()=>this.loginWith("google")},e("img",{src:s("assets/images/google.svg"),class:"social-buttons__image",alt:"login via google"})),this.federatedProviders.indexOf("facebook")>-1&&e("div",{class:"social-buttons__media",onClick:()=>this.loginWith("facebook")},e("img",{src:s("assets/images/fb.svg"),class:"social-buttons__image",alt:"login via facebook"})))));case"tfa/setup":return e("div",{class:"tfa-setup"},e("p",{class:"tfa-setup__text"},"Protect your account with 2-step verification"),this.tfaType===a.TFATypeApp&&e("div",{class:"info-card"},e("div",{class:"info-card__controls"},e("p",{class:"info-card__title"},"Authenticator app"),e("button",{type:"button",class:"info-card__button",onClick:()=>this.setupTFA()},"Setup")),e("p",{class:"info-card__text"},"Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone.")),this.tfaType===a.TFATypeEmail&&e("div",{class:"info-card"},e("div",{class:"info-card__controls"},e("p",{class:"info-card__title"},"Email"),e("button",{type:"button",class:"info-card__button",onClick:()=>this.setupTFA()},"Setup")),e("p",{class:"info-card__subtitle"},this.email),e("p",{class:"info-card__text"}," Use email as 2fa, please check your email, we will send confirmation code to this email.")),this.tfaType===a.TFATypeSMS&&e("div",{class:"tfa-setup__form"},e("p",{class:"tfa-setup__subtitle"}," Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone"),e("input",{type:"phone",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingInput",value:this.phone,placeholder:"Phone",onInput:t=>this.phoneChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.phone)&&this.setupTFA()}),!!this.lastError&&e("div",{class:"error",role:"alert"},null===(r=this.lastError)||void 0===r?void 0:r.detailedMessage),e("button",{onClick:()=>this.setupTFA(),class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.phone},"Setup phone")));case"tfa/verify":return e("div",{class:"tfa-verify"},!(this.tfaType!==a.TFATypeApp)&&e("div",{class:"tfa-verify__title-wrapper"},e("h2",{class:this.provisioningURI?"tfa-verify__title":"tfa-verify__title_mb-40"},this.provisioningURI?"Please scan QR-code with the app":"Use GoogleAuth as 2fa"),!!this.provisioningURI&&e("img",{src:`data:image/png;base64, ${this.provisioningQR}`,alt:this.provisioningURI,class:"tfa-verify__qr-code"})),!(this.tfaType!==a.TFATypeSMS)&&e("div",{class:"tfa-verify__title-wrapper"},e("h2",{class:"tfa-verify__title"},"Enter the code sent to your phone number"),e("p",{class:"tfa-verify__subtitle"},"The code has been sent to ",this.phone)),!(this.tfaType!==a.TFATypeEmail)&&e("div",{class:"tfa-verify__title-wrapper"},e("h2",{class:"tfa-verify__title"},"Enter the code sent to your email address"),e("p",{class:"tfa-verify__subtitle"},"The email has been sent to ",this.email)),e("input",{type:"text",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingCode",value:this.tfaCode,placeholder:"Verify code",onInput:t=>this.tfaCodeChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.tfaCode)&&this.verifyTFA()}),!!this.lastError&&e("div",{class:"error",role:"alert"},null===(n=this.lastError)||void 0===n?void 0:n.detailedMessage),e("button",{type:"button",class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.tfaCode,onClick:()=>this.verifyTFA()},"Confirm"));case"password/forgot":return e("div",{class:"forgot-password"},e("h2",{class:"forgot-password__title"},"Enter the email you gave when you registered"),e("p",{class:"forgot-password__subtitle"},"We will send you a link to create a new password on email"),e("input",{type:"email",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingEmail",value:this.email,placeholder:"Email",onInput:t=>this.emailChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.email)&&this.restorePassword()}),!!this.lastError&&e("div",{class:"error",role:"alert"},null===(l=this.lastError)||void 0===l?void 0:l.detailedMessage),e("button",{type:"button",class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.email,onClick:()=>this.restorePassword()},"Send the link"));case"password/forgot/success":return e("div",{class:"forgot-password-success"},"dark"===this.theme&&e("img",{src:s("./assets/images/email-dark.svg"),alt:"email",class:"forgot-password-success__image"}),"light"===this.theme&&e("img",{src:s("./assets/images/email.svg"),alt:"email",class:"forgot-password-success__image"}),e("p",{class:"forgot-password-success__text"},"We sent you an email with a link to create a new password"));case"password/reset":return e("div",{class:"reset-password"},e("h2",{class:"reset-password__title"},"Set up a new password to log in to the website"),e("p",{class:"reset-password__subtitle"},"Memorize your password and do not give it to anyone."),e("input",{type:"password",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingPassword",value:this.password,placeholder:"Password",onInput:t=>this.passwordChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.password)&&this.setNewPassword()}),!!this.lastError&&e("div",{class:"error",role:"alert"},null===(h=this.lastError)||void 0===h?void 0:h.detailedMessage),e("button",{type:"button",class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.password,onClick:()=>this.setNewPassword()},"Save password"));case"error":return e("div",{class:"error-view"},e("div",{class:"error-view__message"},this.lastError.message),e("div",{class:"error-view__details"},this.lastError.detailedMessage));case"callback":return e("div",{class:"error-view"},e("div",null,"Success"),this.debug&&e("div",null,e("div",null,"Access token: ",this.lastResponse.access_token),e("div",null,"Refresh token: ",this.lastResponse.refresh_token),e("div",null,"User: ",JSON.stringify(this.lastResponse.user))));case"loading":return e("div",{class:"error-view"},e("div",null,"Loading ..."))}}async componentWillLoad(){const t=this.postLogoutRedirectUri||window.location.origin+window.location.pathname;this.auth=new T({appId:this.appId,url:this.url,postLogoutRedirectUri:t});try{const t=await this.auth.api.getAppSettings();this.registrationForbidden=t.registrationForbidden,this.tfaType=t.tfaType,this.federatedProviders=t.federatedProviders}catch(t){this.route="error",this.lastError=t}const i=new URL(window.location.href);if(i.searchParams.get("provider")&&i.searchParams.get("state")){const t=new URL(window.location.href),e=new URLSearchParams,s=i.searchParams.get("appId");e.set("appId",s),window.history.replaceState({},document.title,`${t.pathname}?${e.toString()}`),this.route="loading",this.auth.api.federatedLoginComplete(t.searchParams).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((t=>this.openRoute(t))).catch((t=>this.processError(t)))}}componentWillRender(){if("callback"===this.route){const t=new URL(window.location.href);t.searchParams.set("callbackUrl",this.lastResponse.callbackUrl),window.history.replaceState({},document.title,`${t.pathname}?${t.searchParams.toString()}`),this.complete.emit(this.lastResponse)}"logout"===this.route&&this.complete.emit()}render(){return e(o,null,e("div",{class:{wrapper:"light"===this.theme,"wrapper-dark":"dark"===this.theme}},this.renderRoute(this.route)),e("div",{class:"error-view"},this.debug&&e("div",null,e("br",null),this.appId)))}static get assetsDirs(){return["assets"]}};R.style='.wrapper,.wrapper-dark{--content-width:416px}.wrapper{--main-background:#f7f7f7;--blue-text:#6163f6;--field-background:#fff;--gray-line:#e0e0e0;--social-button:#1b1b1b;--text:#1b1b1b;--upload-photo:#e0e0e0;--content-width:416px}.wrapper-dark{--main-background:#1b1b1b;--blue-text:#8b8dfa;--field-background:#423f3f;--gray-line:#423f3f;--social-button:#423f3f;--text:#fff;--upload-photo:#423f3f;--content-width:416px}*{margin:0;padding:0;-webkit-box-sizing:border-box;box-sizing:border-box;font-family:inherit}.wrapper,.wrapper-dark{display:-ms-flexbox;display:flex;-ms-flex-pack:center;justify-content:center;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.social-buttons{width:100%;position:relative}.social-buttons__text{font-size:14px;line-height:21px;color:#828282;padding:4px 8px;margin-bottom:39px;text-align:center;position:static}.social-buttons__text::before{content:"";position:absolute;height:1px;width:36%;left:0;top:14px;background-color:var(--gray-line)}.social-buttons__text::after{content:"";position:absolute;height:1px;width:36%;right:0;top:14px;background-color:var(--gray-line)}.social-buttons__social-medias{display:-ms-flexbox;display:flex;-ms-flex-pack:center;justify-content:center;-ms-flex-align:center;align-items:center}.social-buttons__media{width:56px;height:56px;border-radius:50%;background-color:var(--social-button);display:-ms-flexbox;display:flex;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;cursor:pointer}.social-buttons__media:not(:last-of-type){margin-right:24px}@media (max-width: 599px){.social-buttons__media{width:44px;height:44px}.social-buttons__text{margin-bottom:36px}.social-buttons__text::before{width:26%}.social-buttons__text::after{width:26%}.social-buttons__image{width:16px;height:16px}}.primary-button{background-color:#6163f6;border:none;outline:none;display:-ms-flexbox;display:flex;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;width:192px;height:64px;border-radius:8px;cursor:pointer;color:#fff;font-size:18px;line-height:26px;-webkit-transition:all 0.4s;transition:all 0.4s}.primary-button:active{-webkit-transform:translateY(-4px);transform:translateY(-4px)}.primary-button:hover{opacity:0.8}.primary-button:disabled{cursor:initial;opacity:0.3}@media (max-width: 599px){.primary-button{width:100%}}.info-card{border:1px solid var(--gray-line);border-radius:8px;padding:24px}.info-card__controls{display:-ms-flexbox;display:flex;-ms-flex-pack:justify;justify-content:space-between}.info-card__title{color:var(--text);font-size:18px;line-height:26px;font-weight:700}.info-card__button{color:var(--blue-text);background:none;border:none;cursor:pointer;font-size:18px;line-height:26px}.info-card__text{color:#828282;font-size:16px;line-height:24px;margin-top:8px}.info-card__subtitle{color:var(--text);font-size:16px;line-height:24px;margin:4px 0 12px}@media (max-width: 599px){.info-card__text{font-size:14px;line-height:21px}}.form-control{width:100%;max-width:var(--content-width);height:72px;background-color:var(--field-background);-webkit-box-shadow:0px 11px 15px rgba(0, 0, 0, 0.04);box-shadow:0px 11px 15px rgba(0, 0, 0, 0.04);border-radius:8px;border:none;outline:none;font-size:18px;line-height:26px;color:var(--text);padding:23px 24px}.form-control::-webkit-inner-spin-button{-webkit-appearance:none;margin:0}.form-control::-webkit-input-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control::-moz-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control:-ms-input-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control::-ms-input-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control::placeholder{font-size:18px;line-height:26px;color:#828282}.form-control-danger{border:1px solid #F66161}@media (max-width: 599px){.form-control{height:64px}}.upload-photo{display:-ms-flexbox;display:flex;-ms-flex-pack:center;justify-content:center;-ms-flex-align:center;align-items:center;margin-bottom:48px}.upload-photo__field{display:none}.upload-photo__label{cursor:pointer;color:var(--blue-text);font-size:16px;line-height:24px}.upload-photo__label:first-of-type{margin-right:16px}.upload-photo__avatar{height:64px;width:64px;border-radius:50%;background-color:var(--upload-photo)}@media (max-width: 599px){.upload-photo{margin-bottom:32px}}.error{min-height:21px;width:100%;font-size:14px;line-height:21px;color:#FF5160}.login-form{display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;width:var(--content-width)}.login-form__register-text{margin-bottom:32px;font-weight:400;font-size:14px;line-height:24px;color:#828282}.login-form__register-link{color:var(--blue-text);cursor:pointer}.login-form .form-control:first-of-type{margin-bottom:32px}.login-form__buttons{margin-top:48px;display:-ms-flexbox;display:flex;width:100%;-ms-flex-align:center;align-items:center;margin-bottom:36px}.login-form__buttons_mt-32{margin-top:32px}.login-form__forgot-pass{color:var(--blue-text);font-size:16px;line-height:24px;cursor:pointer}.login-form .primary-button{margin-right:32px}.login-form .error{margin-top:12px}@media (max-width: 599px){.login-form{width:100%;max-width:var(--content-width);padding:0 24px}.login-form__register-text{font-size:16px}.login-form .form-control:first-of-type{margin-bottom:24px}.login-form__buttons{margin-top:32px;-ms-flex-direction:column;flex-direction:column}.login-form .primary-button{margin-right:0;margin-bottom:36px}}.register-form{width:var(--content-width);padding:64px 0 44px}.register-form .form-control:not(:last-of-type){margin-bottom:32px}.register-form__buttons{display:-ms-flexbox;display:flex;width:100%;-ms-flex-align:center;align-items:center;margin-top:48px}.register-form__buttons_mt-32{margin-top:32px}.register-form__log-in{color:var(--blue-text);font-size:16px;line-height:24px;cursor:pointer}.register-form .primary-button{margin-right:32px}.register-form .error{margin-top:12px}@media (max-width: 599px){.register-form{width:100%;max-width:var(--content-width);padding:48px 24px 32px}.register-form .form-control{margin-bottom:24px}.register-form .primary-button{margin-right:0;margin-bottom:36px}.register-form__buttons{-ms-flex-direction:column;flex-direction:column;margin-top:32px}}.tfa-setup{padding:48px 0 80px;width:var(--content-width);display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.tfa-setup__text{font-size:24px;line-height:32px;color:var(--text);text-align:center;max-width:260px;width:100%}.tfa-setup .error{margin-top:12px}.tfa-setup__form{width:100%}.tfa-setup__form .primary-button{margin:48px auto 0}.tfa-setup__form .primary-button-mt-32{margin-top:32px}.tfa-setup__subtitle{text-align:center;font-size:16px;line-height:24px;color:#828282;max-width:270px;text-align:center;margin:16px auto 48px}.tfa-setup .info-card{margin-top:48px}@media (max-width: 599px){.tfa-setup{width:100%;max-width:var(--content-width);padding:38px 24px 36px}}.tfa-verify{padding:48px 0 52px;width:var(--content-width);display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.tfa-verify .error{margin-top:12px}.tfa-verify__title-wrapper{width:100%;display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.tfa-verify__title,.tfa-verify__title_mb-40{color:var(--text);font-size:24px;line-height:32px;max-width:280px;text-align:center;font-weight:400;margin-bottom:16px}.tfa-verify__title_mb-40{margin-bottom:40px}.tfa-verify__app-button{font-size:18px;line-height:26px;background:none;border:none;color:var(--blue-text);margin-bottom:40px}.tfa-verify .primary-button{margin-top:48px}.tfa-verify .primary-button-mt-32{margin-top:32px}.tfa-verify__back{font-size:16px;line-height:24px;color:var(--blue-text)}.tfa-verify__subtitle{font-size:16px;line-height:24px;margin-bottom:48px;color:#828282;max-width:189px;text-align:center}.tfa-verify__qr-code{width:160px;height:160px;margin-bottom:64px}@media (max-width: 599px){.tfa-verify{padding:102px 24px 41px;width:100%;max-width:var(--content-width)}.tfa-verify .primary-button{margin-top:32px}.tfa-verify__qr-code{margin-bottom:48px}}.otp-login{display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;padding:44px 24px 102px;width:100%;max-width:var(--content-width)}.otp-login__register-text{margin-bottom:32px;font-weight:400;font-size:14px;line-height:24px;color:#828282}.otp-login .form-control{margin-bottom:48px}.otp-login .primary-button{margin-bottom:36px}.error-view{width:100%;max-width:var(--content-width);padding:0 24px;display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center;word-break:break-all}.error-view__message{color:var(--text);font-size:24px;line-height:32px;text-align:center;margin-bottom:16px}.error-view__details{color:#828282;font-size:16px;line-height:24px;text-align:center}.error-view .primary-button{margin-top:64px}.forgot-password{width:var(--content-width)}.forgot-password__title{text-align:center;margin-bottom:16px;font-size:24px;line-height:32px;font-weight:400;color:var(--text)}.forgot-password__subtitle{text-align:center;font-size:16px;line-height:24px;color:#828282;max-width:189px;text-align:center;margin:0 auto 48px}.forgot-password .error{margin-top:12px}.forgot-password .primary-button{margin:48px auto 0}.forgot-password .primary-button-mt-32{margin-top:32px}@media (max-width: 599px){.forgot-password{width:100%;max-width:var(--content-width);padding:0 24px}.forgot-password .primary-button{margin-top:32px}}.forgot-password-success{display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center;padding:0 16px}.forgot-password-success__text{width:100%;max-width:367px;font-size:24px;line-height:32px;text-align:center;color:var(--text)}.forgot-password-success__image{margin-bottom:56px}.reset-password{width:var(--content-width)}.reset-password__title{text-align:center;font-size:24px;line-height:32px;font-weight:400;color:var(--text);max-width:270px;margin:0 auto 16px}.reset-password__subtitle{text-align:center;font-size:16px;line-height:24px;color:#828282;max-width:189px;text-align:center;margin:0 auto 48px}.reset-password .error{margin-top:12px}.reset-password .primary-button{margin:48px auto 0}.reset-password .primary-button-mt-32{margin-top:32px}@media (max-width: 599px){.reset-password{width:100%;max-width:var(--content-width);padding:0 24px}.reset-password .primary-button{margin-top:32px}}';export{R as identifo_form}
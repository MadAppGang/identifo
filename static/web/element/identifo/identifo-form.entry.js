import{r as t,c as s,h as i,g as e,H as r}from"./index-5069a111.js";var o,n,a,l;(n=o||(o={})).PleaseEnableTFA="error.api.request.2fa.please_enable",n.NetworkError="error.network",(l=a||(a={})).TFATypeApp="app",l.TFATypeSMS="sms",l.TFATypeEmail="email";class h extends Error{constructor(t){super((null==t?void 0:t.message)||"Unknown API error"),this.detailedMessage=null==t?void 0:t.detailed_message,this.id=null==t?void 0:t.id,this.status=null==t?void 0:t.status}}var c=Object.defineProperty,u=Object.getOwnPropertySymbols,d=Object.prototype.hasOwnProperty,p=Object.prototype.propertyIsEnumerable,f=(t,s,i)=>s in t?c(t,s,{enumerable:!0,configurable:!0,writable:!0,value:i}):t[s]=i,g=(t,s)=>{for(var i in s||(s={}))d.call(s,i)&&f(t,i,s[i]);if(u)for(var i of u(s))p.call(s,i)&&f(t,i,s[i]);return t},w=(t,s,i)=>new Promise(((e,r)=>{var o=t=>{try{a(i.next(t))}catch(t){r(t)}},n=t=>{try{a(i.throw(t))}catch(t){r(t)}},a=t=>t.done?e(t.value):Promise.resolve(t.value).then(o,n);a((i=i.apply(t,s)).next())}));class v{constructor(t,s){this.config=t,this.tokenService=s,this.defaultHeaders={"X-Identifo-Clientid":"",Accept:"application/json","Content-Type":"application/json"},this.catchNetworkErrorHandler=t=>{if("Network Error"===t.message||"Failed to fetch"===t.message||"Preflight response is not successful"===t.message||t.message.indexOf("is not allowed by Access-Control-Allow-Origin")>-1)throw console.error(t.message),new h({id:o.NetworkError,status:0,message:"Configuration error",detailed_message:`Please check Identifo URL and add "${window.location.protocol}//${window.location.host}" to "REDIRECT URLS" in Identifo app settings.`});throw t},this.checkStatusCodeAndGetJSON=t=>w(this,null,(function*(){if(!t.ok){const s=yield t.json();throw new h(null==s?void 0:s.error)}return t.json()})),this.baseUrl=t.url.replace(/\/$/,""),this.defaultHeaders["X-Identifo-Clientid"]=t.appId,this.appId=t.appId}get(t,s){return this.send(t,g({method:"GET"},s))}put(t,s,i){return this.send(t,g({method:"PUT",body:JSON.stringify(s)},i))}post(t,s,i){return this.send(t,g({method:"POST",body:JSON.stringify(s)},i))}send(t,s){const i=g({},s);return i.credentials="include",i.headers=g(g({},i.headers),this.defaultHeaders),fetch(`${this.baseUrl}${t}`,i).catch(this.catchNetworkErrorHandler).then(this.checkStatusCodeAndGetJSON).then((t=>t))}getUser(){return w(this,null,(function*(){var t,s;if(!(null==(t=this.tokenService.getToken())?void 0:t.token))throw new Error("No token in token service.");return this.get("/me",{headers:{Authorization:`Bearer ${null==(s=this.tokenService.getToken())?void 0:s.token}`}})}))}renewToken(){return w(this,null,(function*(){var t,s;if(!(null==(t=this.tokenService.getToken("refresh"))?void 0:t.token))throw new Error("No token in token service.");return this.post("/auth/token",{scopes:this.config.scopes},{headers:{Authorization:`Bearer ${null==(s=this.tokenService.getToken("refresh"))?void 0:s.token}`}}).then((t=>this.storeToken(t)))}))}updateUser(t){return w(this,null,(function*(){var s,i;if(!(null==(s=this.tokenService.getToken())?void 0:s.token))throw new Error("No token in token service.");return this.put("/me",t,{headers:{Authorization:`Bearer ${null==(i=this.tokenService.getToken("access"))?void 0:i.token}`}})}))}login(t,s,i,e){return w(this,null,(function*(){return this.post("/auth/login",{email:t,password:s,device_token:i,scopes:e}).then((t=>this.storeToken(t)))}))}federatedLogin(t,s,i,e){return w(this,arguments,(function*(t,s,i,e,r={width:600,height:800,popUp:!1}){var o=document.createElement("form");o.style.display="none",r.popUp&&(o.target="TargetWindow"),o.method="POST";const n=new URLSearchParams;if(n.set("appId",this.config.appId),n.set("provider",t),n.set("scopes",s.join(",")),n.set("redirectUrl",i),e&&n.set("callbackUrl",e),o.action=`${this.baseUrl}/auth/federated?${n.toString()}`,document.body.appendChild(o),r.popUp){const t=window.screenX+window.outerWidth/2-(r.width||600)/2,s=window.screenY+window.outerHeight/2-(r.height||800)/2;window.open("","TargetWindow",`status=0,title=0,height=${r.height},width=${r.width},top=${s},left=${t},scrollbars=1`)&&o.submit()}else window.location.assign(`${this.baseUrl}/auth/federated?${n.toString()}`)}))}federatedLoginComplete(t){return w(this,null,(function*(){return this.get(`/auth/federated/complete?${t.toString()}`).then((t=>this.storeToken(t)))}))}register(t,s,i){return w(this,null,(function*(){return this.post("/auth/register",{email:t,password:s,scopes:i}).then((t=>this.storeToken(t)))}))}requestResetPassword(t){return w(this,null,(function*(){return this.post("/auth/request_reset_password",{email:t})}))}resetPassword(t){return w(this,null,(function*(){var s,i;if(!(null==(s=this.tokenService.getToken())?void 0:s.token))throw new Error("No token in token service.");return this.post("/auth/reset_password",{password:t},{headers:{Authorization:`Bearer ${null==(i=this.tokenService.getToken())?void 0:i.token}`}})}))}getAppSettings(){return w(this,null,(function*(){return this.get("/auth/app_settings")}))}enableTFA(){return w(this,null,(function*(){var t,s;if(!(null==(t=this.tokenService.getToken())?void 0:t.token))throw new Error("No token in token service.");return this.put("/auth/tfa/enable",{},{headers:{Authorization:`BEARER ${null==(s=this.tokenService.getToken())?void 0:s.token}`}})}))}verifyTFA(t,s){return w(this,null,(function*(){var i,e;if(!(null==(i=this.tokenService.getToken())?void 0:i.token))throw new Error("No token in token service.");return this.post("/auth/tfa/login",{tfa_code:t,scopes:s},{headers:{Authorization:`BEARER ${null==(e=this.tokenService.getToken())?void 0:e.token}`}}).then((t=>this.storeToken(t)))}))}logout(){return w(this,null,(function*(){var t,s,i;if(!(null==(t=this.tokenService.getToken())?void 0:t.token))throw new Error("No token in token service.");return this.post("/me/logout",{refresh_token:null==(s=this.tokenService.getToken("refresh"))?void 0:s.token},{headers:{Authorization:`Bearer ${null==(i=this.tokenService.getToken())?void 0:i.token}`}})}))}storeToken(t){return t.access_token&&this.tokenService.saveToken(t.access_token,"access"),t.refresh_token&&this.tokenService.saveToken(t.refresh_token,"refresh"),t}}const m=/^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-=]*$)/;class _ extends class{constructor(t,s,i){this.preffix="identifo_",this.storageType="localStorage",this.access=`${this.preffix}access_token`,this.refresh=`${this.preffix}refresh_token`,this.isAccessible=!0,this.access=s?this.preffix+s:this.access,this.refresh=i?this.preffix+i:this.refresh,this.storageType=t}saveToken(t,s){return!!t&&(window[this.storageType].setItem(this[s],t),!0)}getToken(t){var s;return null!=(s=window[this.storageType].getItem(this[t]))?s:""}deleteToken(t){window[this.storageType].removeItem(this[t])}}{constructor(t,s){super("localStorage",t,s)}}var b=(t,s,i)=>new Promise(((e,r)=>{var o=t=>{try{a(i.next(t))}catch(t){r(t)}},n=t=>{try{a(i.throw(t))}catch(t){r(t)}},a=t=>t.done?e(t.value):Promise.resolve(t.value).then(o,n);a((i=i.apply(t,s)).next())}));class y{constructor(t){this.tokenManager=t||new _}handleVerification(t,s,i){return b(this,null,(function*(){if(!this.tokenManager.isAccessible)return!0;try{return yield this.validateToken(t,s,i),this.saveToken(t),!0}catch(t){return this.removeToken(),Promise.reject(t)}}))}validateToken(t,s,i){return b(this,null,(function*(){var e;if(!t)throw new Error("Empty or invalid token");const r=this.parseJWT(t),o=this.isJWTExpired(r);if((null==(e=r.aud)?void 0:e.includes(s))&&(!i||r.iss===i)&&!o)return Promise.resolve(!0);throw new Error("Empty or invalid token")}))}parseJWT(t){const s=t.split(".")[1];if(!s)return{aud:[],iss:"",exp:10};const i=s.replace(/-/g,"+").replace(/_/g,"/"),e=decodeURIComponent(atob(i).split("").map((t=>`%${`00${t.charCodeAt(0).toString(16)}`.slice(-2)}`)).join(""));return JSON.parse(e)}isJWTExpired(t){const s=(new Date).getTime()/1e3;return!!(t.exp&&s>t.exp)}isAuthenticated(t,s){if(!this.tokenManager.isAccessible)return Promise.resolve(!0);const i=this.tokenManager.getToken("access");return this.validateToken(i,t,s)}saveToken(t,s="access"){return this.tokenManager.saveToken(t,s)}removeToken(t="access"){this.tokenManager.deleteToken(t)}getToken(t="access"){const s=this.tokenManager.getToken(t);return s?{token:s,payload:this.parseJWT(s)}:null}}class k{constructor(t){this.config=t}getUrl(t){var s,i;const e=(null==(s=this.config.scopes)?void 0:s.join())||"",r=encodeURIComponent(null!=(i=this.config.redirectUri)?i:window.location.href),o=`appId=${this.config.appId}&scopes=${e}`,n=`${o}&callbackUrl=${r}`,a=this.config.postLogoutRedirectUri?`&callbackUrl=${encodeURIComponent(this.config.postLogoutRedirectUri)}`:`&callbackUrl=${r}&redirectUri=${this.config.url}/web/login?${encodeURIComponent(o)}`,l={signup:`${this.config.url}/web/register?${n}`,signin:`${this.config.url}/web/login?${n}`,logout:`${this.config.url}/web/logout?${o}${a}`,renew:`${this.config.url}/web/token/renew?${o}&redirectUri=${r}`,default:"default"};return l[t]||l.default}createSignupUrl(){return this.getUrl("signup")}createSigninUrl(){return this.getUrl("signin")}createLogoutUrl(){return this.getUrl("logout")}createRenewSessionUrl(){return this.getUrl("renew")}}var $=Object.defineProperty,P=Object.defineProperties,E=Object.getOwnPropertyDescriptors,U=Object.getOwnPropertySymbols,C=Object.prototype.hasOwnProperty,A=Object.prototype.propertyIsEnumerable,T=(t,s,i)=>s in t?$(t,s,{enumerable:!0,configurable:!0,writable:!0,value:i}):t[s]=i,I=(t,s,i)=>new Promise(((e,r)=>{var o=t=>{try{a(i.next(t))}catch(t){r(t)}},n=t=>{try{a(i.throw(t))}catch(t){r(t)}},a=t=>t.done?e(t.value):Promise.resolve(t.value).then(o,n);a((i=i.apply(t,s)).next())}));class R{constructor(t){var s,i,e,r;this.token=null,this.isAuth=!1,this.config=(e=((t,s)=>{for(var i in s||(s={}))C.call(s,i)&&T(t,i,s[i]);if(U)for(var i of U(s))A.call(s,i)&&T(t,i,s[i]);return t})({},t),r={autoRenew:null==(s=t.autoRenew)||s},P(e,E(r))),this.tokenService=new y(t.tokenManager),this.urlBuilder=new k(this.config),this.api=new v(t,this.tokenService),this.handleToken((null==(i=this.tokenService.getToken())?void 0:i.token)||"","access")}handleToken(t,s){if(t)if("access"===s){const s=this.tokenService.parseJWT(t);this.token={token:t,payload:s},this.isAuth=!0,this.tokenService.saveToken(t)}else this.tokenService.saveToken(t,"refresh")}resetAuthValues(){this.token=null,this.isAuth=!1,this.tokenService.removeToken(),this.tokenService.removeToken("refresh")}signup(){window.location.href=this.urlBuilder.createSignupUrl()}signin(){window.location.href=this.urlBuilder.createSigninUrl()}logout(){this.resetAuthValues(),window.location.href=this.urlBuilder.createLogoutUrl()}handleAuthentication(){return I(this,null,(function*(){const{access:t,refresh:s}=this.getTokenFromUrl();if(!t)return this.resetAuthValues(),Promise.reject();try{return yield this.tokenService.handleVerification(t,this.config.appId,this.config.issuer),this.handleToken(t,"access"),s&&this.handleToken(s,"refresh"),yield Promise.resolve(!0)}catch(t){return this.resetAuthValues(),yield Promise.reject()}finally{window.history.pushState({},document.title,window.location.pathname)}}))}getTokenFromUrl(){const t=new URLSearchParams(window.location.search),s={access:"",refresh:""},i=t.get("token"),e=t.get("refresh_token");return e&&m.test(e)&&(s.refresh=e),i&&m.test(i)&&(s.access=i),s}getToken(){return I(this,null,(function*(){const t=this.tokenService.getToken(),s=this.tokenService.getToken("refresh");if(t){if(this.tokenService.isJWTExpired(t.payload)&&s)try{return yield this.renewSession(),yield Promise.resolve(this.token)}catch(t){throw this.resetAuthValues(),new Error("No token")}return Promise.resolve(t)}return Promise.resolve(null)}))}renewSession(){return I(this,null,(function*(){try{const{access:t,refresh:s}=yield this.renewSessionWithToken();return this.handleToken(t,"access"),this.handleToken(s,"refresh"),yield Promise.resolve(t)}catch(t){return Promise.reject()}}))}renewSessionWithToken(){return I(this,null,(function*(){try{return yield this.api.renewToken().then((t=>({access:t.access_token||"",refresh:t.refresh_token||""})))}catch(t){return Promise.reject(t)}}))}}const S=/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,x=class{constructor(i){t(this,i),this.complete=s(this,"complete",7),this.error=s(this,"error",7),this.route="login",this.theme="light",this.scopes="",this.afterLoginRedirect=t=>{if(this.phone=t.user.phone||"",this.email=t.user.email||"",this.lastResponse=t,t.require_2fa){if(!t.enabled_2fa)return"tfa/setup";if(t.enabled_2fa)return"tfa/verify"}return t.access_token&&t.refresh_token||t.access_token&&!t.refresh_token?"callback":void 0},this.loginCatchRedirect=t=>{if(t.id===o.PleaseEnableTFA)return"tfa/setup";throw t}}processError(t){this.lastError=t,this.error.emit(t)}async signIn(){await this.auth.api.login(this.email,this.password,"",this.scopes.split(",")).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((t=>this.openRoute(t))).catch((t=>this.processError(t)))}async loginWith(t){this.route="loading";const s=this.federatedRedirectUrl||window.location.origin+window.location.pathname;this.auth.api.federatedLogin(t,this.scopes.split(","),s,this.callbackUrl)}async signUp(){this.validateEmail(this.email)&&await this.auth.api.register(this.email,this.password,this.scopes.split(",")).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((t=>this.openRoute(t))).catch((t=>this.processError(t)))}async verifyTFA(){this.auth.api.verifyTFA(this.tfaCode,[]).then((()=>this.openRoute("callback"))).catch((t=>this.processError(t)))}async setupTFA(){if(this.tfaType==a.TFATypeSMS)try{await this.auth.api.updateUser({new_phone:this.phone})}catch(t){return void this.processError(t)}await this.auth.api.enableTFA().then((t=>{t.provisioning_uri||this.openRoute("tfa/verify"),t.provisioning_uri&&(this.provisioningURI=t.provisioning_uri,this.provisioningQR=t.provisioning_qr,this.openRoute("tfa/verify"))}))}restorePassword(){this.auth.api.requestResetPassword(this.email).then((()=>{this.success=!0,this.openRoute("password/forgot/success")})).catch((t=>this.processError(t)))}setNewPassword(){this.token&&this.auth.tokenService.saveToken(this.token,"access"),this.auth.api.resetPassword(this.password).then((()=>{this.success=!0,this.openRoute("login"),this.password=""})).catch((t=>this.processError(t)))}openRoute(t){this.lastError=void 0,this.route=t}usernameChange(t){this.username=t.target.value}passwordChange(t){this.password=t.target.value}emailChange(t){this.email=t.target.value}phoneChange(t){this.phone=t.target.value}tfaCodeChange(t){this.tfaCode=t.target.value}validateEmail(t){return!!S.test(t)||(this.processError({detailedMessage:"Email address is not valid",name:"Validation error",message:"Email address is not valid"}),!1)}renderRoute(t){var s,r,o,n,l,h;switch(t){case"login":return i("div",{class:"login-form"},!this.registrationForbidden&&i("p",{class:"login-form__register-text"},"Don't have an account?",i("a",{onClick:()=>this.openRoute("register"),class:"login-form__register-link"}," ","Sign Up")),i("input",{type:"text",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingInput",value:this.email,placeholder:"Email",onInput:t=>this.emailChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.email||!this.password)&&this.signIn()}),i("input",{type:"password",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingPassword",value:this.password,placeholder:"Password",onInput:t=>this.passwordChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.email||!this.password)&&this.signIn()}),!!this.lastError&&i("div",{class:"error",role:"alert"},null===(s=this.lastError)||void 0===s?void 0:s.detailedMessage),i("div",{class:"login-form__buttons "+(this.lastError?"login-form__buttons_mt-32":"")},i("button",{onClick:()=>this.signIn(),class:"primary-button",disabled:!this.email||!this.password},"Login"),i("a",{onClick:()=>this.openRoute("password/forgot"),class:"login-form__forgot-pass"},"Forgot password")),this.federatedProviders.length>0&&i("div",{class:"social-buttons"},i("p",{class:"social-buttons__text"},"or continue with"),i("div",{class:"social-buttons__social-medias"},this.federatedProviders.indexOf("apple")>-1&&i("div",{class:"social-buttons__media",onClick:()=>this.loginWith("apple")},i("img",{src:e("assets/images/apple.svg"),class:"social-buttons__image",alt:"login via apple"})),this.federatedProviders.indexOf("google")>-1&&i("div",{class:"social-buttons__media",onClick:()=>this.loginWith("google")},i("img",{src:e("assets/images/google.svg"),class:"social-buttons__image",alt:"login via google"})),this.federatedProviders.indexOf("facebook")>-1&&i("div",{class:"social-buttons__media",onClick:()=>this.loginWith("facebook")},i("img",{src:e("assets/images/fb.svg"),class:"social-buttons__image",alt:"login via facebook"})))));case"register":return i("div",{class:"register-form"},i("input",{type:"text",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingInput",value:this.email,placeholder:"Email",onInput:t=>this.emailChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.password||!this.email)&&this.signUp()}),i("input",{type:"password",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingPassword",value:this.password,placeholder:"Password",onInput:t=>this.passwordChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.password||!this.email)&&this.signUp()}),!!this.lastError&&i("div",{class:"error",role:"alert"},null===(r=this.lastError)||void 0===r?void 0:r.detailedMessage),i("div",{class:"register-form__buttons "+(this.lastError?"register-form__buttons_mt-32":"")},i("button",{onClick:()=>this.signUp(),class:"primary-button",disabled:!this.email||!this.password},"Continue"),i("a",{onClick:()=>this.openRoute("login"),class:"register-form__log-in"},"Go back to login")));case"otp/login":return i("div",{class:"otp-login"},!this.registrationForbidden&&i("p",{class:"otp-login__register-text"},"Don't have an account?",i("a",{onClick:()=>this.openRoute("register"),class:"login-form__register-link"}," ","Sign Up")),i("input",{type:"phone",class:"form-control",id:"floatingInput",value:this.phone,placeholder:"Phone number",onInput:t=>this.phoneChange(t)}),i("button",{onClick:()=>this.openRoute("tfa/verify"),class:"primary-button",disabled:!this.phone},"Continue"),this.federatedProviders.length>0&&i("div",{class:"social-buttons"},i("p",{class:"social-buttons__text"},"or continue with"),i("div",{class:"social-buttons__social-medias"},this.federatedProviders.indexOf("apple")>-1&&i("div",{class:"social-buttons__media",onClick:()=>this.loginWith("apple")},i("img",{src:e("assets/images/apple.svg"),class:"social-buttons__image",alt:"login via apple"})),this.federatedProviders.indexOf("google")>-1&&i("div",{class:"social-buttons__media",onClick:()=>this.loginWith("google")},i("img",{src:e("assets/images/google.svg"),class:"social-buttons__image",alt:"login via google"})),this.federatedProviders.indexOf("facebook")>-1&&i("div",{class:"social-buttons__media",onClick:()=>this.loginWith("facebook")},i("img",{src:e("assets/images/fb.svg"),class:"social-buttons__image",alt:"login via facebook"})))));case"tfa/setup":return i("div",{class:"tfa-setup"},i("p",{class:"tfa-setup__text"},"Protect your account with 2-step verification"),this.tfaType===a.TFATypeApp&&i("div",{class:"info-card"},i("div",{class:"info-card__controls"},i("p",{class:"info-card__title"},"Authenticator app"),i("button",{type:"button",class:"info-card__button",onClick:()=>this.setupTFA()},"Setup")),i("p",{class:"info-card__text"},"Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone.")),this.tfaType===a.TFATypeEmail&&i("div",{class:"info-card"},i("div",{class:"info-card__controls"},i("p",{class:"info-card__title"},"Email"),i("button",{type:"button",class:"info-card__button",onClick:()=>this.setupTFA()},"Setup")),i("p",{class:"info-card__subtitle"},this.email),i("p",{class:"info-card__text"}," Use email as 2fa, please check your email, we will send confirmation code to this email.")),this.tfaType===a.TFATypeSMS&&i("div",{class:"tfa-setup__form"},i("p",{class:"tfa-setup__subtitle"}," Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone"),i("input",{type:"phone",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingInput",value:this.phone,placeholder:"Phone",onInput:t=>this.phoneChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.phone)&&this.setupTFA()}),!!this.lastError&&i("div",{class:"error",role:"alert"},null===(o=this.lastError)||void 0===o?void 0:o.detailedMessage),i("button",{onClick:()=>this.setupTFA(),class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.phone},"Setup phone")));case"tfa/verify":return i("div",{class:"tfa-verify"},!(this.tfaType!==a.TFATypeApp)&&i("div",{class:"tfa-verify__title-wrapper"},i("h2",{class:this.provisioningURI?"tfa-verify__title":"tfa-verify__title_mb-40"},this.provisioningURI?"Please scan QR-code with the app":"Use GoogleAuth as 2fa"),!!this.provisioningURI&&i("img",{src:`data:image/png;base64, ${this.provisioningQR}`,alt:this.provisioningURI,class:"tfa-verify__qr-code"})),!(this.tfaType!==a.TFATypeSMS)&&i("div",{class:"tfa-verify__title-wrapper"},i("h2",{class:"tfa-verify__title"},"Enter the code sent to your phone number"),i("p",{class:"tfa-verify__subtitle"},"The code has been sent to ",this.phone)),!(this.tfaType!==a.TFATypeEmail)&&i("div",{class:"tfa-verify__title-wrapper"},i("h2",{class:"tfa-verify__title"},"Enter the code sent to your email address"),i("p",{class:"tfa-verify__subtitle"},"The email has been sent to ",this.email)),i("input",{type:"text",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingCode",value:this.tfaCode,placeholder:"Verify code",onInput:t=>this.tfaCodeChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.tfaCode)&&this.verifyTFA()}),!!this.lastError&&i("div",{class:"error",role:"alert"},null===(n=this.lastError)||void 0===n?void 0:n.detailedMessage),i("button",{type:"button",class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.tfaCode,onClick:()=>this.verifyTFA()},"Confirm"));case"password/forgot":return i("div",{class:"forgot-password"},i("h2",{class:"forgot-password__title"},"Enter the email you gave when you registered"),i("p",{class:"forgot-password__subtitle"},"We will send you a link to create a new password on email"),i("input",{type:"email",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingEmail",value:this.email,placeholder:"Email",onInput:t=>this.emailChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.email)&&this.restorePassword()}),!!this.lastError&&i("div",{class:"error",role:"alert"},null===(l=this.lastError)||void 0===l?void 0:l.detailedMessage),i("button",{type:"button",class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.email,onClick:()=>this.restorePassword()},"Send the link"));case"password/forgot/success":return i("div",{class:"forgot-password-success"},"dark"===this.theme&&i("img",{src:e("./assets/images/email-dark.svg"),alt:"email",class:"forgot-password-success__image"}),"light"===this.theme&&i("img",{src:e("./assets/images/email.svg"),alt:"email",class:"forgot-password-success__image"}),i("p",{class:"forgot-password-success__text"},"We sent you an email with a link to create a new password"));case"password/reset":return i("div",{class:"reset-password"},i("h2",{class:"reset-password__title"},"Set up a new password to log in to the website"),i("p",{class:"reset-password__subtitle"},"Memorize your password and do not give it to anyone."),i("input",{type:"password",class:`form-control ${this.lastError&&"form-control-danger"}`,id:"floatingPassword",value:this.password,placeholder:"Password",onInput:t=>this.passwordChange(t),onKeyPress:t=>!("Enter"!==t.key||!this.password)&&this.setNewPassword()}),!!this.lastError&&i("div",{class:"error",role:"alert"},null===(h=this.lastError)||void 0===h?void 0:h.detailedMessage),i("button",{type:"button",class:`primary-button ${this.lastError&&"primary-button-mt-32"}`,disabled:!this.password,onClick:()=>this.setNewPassword()},"Save password"));case"error":return i("div",{class:"error-view"},i("div",{class:"error-view__message"},this.lastError.message),i("div",{class:"error-view__details"},this.lastError.detailedMessage));case"callback":return i("div",{class:"error-view"},i("div",null,"Success"),this.debug&&i("div",null,i("div",null,"Access token: ",this.lastResponse.access_token),i("div",null,"Refresh token: ",this.lastResponse.refresh_token),i("div",null,"User: ",JSON.stringify(this.lastResponse.user))));case"loading":return i("div",{class:"error-view"},i("div",null,"Loading ..."))}}async componentWillLoad(){const t=this.postLogoutRedirectUri||window.location.origin+window.location.pathname;this.auth=new R({appId:this.appId,url:this.url,postLogoutRedirectUri:t});try{const t=await this.auth.api.getAppSettings();this.registrationForbidden=t.registrationForbidden,this.tfaType=t.tfaType,this.federatedProviders=t.federatedProviders}catch(t){this.route="error",this.lastError=t}const s=new URL(window.location.href);if(s.searchParams.get("provider")&&s.searchParams.get("state")){const t=new URL(window.location.href),i=new URLSearchParams,e=s.searchParams.get("appId");i.set("appId",e),window.history.replaceState({},document.title,`${t.pathname}?${i.toString()}`),this.route="loading",this.auth.api.federatedLoginComplete(t.searchParams).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((t=>this.openRoute(t))).catch((t=>this.processError(t)))}}componentWillRender(){if("callback"===this.route){const t=new URL(window.location.href);t.searchParams.set("callbackUrl",this.lastResponse.callbackUrl),window.history.replaceState({},document.title,`${t.pathname}?${t.searchParams.toString()}`),this.complete.emit(this.lastResponse)}"logout"===this.route&&this.complete.emit()}render(){return i(r,null,i("div",{class:{wrapper:"light"===this.theme,"wrapper-dark":"dark"===this.theme}},this.renderRoute(this.route)),i("div",{class:"error-view"},this.debug&&i("div",null,i("br",null),this.appId)))}static get assetsDirs(){return["assets"]}};export{x as identifo_form}
var __extends=this&&this.__extends||function(){var e=function(t,r){e=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(e,t){e.__proto__=t}||function(e,t){for(var r in t)if(Object.prototype.hasOwnProperty.call(t,r))e[r]=t[r]};return e(t,r)};return function(t,r){if(typeof r!=="function"&&r!==null)throw new TypeError("Class extends value "+String(r)+" is not a constructor or null");e(t,r);function n(){this.constructor=t}t.prototype=r===null?Object.create(r):(n.prototype=r.prototype,new n)}}();var __awaiter=this&&this.__awaiter||function(e,t,r,n){function o(e){return e instanceof r?e:new r((function(t){t(e)}))}return new(r||(r=Promise))((function(r,i){function s(e){try{c(n.next(e))}catch(e){i(e)}}function a(e){try{c(n["throw"](e))}catch(e){i(e)}}function c(e){e.done?r(e.value):o(e.value).then(s,a)}c((n=n.apply(e,t||[])).next())}))};var __generator=this&&this.__generator||function(e,t){var r={label:0,sent:function(){if(i[0]&1)throw i[1];return i[1]},trys:[],ops:[]},n,o,i,s;return s={next:a(0),throw:a(1),return:a(2)},typeof Symbol==="function"&&(s[Symbol.iterator]=function(){return this}),s;function a(e){return function(t){return c([e,t])}}function c(s){if(n)throw new TypeError("Generator is already executing.");while(r)try{if(n=1,o&&(i=s[0]&2?o["return"]:s[0]?o["throw"]||((i=o["return"])&&i.call(o),0):o.next)&&!(i=i.call(o,s[1])).done)return i;if(o=0,i)s=[s[0]&2,i.value];switch(s[0]){case 0:case 1:i=s;break;case 4:r.label++;return{value:s[1],done:false};case 5:r.label++;o=s[1];s=[0];continue;case 7:s=r.ops.pop();r.trys.pop();continue;default:if(!(i=r.trys,i=i.length>0&&i[i.length-1])&&(s[0]===6||s[0]===2)){r=0;continue}if(s[0]===3&&(!i||s[1]>i[0]&&s[1]<i[3])){r.label=s[1];break}if(s[0]===6&&r.label<i[1]){r.label=i[1];i=s;break}if(i&&r.label<i[2]){r.label=i[2];r.ops.push(s);break}if(i[2])r.ops.pop();r.trys.pop();continue}s=t.call(e,r)}catch(e){s=[6,e];o=0}finally{n=i=0}if(s[0]&5)throw s[1];return{value:s[0]?s[1]:void 0,done:true}}};System.register(["./index-7003c359.system.js"],(function(e){"use strict";var t,r,n,o,i;return{setters:[function(e){t=e.r;r=e.c;n=e.h;o=e.g;i=e.H}],execute:function(){var s;(function(e){e["PleaseEnableTFA"]="error.api.request.2fa.please_enable";e["NetworkError"]="error.network"})(s||(s={}));var a;(function(e){e["TFATypeApp"]="app";e["TFATypeSMS"]="sms";e["TFATypeEmail"]="email"})(a||(a={}));var c=function(e){__extends(t,e);function t(t){var r=e.call(this,(t==null?void 0:t.message)||"Unknown API error")||this;r.detailedMessage=t==null?void 0:t.detailed_message;r.id=t==null?void 0:t.id;r.status=t==null?void 0:t.status;return r}return t}(Error);var u=Object.defineProperty;var l=Object.getOwnPropertySymbols;var h=Object.prototype.hasOwnProperty;var p=Object.prototype.propertyIsEnumerable;var f=function(e,t,r){return t in e?u(e,t,{enumerable:true,configurable:true,writable:true,value:r}):e[t]=r};var d=function(e,t){for(var r in t||(t={}))if(h.call(t,r))f(e,r,t[r]);if(l)for(var n=0,o=l(t);n<o.length;n++){var r=o[n];if(p.call(t,r))f(e,r,t[r])}return e};var g=function(e,t,r){return new Promise((function(n,o){var i=function(e){try{a(r.next(e))}catch(e){o(e)}};var s=function(e){try{a(r.throw(e))}catch(e){o(e)}};var a=function(e){return e.done?n(e.value):Promise.resolve(e.value).then(i,s)};a((r=r.apply(e,t)).next())}))};var v="X-Identifo-Clientid";var w="Authorization";var m=function(){function e(e,t){var r;var n=this;this.config=e;this.tokenService=t;this.defaultHeaders=(r={},r[v]="",r.Accept="application/json",r["Content-Type"]="application/json",r);this.catchNetworkErrorHandler=function(e){if(e.message==="Network Error"||e.message==="Failed to fetch"||e.message==="Preflight response is not successful"||e.message.indexOf("is not allowed by Access-Control-Allow-Origin")>-1){console.error(e.message);throw new c({id:s.NetworkError,status:0,message:"Configuration error",detailed_message:'Please check Identifo URL and add "'+window.location.protocol+"//"+window.location.host+'" to "REDIRECT URLS" in Identifo app settings.'})}throw e};this.checkStatusCodeAndGetJSON=function(e){return g(n,null,(function(){var t;return __generator(this,(function(r){switch(r.label){case 0:if(!!e.ok)return[3,2];return[4,e.json()];case 1:t=r.sent();throw new c(t==null?void 0:t.error);case 2:return[2,e.json()]}}))}))};this.baseUrl=e.url.replace(/\/$/,"");this.defaultHeaders[v]=e.appId;this.appId=e.appId}e.prototype.get=function(e,t){return this.send(e,d({method:"GET"},t))};e.prototype.put=function(e,t,r){return this.send(e,d({method:"PUT",body:JSON.stringify(t)},r))};e.prototype.post=function(e,t,r){return this.send(e,d({method:"POST",body:JSON.stringify(t)},r))};e.prototype.send=function(e,t){var r=d({},t);r.credentials="include";r.headers=d(d({},r.headers),this.defaultHeaders);return fetch(""+this.baseUrl+e,r).catch(this.catchNetworkErrorHandler).then(this.checkStatusCodeAndGetJSON).then((function(e){return e}))};e.prototype.getUser=function(){return g(this,null,(function(){var e,t;var r;return __generator(this,(function(n){if(!((e=this.tokenService.getToken())==null?void 0:e.token)){throw new Error("No token in token service.")}return[2,this.get("/me",{headers:(r={},r[w]="Bearer "+((t=this.tokenService.getToken())==null?void 0:t.token),r)})]}))}))};e.prototype.renewToken=function(){return g(this,null,(function(){var e,t;var r;var n=this;return __generator(this,(function(o){if(!((e=this.tokenService.getToken("refresh"))==null?void 0:e.token)){throw new Error("No token in token service.")}return[2,this.post("/auth/token",{scopes:this.config.scopes},{headers:(r={},r[w]="Bearer "+((t=this.tokenService.getToken("refresh"))==null?void 0:t.token),r)}).then((function(e){return n.storeToken(e)}))]}))}))};e.prototype.updateUser=function(e){return g(this,null,(function(){var t,r;var n;return __generator(this,(function(o){if(!((t=this.tokenService.getToken())==null?void 0:t.token)){throw new Error("No token in token service.")}return[2,this.put("/me",e,{headers:(n={},n[w]="Bearer "+((r=this.tokenService.getToken("access"))==null?void 0:r.token),n)})]}))}))};e.prototype.login=function(e,t,r,n){return g(this,null,(function(){var o;var i=this;return __generator(this,(function(s){o={email:e,password:t,device_token:r,scopes:n};return[2,this.post("/auth/login",o).then((function(e){return i.storeToken(e)}))]}))}))};e.prototype.federatedLogin=function(e,t,r,n){return g(this,arguments,(function(e,t,r,n,o){var i,s,a,c,u;if(o===void 0){o={width:600,height:800,popUp:false}}return __generator(this,(function(l){i=document.createElement("form");i.style.display="none";if(o.popUp){i.target="TargetWindow"}i.method="POST";s=new URLSearchParams;s.set("appId",this.config.appId);s.set("provider",e);s.set("scopes",t.join(","));s.set("redirectUrl",r);if(n){s.set("callbackUrl",n)}i.action=this.baseUrl+"/auth/federated?"+s.toString();document.body.appendChild(i);if(o.popUp){a=window.screenX+window.outerWidth/2-(o.width||600)/2;c=window.screenY+window.outerHeight/2-(o.height||800)/2;u=window.open("","TargetWindow","status=0,title=0,height="+o.height+",width="+o.width+",top="+c+",left="+a+",scrollbars=1");if(u){i.submit()}}else{window.location.assign(this.baseUrl+"/auth/federated?"+s.toString())}return[2]}))}))};e.prototype.federatedLoginComplete=function(e){return g(this,null,(function(){var t=this;return __generator(this,(function(r){return[2,this.get("/auth/federated/complete?"+e.toString()).then((function(e){return t.storeToken(e)}))]}))}))};e.prototype.register=function(e,t,r){return g(this,null,(function(){var n;var o=this;return __generator(this,(function(i){n={email:e,password:t,scopes:r};return[2,this.post("/auth/register",n).then((function(e){return o.storeToken(e)}))]}))}))};e.prototype.requestResetPassword=function(e){return g(this,null,(function(){var t;return __generator(this,(function(r){t={email:e};return[2,this.post("/auth/request_reset_password",t)]}))}))};e.prototype.resetPassword=function(e){return g(this,null,(function(){var t,r,n;var o;return __generator(this,(function(i){if(!((t=this.tokenService.getToken())==null?void 0:t.token)){throw new Error("No token in token service.")}n={password:e};return[2,this.post("/auth/reset_password",n,{headers:(o={},o[w]="Bearer "+((r=this.tokenService.getToken())==null?void 0:r.token),o)})]}))}))};e.prototype.getAppSettings=function(){return g(this,null,(function(){return __generator(this,(function(e){return[2,this.get("/auth/app_settings")]}))}))};e.prototype.enableTFA=function(){return g(this,null,(function(){var e,t;var r;return __generator(this,(function(n){if(!((e=this.tokenService.getToken())==null?void 0:e.token)){throw new Error("No token in token service.")}return[2,this.put("/auth/tfa/enable",{},{headers:(r={},r[w]="BEARER "+((t=this.tokenService.getToken())==null?void 0:t.token),r)})]}))}))};e.prototype.verifyTFA=function(e,t){return g(this,null,(function(){var r,n;var o;var i=this;return __generator(this,(function(s){if(!((r=this.tokenService.getToken())==null?void 0:r.token)){throw new Error("No token in token service.")}return[2,this.post("/auth/tfa/login",{tfa_code:e,scopes:t},{headers:(o={},o[w]="BEARER "+((n=this.tokenService.getToken())==null?void 0:n.token),o)}).then((function(e){return i.storeToken(e)}))]}))}))};e.prototype.logout=function(){return g(this,null,(function(){var e,t,r;var n;return __generator(this,(function(o){if(!((e=this.tokenService.getToken())==null?void 0:e.token)){throw new Error("No token in token service.")}return[2,this.post("/me/logout",{refresh_token:(t=this.tokenService.getToken("refresh"))==null?void 0:t.token},{headers:(n={},n[w]="Bearer "+((r=this.tokenService.getToken())==null?void 0:r.token),n)})]}))}))};e.prototype.storeToken=function(e){if(e.access_token){this.tokenService.saveToken(e.access_token,"access")}if(e.refresh_token){this.tokenService.saveToken(e.refresh_token,"refresh")}return e};return e}();var _=/^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-=]*$)/;var y="Empty or invalid token";var k="token";var b="refresh_token";var T=function(){function e(e,t,r){this.preffix="identifo_";this.storageType="localStorage";this.access=this.preffix+"access_token";this.refresh=this.preffix+"refresh_token";this.isAccessible=true;this.access=t?this.preffix+t:this.access;this.refresh=r?this.preffix+r:this.refresh;this.storageType=e}e.prototype.saveToken=function(e,t){if(e){window[this.storageType].setItem(this[t],e);return true}return false};e.prototype.getToken=function(e){var t;return(t=window[this.storageType].getItem(this[e]))!=null?t:""};e.prototype.deleteToken=function(e){window[this.storageType].removeItem(this[e])};return e}();var E=function(e){__extends(t,e);function t(t,r){return e.call(this,"localStorage",t,r)||this}return t}(T);var S=function(e,t,r){return new Promise((function(n,o){var i=function(e){try{a(r.next(e))}catch(e){o(e)}};var s=function(e){try{a(r.throw(e))}catch(e){o(e)}};var a=function(e){return e.done?n(e.value):Promise.resolve(e.value).then(i,s)};a((r=r.apply(e,t)).next())}))};var P=function(){function e(e){this.tokenManager=e||new E}e.prototype.handleVerification=function(e,t,r){return S(this,null,(function(){var n;return __generator(this,(function(o){switch(o.label){case 0:if(!this.tokenManager.isAccessible)return[2,true];o.label=1;case 1:o.trys.push([1,3,,4]);return[4,this.validateToken(e,t,r)];case 2:o.sent();this.saveToken(e);return[2,true];case 3:n=o.sent();this.removeToken();return[2,Promise.reject(n)];case 4:return[2]}}))}))};e.prototype.validateToken=function(e,t,r){return S(this,null,(function(){var n,o,i;return __generator(this,(function(s){if(!e)throw new Error(y);o=this.parseJWT(e);i=this.isJWTExpired(o);if(((n=o.aud)==null?void 0:n.includes(t))&&(!r||o.iss===r)&&!i){return[2,Promise.resolve(true)]}throw new Error(y)}))}))};e.prototype.parseJWT=function(e){var t=e.split(".")[1];if(!t)return{aud:[],iss:"",exp:10};var r=t.replace(/-/g,"+").replace(/_/g,"/");var n=decodeURIComponent(atob(r).split("").map((function(e){return"%"+("00"+e.charCodeAt(0).toString(16)).slice(-2)})).join(""));return JSON.parse(n)};e.prototype.isJWTExpired=function(e){var t=(new Date).getTime()/1e3;if(e.exp&&t>e.exp){return true}return false};e.prototype.isAuthenticated=function(e,t){if(!this.tokenManager.isAccessible)return Promise.resolve(true);var r=this.tokenManager.getToken("access");return this.validateToken(r,e,t)};e.prototype.saveToken=function(e,t){if(t===void 0){t="access"}return this.tokenManager.saveToken(e,t)};e.prototype.removeToken=function(e){if(e===void 0){e="access"}this.tokenManager.deleteToken(e)};e.prototype.getToken=function(e){if(e===void 0){e="access"}var t=this.tokenManager.getToken(e);if(!t)return null;var r=this.parseJWT(t);return{token:t,payload:r}};return e}();var U=function(){function e(e){this.config=e}e.prototype.getUrl=function(e){var t,r;var n=((t=this.config.scopes)==null?void 0:t.join())||"";var o=encodeURIComponent((r=this.config.redirectUri)!=null?r:window.location.href);var i="appId="+this.config.appId+"&scopes="+n;var s=i+"&callbackUrl="+o;var a=this.config.postLogoutRedirectUri?"&callbackUrl="+encodeURIComponent(this.config.postLogoutRedirectUri):"&callbackUrl="+o+"&redirectUri="+this.config.url+"/web/login?"+encodeURIComponent(i);var c={signup:this.config.url+"/web/register?"+s,signin:this.config.url+"/web/login?"+s,logout:this.config.url+"/web/logout?"+i+a,renew:this.config.url+"/web/token/renew?"+i+"&redirectUri="+o,default:"default"};return c[e]||c.default};e.prototype.createSignupUrl=function(){return this.getUrl("signup")};e.prototype.createSigninUrl=function(){return this.getUrl("signin")};e.prototype.createLogoutUrl=function(){return this.getUrl("logout")};e.prototype.createRenewSessionUrl=function(){return this.getUrl("renew")};return e}();var R=Object.defineProperty;var A=Object.defineProperties;var C=Object.getOwnPropertyDescriptors;var I=Object.getOwnPropertySymbols;var x=Object.prototype.hasOwnProperty;var O=Object.prototype.propertyIsEnumerable;var F=function(e,t,r){return t in e?R(e,t,{enumerable:true,configurable:true,writable:true,value:r}):e[t]=r};var j=function(e,t){for(var r in t||(t={}))if(x.call(t,r))F(e,r,t[r]);if(I)for(var n=0,o=I(t);n<o.length;n++){var r=o[n];if(O.call(t,r))F(e,r,t[r])}return e};var L=function(e,t){return A(e,C(t))};var W=function(e,t,r){return new Promise((function(n,o){var i=function(e){try{a(r.next(e))}catch(e){o(e)}};var s=function(e){try{a(r.throw(e))}catch(e){o(e)}};var a=function(e){return e.done?n(e.value):Promise.resolve(e.value).then(i,s)};a((r=r.apply(e,t)).next())}))};var M=function(){function e(e){this.token=null;this.isAuth=false;var t,r;this.config=L(j({},e),{autoRenew:(t=e.autoRenew)!=null?t:true});this.tokenService=new P(e.tokenManager);this.urlBuilder=new U(this.config);this.api=new m(e,this.tokenService);this.handleToken(((r=this.tokenService.getToken())==null?void 0:r.token)||"","access")}e.prototype.handleToken=function(e,t){if(e){if(t==="access"){var r=this.tokenService.parseJWT(e);this.token={token:e,payload:r};this.isAuth=true;this.tokenService.saveToken(e)}else{this.tokenService.saveToken(e,"refresh")}}};e.prototype.resetAuthValues=function(){this.token=null;this.isAuth=false;this.tokenService.removeToken();this.tokenService.removeToken("refresh")};e.prototype.signup=function(){window.location.href=this.urlBuilder.createSignupUrl()};e.prototype.signin=function(){window.location.href=this.urlBuilder.createSigninUrl()};e.prototype.logout=function(){this.resetAuthValues();window.location.href=this.urlBuilder.createLogoutUrl()};e.prototype.handleAuthentication=function(){return W(this,null,(function(){var e,t,r,n;return __generator(this,(function(o){switch(o.label){case 0:e=this.getTokenFromUrl(),t=e.access,r=e.refresh;if(!t){this.resetAuthValues();return[2,Promise.reject()]}o.label=1;case 1:o.trys.push([1,4,6,7]);return[4,this.tokenService.handleVerification(t,this.config.appId,this.config.issuer)];case 2:o.sent();this.handleToken(t,"access");if(r){this.handleToken(r,"refresh")}return[4,Promise.resolve(true)];case 3:return[2,o.sent()];case 4:n=o.sent();this.resetAuthValues();return[4,Promise.reject()];case 5:return[2,o.sent()];case 6:window.history.pushState({},document.title,window.location.pathname);return[7];case 7:return[2]}}))}))};e.prototype.getTokenFromUrl=function(){var e=new URLSearchParams(window.location.search);var t={access:"",refresh:""};var r=e.get(k);var n=e.get(b);if(n&&_.test(n)){t.refresh=n}if(r&&_.test(r)){t.access=r}return t};e.prototype.getToken=function(){return W(this,null,(function(){var e,t,r,n;return __generator(this,(function(o){switch(o.label){case 0:e=this.tokenService.getToken();t=this.tokenService.getToken("refresh");if(!e)return[3,6];r=this.tokenService.isJWTExpired(e.payload);if(!(r&&t))return[3,5];o.label=1;case 1:o.trys.push([1,4,,5]);return[4,this.renewSession()];case 2:o.sent();return[4,Promise.resolve(this.token)];case 3:return[2,o.sent()];case 4:n=o.sent();this.resetAuthValues();throw new Error("No token");case 5:return[2,Promise.resolve(e)];case 6:return[2,Promise.resolve(null)]}}))}))};e.prototype.renewSession=function(){return W(this,null,(function(){var e,t,r,n;return __generator(this,(function(o){switch(o.label){case 0:o.trys.push([0,3,,4]);return[4,this.renewSessionWithToken()];case 1:e=o.sent(),t=e.access,r=e.refresh;this.handleToken(t,"access");this.handleToken(r,"refresh");return[4,Promise.resolve(t)];case 2:return[2,o.sent()];case 3:n=o.sent();return[2,Promise.reject()];case 4:return[2]}}))}))};e.prototype.renewSessionWithToken=function(){return W(this,null,(function(){var e,t;return __generator(this,(function(r){switch(r.label){case 0:r.trys.push([0,2,,3]);return[4,this.api.renewToken().then((function(e){return{access:e.access_token||"",refresh:e.refresh_token||""}}))];case 1:e=r.sent();return[2,e];case 2:t=r.sent();return[2,Promise.reject(t)];case 3:return[2]}}))}))};return e}();var N=/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;var J=e("identifo_form",function(){function e(e){var n=this;t(this,e);this.complete=r(this,"complete",7);this.error=r(this,"error",7);this.route="login";this.theme="light";this.scopes="";this.afterLoginRedirect=function(e){n.phone=e.user.phone||"";n.email=e.user.email||"";n.lastResponse=e;if(e.require_2fa){if(!e.enabled_2fa){return"tfa/setup"}if(e.enabled_2fa){return"tfa/verify"}}if(e.access_token&&e.refresh_token){return"callback"}if(e.access_token&&!e.refresh_token){return"callback"}};this.loginCatchRedirect=function(e){if(e.id===s.PleaseEnableTFA){return"tfa/setup"}throw e}}e.prototype.processError=function(e){this.lastError=e;this.error.emit(e)};e.prototype.signIn=function(){return __awaiter(this,void 0,void 0,(function(){var e=this;return __generator(this,(function(t){switch(t.label){case 0:return[4,this.auth.api.login(this.email,this.password,"",this.scopes.split(",")).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((function(t){return e.openRoute(t)})).catch((function(t){return e.processError(t)}))];case 1:t.sent();return[2]}}))}))};e.prototype.loginWith=function(e){return __awaiter(this,void 0,void 0,(function(){var t;return __generator(this,(function(r){this.route="loading";t=this.federatedRedirectUrl||window.location.origin+window.location.pathname;this.auth.api.federatedLogin(e,this.scopes.split(","),t,this.callbackUrl);return[2]}))}))};e.prototype.signUp=function(){return __awaiter(this,void 0,void 0,(function(){var e=this;return __generator(this,(function(t){switch(t.label){case 0:if(!this.validateEmail(this.email)){return[2]}return[4,this.auth.api.register(this.email,this.password,this.scopes.split(",")).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((function(t){return e.openRoute(t)})).catch((function(t){return e.processError(t)}))];case 1:t.sent();return[2]}}))}))};e.prototype.verifyTFA=function(){return __awaiter(this,void 0,void 0,(function(){var e=this;return __generator(this,(function(t){this.auth.api.verifyTFA(this.tfaCode,[]).then((function(){return e.openRoute("callback")})).catch((function(t){return e.processError(t)}));return[2]}))}))};e.prototype.setupTFA=function(){return __awaiter(this,void 0,void 0,(function(){var e;var t=this;return __generator(this,(function(r){switch(r.label){case 0:if(!(this.tfaType==a.TFATypeSMS))return[3,4];r.label=1;case 1:r.trys.push([1,3,,4]);return[4,this.auth.api.updateUser({new_phone:this.phone})];case 2:r.sent();return[3,4];case 3:e=r.sent();this.processError(e);return[2];case 4:return[4,this.auth.api.enableTFA().then((function(e){if(!e.provisioning_uri){t.openRoute("tfa/verify")}if(e.provisioning_uri){t.provisioningURI=e.provisioning_uri;t.provisioningQR=e.provisioning_qr;t.openRoute("tfa/verify")}}))];case 5:r.sent();return[2]}}))}))};e.prototype.restorePassword=function(){var e=this;this.auth.api.requestResetPassword(this.email).then((function(){e.success=true;e.openRoute("password/forgot/success")})).catch((function(t){return e.processError(t)}))};e.prototype.setNewPassword=function(){var e=this;if(this.token){this.auth.tokenService.saveToken(this.token,"access")}this.auth.api.resetPassword(this.password).then((function(){e.success=true;e.openRoute("login");e.password=""})).catch((function(t){return e.processError(t)}))};e.prototype.openRoute=function(e){this.lastError=undefined;this.route=e};e.prototype.usernameChange=function(e){this.username=e.target.value};e.prototype.passwordChange=function(e){this.password=e.target.value};e.prototype.emailChange=function(e){this.email=e.target.value};e.prototype.phoneChange=function(e){this.phone=e.target.value};e.prototype.tfaCodeChange=function(e){this.tfaCode=e.target.value};e.prototype.validateEmail=function(e){if(!N.test(e)){this.processError({detailedMessage:"Email address is not valid",name:"Validation error",message:"Email address is not valid"});return false}return true};e.prototype.renderRoute=function(e){var t=this;var r,i,s,c,u,l;switch(e){case"login":return n("div",{class:"login-form"},!this.registrationForbidden&&n("p",{class:"login-form__register-text"},"Don't have an account?",n("a",{onClick:function(){return t.openRoute("register")},class:"login-form__register-link"}," ","Sign Up")),n("input",{type:"text",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingInput",value:this.email,placeholder:"Email",onInput:function(e){return t.emailChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.email&&t.password)&&t.signIn()}}),n("input",{type:"password",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingPassword",value:this.password,placeholder:"Password",onInput:function(e){return t.passwordChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.email&&t.password)&&t.signIn()}}),!!this.lastError&&n("div",{class:"error",role:"alert"},(r=this.lastError)===null||r===void 0?void 0:r.detailedMessage),n("div",{class:"login-form__buttons "+(!!this.lastError?"login-form__buttons_mt-32":"")},n("button",{onClick:function(){return t.signIn()},class:"primary-button",disabled:!this.email||!this.password},"Login"),n("a",{onClick:function(){return t.openRoute("password/forgot")},class:"login-form__forgot-pass"},"Forgot password")),this.federatedProviders.length>0&&n("div",{class:"social-buttons"},n("p",{class:"social-buttons__text"},"or continue with"),n("div",{class:"social-buttons__social-medias"},this.federatedProviders.indexOf("apple")>-1&&n("div",{class:"social-buttons__media",onClick:function(){return t.loginWith("apple")}},n("img",{src:o("assets/images/"+"apple.svg"),class:"social-buttons__image",alt:"login via apple"})),this.federatedProviders.indexOf("google")>-1&&n("div",{class:"social-buttons__media",onClick:function(){return t.loginWith("google")}},n("img",{src:o("assets/images/"+"google.svg"),class:"social-buttons__image",alt:"login via google"})),this.federatedProviders.indexOf("facebook")>-1&&n("div",{class:"social-buttons__media",onClick:function(){return t.loginWith("facebook")}},n("img",{src:o("assets/images/"+"fb.svg"),class:"social-buttons__image",alt:"login via facebook"})))));case"register":return n("div",{class:"register-form"},n("input",{type:"text",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingInput",value:this.email,placeholder:"Email",onInput:function(e){return t.emailChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.password&&t.email)&&t.signUp()}}),n("input",{type:"password",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingPassword",value:this.password,placeholder:"Password",onInput:function(e){return t.passwordChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.password&&t.email)&&t.signUp()}}),!!this.lastError&&n("div",{class:"error",role:"alert"},(i=this.lastError)===null||i===void 0?void 0:i.detailedMessage),n("div",{class:"register-form__buttons "+(!!this.lastError?"register-form__buttons_mt-32":"")},n("button",{onClick:function(){return t.signUp()},class:"primary-button",disabled:!this.email||!this.password},"Continue"),n("a",{onClick:function(){return t.openRoute("login")},class:"register-form__log-in"},"Go back to login")));case"otp/login":return n("div",{class:"otp-login"},!this.registrationForbidden&&n("p",{class:"otp-login__register-text"},"Don't have an account?",n("a",{onClick:function(){return t.openRoute("register")},class:"login-form__register-link"}," ","Sign Up")),n("input",{type:"phone",class:"form-control",id:"floatingInput",value:this.phone,placeholder:"Phone number",onInput:function(e){return t.phoneChange(e)}}),n("button",{onClick:function(){return t.openRoute("tfa/verify")},class:"primary-button",disabled:!this.phone},"Continue"),this.federatedProviders.length>0&&n("div",{class:"social-buttons"},n("p",{class:"social-buttons__text"},"or continue with"),n("div",{class:"social-buttons__social-medias"},this.federatedProviders.indexOf("apple")>-1&&n("div",{class:"social-buttons__media",onClick:function(){return t.loginWith("apple")}},n("img",{src:o("assets/images/"+"apple.svg"),class:"social-buttons__image",alt:"login via apple"})),this.federatedProviders.indexOf("google")>-1&&n("div",{class:"social-buttons__media",onClick:function(){return t.loginWith("google")}},n("img",{src:o("assets/images/"+"google.svg"),class:"social-buttons__image",alt:"login via google"})),this.federatedProviders.indexOf("facebook")>-1&&n("div",{class:"social-buttons__media",onClick:function(){return t.loginWith("facebook")}},n("img",{src:o("assets/images/"+"fb.svg"),class:"social-buttons__image",alt:"login via facebook"})))));case"tfa/setup":return n("div",{class:"tfa-setup"},n("p",{class:"tfa-setup__text"},"Protect your account with 2-step verification"),this.tfaType===a.TFATypeApp&&n("div",{class:"info-card"},n("div",{class:"info-card__controls"},n("p",{class:"info-card__title"},"Authenticator app"),n("button",{type:"button",class:"info-card__button",onClick:function(){return t.setupTFA()}},"Setup")),n("p",{class:"info-card__text"},"Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone.")),this.tfaType===a.TFATypeEmail&&n("div",{class:"info-card"},n("div",{class:"info-card__controls"},n("p",{class:"info-card__title"},"Email"),n("button",{type:"button",class:"info-card__button",onClick:function(){return t.setupTFA()}},"Setup")),n("p",{class:"info-card__subtitle"},this.email),n("p",{class:"info-card__text"}," Use email as 2fa, please check your email, we will send confirmation code to this email.")),this.tfaType===a.TFATypeSMS&&n("div",{class:"tfa-setup__form"},n("p",{class:"tfa-setup__subtitle"}," Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone"),n("input",{type:"phone",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingInput",value:this.phone,placeholder:"Phone",onInput:function(e){return t.phoneChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.phone)&&t.setupTFA()}}),!!this.lastError&&n("div",{class:"error",role:"alert"},(s=this.lastError)===null||s===void 0?void 0:s.detailedMessage),n("button",{onClick:function(){return t.setupTFA()},class:"primary-button "+(this.lastError&&"primary-button-mt-32"),disabled:!this.phone},"Setup phone")));case"tfa/verify":return n("div",{class:"tfa-verify"},!!(this.tfaType===a.TFATypeApp)&&n("div",{class:"tfa-verify__title-wrapper"},n("h2",{class:this.provisioningURI?"tfa-verify__title":"tfa-verify__title_mb-40"},!!this.provisioningURI?"Please scan QR-code with the app":"Use GoogleAuth as 2fa"),!!this.provisioningURI&&n("img",{src:"data:image/png;base64, "+this.provisioningQR,alt:this.provisioningURI,class:"tfa-verify__qr-code"})),!!(this.tfaType===a.TFATypeSMS)&&n("div",{class:"tfa-verify__title-wrapper"},n("h2",{class:"tfa-verify__title"},"Enter the code sent to your phone number"),n("p",{class:"tfa-verify__subtitle"},"The code has been sent to ",this.phone)),!!(this.tfaType===a.TFATypeEmail)&&n("div",{class:"tfa-verify__title-wrapper"},n("h2",{class:"tfa-verify__title"},"Enter the code sent to your email address"),n("p",{class:"tfa-verify__subtitle"},"The email has been sent to ",this.email)),n("input",{type:"text",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingCode",value:this.tfaCode,placeholder:"Verify code",onInput:function(e){return t.tfaCodeChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.tfaCode)&&t.verifyTFA()}}),!!this.lastError&&n("div",{class:"error",role:"alert"},(c=this.lastError)===null||c===void 0?void 0:c.detailedMessage),n("button",{type:"button",class:"primary-button "+(this.lastError&&"primary-button-mt-32"),disabled:!this.tfaCode,onClick:function(){return t.verifyTFA()}},"Confirm"));case"password/forgot":return n("div",{class:"forgot-password"},n("h2",{class:"forgot-password__title"},"Enter the email you gave when you registered"),n("p",{class:"forgot-password__subtitle"},"We will send you a link to create a new password on email"),n("input",{type:"email",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingEmail",value:this.email,placeholder:"Email",onInput:function(e){return t.emailChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.email)&&t.restorePassword()}}),!!this.lastError&&n("div",{class:"error",role:"alert"},(u=this.lastError)===null||u===void 0?void 0:u.detailedMessage),n("button",{type:"button",class:"primary-button "+(this.lastError&&"primary-button-mt-32"),disabled:!this.email,onClick:function(){return t.restorePassword()}},"Send the link"));case"password/forgot/success":return n("div",{class:"forgot-password-success"},this.theme==="dark"&&n("img",{src:o("./assets/images/"+"email-dark.svg"),alt:"email",class:"forgot-password-success__image"}),this.theme==="light"&&n("img",{src:o("./assets/images/"+"email.svg"),alt:"email",class:"forgot-password-success__image"}),n("p",{class:"forgot-password-success__text"},"We sent you an email with a link to create a new password"));case"password/reset":return n("div",{class:"reset-password"},n("h2",{class:"reset-password__title"},"Set up a new password to log in to the website"),n("p",{class:"reset-password__subtitle"},"Memorize your password and do not give it to anyone."),n("input",{type:"password",class:"form-control "+(this.lastError&&"form-control-danger"),id:"floatingPassword",value:this.password,placeholder:"Password",onInput:function(e){return t.passwordChange(e)},onKeyPress:function(e){return!!(e.key==="Enter"&&t.password)&&t.setNewPassword()}}),!!this.lastError&&n("div",{class:"error",role:"alert"},(l=this.lastError)===null||l===void 0?void 0:l.detailedMessage),n("button",{type:"button",class:"primary-button "+(this.lastError&&"primary-button-mt-32"),disabled:!this.password,onClick:function(){return t.setNewPassword()}},"Save password"));case"error":return n("div",{class:"error-view"},n("div",{class:"error-view__message"},this.lastError.message),n("div",{class:"error-view__details"},this.lastError.detailedMessage));case"callback":return n("div",{class:"error-view"},n("div",null,"Success"),this.debug&&n("div",null,n("div",null,"Access token: ",this.lastResponse.access_token),n("div",null,"Refresh token: ",this.lastResponse.refresh_token),n("div",null,"User: ",JSON.stringify(this.lastResponse.user))));case"loading":return n("div",{class:"error-view"},n("div",null,"Loading ..."))}};e.prototype.componentWillLoad=function(){return __awaiter(this,void 0,void 0,(function(){var e,t,r,n,o,i,s;var a=this;return __generator(this,(function(c){switch(c.label){case 0:e=this.postLogoutRedirectUri||window.location.origin+window.location.pathname;this.auth=new M({appId:this.appId,url:this.url,postLogoutRedirectUri:e});c.label=1;case 1:c.trys.push([1,3,,4]);return[4,this.auth.api.getAppSettings()];case 2:t=c.sent();this.registrationForbidden=t.registrationForbidden;this.tfaType=t.tfaType;this.federatedProviders=t.federatedProviders;return[3,4];case 3:r=c.sent();this.route="error";this.lastError=r;return[3,4];case 4:n=new URL(window.location.href);if(!!n.searchParams.get("provider")&&!!n.searchParams.get("state")){o=new URL(window.location.href);i=new URLSearchParams;s=n.searchParams.get("appId");i.set("appId",s);window.history.replaceState({},document.title,o.pathname+"?"+i.toString());this.route="loading";this.auth.api.federatedLoginComplete(o.searchParams).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).then((function(e){return a.openRoute(e)})).catch((function(e){return a.processError(e)}))}return[2]}}))}))};e.prototype.componentWillRender=function(){if(this.route==="callback"){var e=new URL(window.location.href);e.searchParams.set("callbackUrl",this.lastResponse.callbackUrl);window.history.replaceState({},document.title,e.pathname+"?"+e.searchParams.toString());this.complete.emit(this.lastResponse)}if(this.route==="logout"){this.complete.emit()}};e.prototype.render=function(){return n(i,null,n("div",{class:{wrapper:this.theme==="light","wrapper-dark":this.theme==="dark"}},this.renderRoute(this.route)),n("div",{class:"error-view"},this.debug&&n("div",null,n("br",null),this.appId)))};Object.defineProperty(e,"assetsDirs",{get:function(){return["assets"]},enumerable:false,configurable:true});return e}())}}}));
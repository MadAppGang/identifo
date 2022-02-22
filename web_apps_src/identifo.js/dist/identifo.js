'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var rxjs = require('rxjs');

exports.APIErrorCodes = void 0;
(function(APIErrorCodes2) {
  APIErrorCodes2["PleaseEnableTFA"] = "error.api.request.2fa.please_enable";
  APIErrorCodes2["InvalidCallbackURL"] = "error.api.request.callbackurl.invalid";
  APIErrorCodes2["NetworkError"] = "error.network";
})(exports.APIErrorCodes || (exports.APIErrorCodes = {}));
exports.TFAType = void 0;
(function(TFAType2) {
  TFAType2["TFATypeApp"] = "app";
  TFAType2["TFATypeSMS"] = "sms";
  TFAType2["TFATypeEmail"] = "email";
})(exports.TFAType || (exports.TFAType = {}));
exports.TFAStatus = void 0;
(function(TFAStatus2) {
  TFAStatus2["DISABLED"] = "disabled";
  TFAStatus2["OPTIONAL"] = "optional";
  TFAStatus2["MANDATORY"] = "mandatory";
})(exports.TFAStatus || (exports.TFAStatus = {}));
class ApiError extends Error {
  constructor(error) {
    super((error == null ? void 0 : error.message) || "Unknown API error");
    this.detailedMessage = error == null ? void 0 : error.detailed_message;
    this.id = error == null ? void 0 : error.id;
    this.status = error == null ? void 0 : error.status;
  }
}

var __defProp$2 = Object.defineProperty;
var __getOwnPropSymbols$2 = Object.getOwnPropertySymbols;
var __hasOwnProp$2 = Object.prototype.hasOwnProperty;
var __propIsEnum$2 = Object.prototype.propertyIsEnumerable;
var __defNormalProp$2 = (obj, key, value) => key in obj ? __defProp$2(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __spreadValues$2 = (a, b) => {
  for (var prop in b || (b = {}))
    if (__hasOwnProp$2.call(b, prop))
      __defNormalProp$2(a, prop, b[prop]);
  if (__getOwnPropSymbols$2)
    for (var prop of __getOwnPropSymbols$2(b)) {
      if (__propIsEnum$2.call(b, prop))
        __defNormalProp$2(a, prop, b[prop]);
    }
  return a;
};
var __async$3 = (__this, __arguments, generator) => {
  return new Promise((resolve, reject) => {
    var fulfilled = (value) => {
      try {
        step(generator.next(value));
      } catch (e) {
        reject(e);
      }
    };
    var rejected = (value) => {
      try {
        step(generator.throw(value));
      } catch (e) {
        reject(e);
      }
    };
    var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
    step((generator = generator.apply(__this, __arguments)).next());
  });
};
const APP_ID_HEADER_KEY = "X-Identifo-Clientid";
const AUTHORIZATION_HEADER_KEY = "Authorization";
class API {
  constructor(config, tokenService) {
    this.config = config;
    this.tokenService = tokenService;
    this.defaultHeaders = {
      [APP_ID_HEADER_KEY]: "",
      Accept: "application/json",
      "Content-Type": "application/json"
    };
    this.catchNetworkErrorHandler = (e) => {
      if (e.message === "Network Error" || e.message === "Failed to fetch" || e.message === "Preflight response is not successful" || e.message.indexOf("is not allowed by Access-Control-Allow-Origin") > -1) {
        console.error(e.message);
        throw new ApiError({
          id: exports.APIErrorCodes.NetworkError,
          status: 0,
          message: "Configuration error",
          detailed_message: `Please check Identifo URL and add "${window.location.protocol}//${window.location.host}" to "REDIRECT URLS" in Identifo app settings.`
        });
      }
      throw e;
    };
    this.checkStatusCodeAndGetJSON = (r) => __async$3(this, null, function* () {
      if (!r.ok) {
        const error = yield r.json();
        throw new ApiError(error == null ? void 0 : error.error);
      }
      return r.json();
    });
    this.baseUrl = config.url.replace(/\/$/, "");
    this.defaultHeaders[APP_ID_HEADER_KEY] = config.appId;
    this.appId = config.appId;
  }
  get(path, options) {
    return this.send(path, __spreadValues$2({ method: "GET" }, options));
  }
  put(path, data, options) {
    return this.send(path, __spreadValues$2({ method: "PUT", body: JSON.stringify(data) }, options));
  }
  post(path, data, options) {
    return this.send(path, __spreadValues$2({ method: "POST", body: JSON.stringify(data) }, options));
  }
  send(path, options) {
    const init = __spreadValues$2({}, options);
    init.credentials = "include";
    init.headers = __spreadValues$2(__spreadValues$2({}, init.headers), this.defaultHeaders);
    return fetch(`${this.baseUrl}${path}`, init).catch(this.catchNetworkErrorHandler).then(this.checkStatusCodeAndGetJSON).then((value) => value);
  }
  getUser() {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.get("/me", {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}`
        }
      });
    });
  }
  renewToken() {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken("refresh")) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.post("/auth/token", { scopes: this.config.scopes }, {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${(_b = this.tokenService.getToken("refresh")) == null ? void 0 : _b.token}`
        }
      }).then((r) => this.storeToken(r));
    });
  }
  updateUser(user) {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.put("/me", user, {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${(_b = this.tokenService.getToken("access")) == null ? void 0 : _b.token}`
        }
      });
    });
  }
  login(email, password, deviceToken, scopes) {
    return __async$3(this, null, function* () {
      const data = {
        email,
        password,
        device_token: deviceToken,
        scopes
      };
      return this.post("/auth/login", data).then((r) => this.storeToken(r));
    });
  }
  requestPhoneCode(phone) {
    return __async$3(this, null, function* () {
      const data = {
        phone_number: phone
      };
      return this.post("/auth/request_phone_code", data);
    });
  }
  phoneLogin(phone, code, scopes) {
    return __async$3(this, null, function* () {
      const data = {
        phone_number: phone,
        code,
        scopes
      };
      return this.post("/auth/phone_login", data).then((r) => this.storeToken(r));
    });
  }
  federatedLogin(_0, _1, _2, _3) {
    return __async$3(this, arguments, function* (provider, scopes, redirectUrl, callbackUrl, opts = { width: 600, height: 800, popUp: false }) {
      const dataForm = document.createElement("form");
      dataForm.style.display = "none";
      if (opts.popUp) {
        dataForm.target = "TargetWindow";
      }
      dataForm.method = "POST";
      const params = new URLSearchParams();
      params.set("appId", this.config.appId);
      params.set("provider", provider);
      params.set("scopes", scopes.join(","));
      params.set("redirectUrl", redirectUrl);
      if (callbackUrl) {
        params.set("callbackUrl", callbackUrl);
      }
      dataForm.action = `${this.baseUrl}/auth/federated?${params.toString()}`;
      document.body.appendChild(dataForm);
      if (opts.popUp) {
        const left = window.screenX + window.outerWidth / 2 - (opts.width || 600) / 2;
        const top = window.screenY + window.outerHeight / 2 - (opts.height || 800) / 2;
        const postWindow = window.open("", "TargetWindow", `status=0,title=0,height=${opts.height},width=${opts.width},top=${top},left=${left},scrollbars=1`);
        if (postWindow) {
          dataForm.submit();
        }
      } else {
        window.location.assign(`${this.baseUrl}/auth/federated?${params.toString()}`);
      }
    });
  }
  federatedLoginComplete(params) {
    return __async$3(this, null, function* () {
      return this.get(`/auth/federated/complete?${params.toString()}`).then((r) => this.storeToken(r));
    });
  }
  register(email, password, scopes, invite) {
    return __async$3(this, null, function* () {
      const data = {
        email,
        password,
        scopes
      };
      if (invite) {
        data.invite = invite;
      }
      return this.post("/auth/register", data).then((r) => this.storeToken(r));
    });
  }
  requestResetPassword(email, tfaCode) {
    return __async$3(this, null, function* () {
      const data = {
        email,
        tfa_code: tfaCode
      };
      return this.post("/auth/request_reset_password", data);
    });
  }
  resetPassword(password) {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      const data = {
        password
      };
      return this.post("/auth/reset_password", data, {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}`
        }
      });
    });
  }
  getAppSettings(callbackUrl) {
    return __async$3(this, null, function* () {
      return this.get(`/auth/app_settings?${new URLSearchParams({ callbackUrl }).toString()}`);
    });
  }
  enableTFA(data) {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.put("/auth/tfa/enable", data, {
        headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}` }
      }).then((r) => this.storeToken(r));
    });
  }
  verifyTFA(code, scopes) {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.post("/auth/tfa/login", { tfa_code: code, scopes }, { headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}` } }).then((r) => this.storeToken(r));
    });
  }
  resendTFA() {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.post("/auth/tfa/resend", null, {
        headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}` }
      }).then((r) => this.storeToken(r));
    });
  }
  logout() {
    return __async$3(this, null, function* () {
      var _a, _b, _c;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.post("/me/logout", {
        refresh_token: (_b = this.tokenService.getToken("refresh")) == null ? void 0 : _b.token
      }, {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${(_c = this.tokenService.getToken()) == null ? void 0 : _c.token}`
        }
      }).then((r) => {
        this.tokenService.removeToken();
        this.tokenService.removeToken("refresh");
        return r;
      });
    });
  }
  invite(email, role, callbackUrl) {
    return __async$3(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.post("/auth/invite", {
        email,
        access_role: role,
        callback_url: callbackUrl
      }, {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}`
        }
      });
    });
  }
  storeToken(response) {
    if (response.access_token) {
      this.tokenService.saveToken(response.access_token, "access");
    }
    if (response.refresh_token) {
      this.tokenService.saveToken(response.refresh_token, "refresh");
    }
    return response;
  }
}

const jwtRegex = /^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-=]*$)/;
const INVALID_TOKEN_ERROR = "Empty or invalid token";
const TOKEN_QUERY_KEY = "token";
const REFRESH_TOKEN_QUERY_KEY = "refresh_token";

class CookieStorage {
  constructor() {
    this.isAccessible = false;
  }
  saveToken() {
    return true;
  }
  getToken() {
    throw new Error("Can not get token from HttpOnly");
  }
  deleteToken() {
  }
}

class StorageManager {
  constructor(storageType, accessKey, refreshKey) {
    this.preffix = "identifo_";
    this.storageType = "localStorage";
    this.access = `${this.preffix}access_token`;
    this.refresh = `${this.preffix}refresh_token`;
    this.isAccessible = true;
    this.access = accessKey ? this.preffix + accessKey : this.access;
    this.refresh = refreshKey ? this.preffix + refreshKey : this.refresh;
    this.storageType = storageType;
  }
  saveToken(token, tokenType) {
    if (token) {
      window[this.storageType].setItem(this[tokenType], token);
      return true;
    }
    return false;
  }
  getToken(tokenType) {
    var _a;
    return (_a = window[this.storageType].getItem(this[tokenType])) != null ? _a : "";
  }
  deleteToken(tokenType) {
    window[this.storageType].removeItem(this[tokenType]);
  }
}

class LocalStorage extends StorageManager {
  constructor(accessKey, refreshKey) {
    super("localStorage", accessKey, refreshKey);
  }
}

class SessionStorage extends StorageManager {
  constructor(accessKey, refreshKey) {
    super("sessionStorage", accessKey, refreshKey);
  }
}

var __async$2 = (__this, __arguments, generator) => {
  return new Promise((resolve, reject) => {
    var fulfilled = (value) => {
      try {
        step(generator.next(value));
      } catch (e) {
        reject(e);
      }
    };
    var rejected = (value) => {
      try {
        step(generator.throw(value));
      } catch (e) {
        reject(e);
      }
    };
    var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
    step((generator = generator.apply(__this, __arguments)).next());
  });
};
class TokenService {
  constructor(tokenManager) {
    this.isAuth = false;
    this.tokenManager = tokenManager || new LocalStorage();
  }
  handleVerification(token, audience, issuer) {
    return __async$2(this, null, function* () {
      if (!this.tokenManager.isAccessible)
        return true;
      try {
        yield this.validateToken(token, audience, issuer);
        this.saveToken(token);
        return true;
      } catch (err) {
        this.removeToken();
        return Promise.reject(err);
      }
    });
  }
  validateToken(token, audience, issuer) {
    return __async$2(this, null, function* () {
      var _a;
      if (!token)
        throw new Error(INVALID_TOKEN_ERROR);
      const jwtPayload = this.parseJWT(token);
      const isJwtExpired = this.isJWTExpired(jwtPayload);
      if (((_a = jwtPayload.aud) == null ? void 0 : _a.includes(audience)) && (!issuer || jwtPayload.iss === issuer) && !isJwtExpired) {
        return Promise.resolve(true);
      }
      throw new Error(INVALID_TOKEN_ERROR);
    });
  }
  parseJWT(token) {
    const base64Url = token.split(".")[1];
    if (!base64Url)
      return { aud: [], iss: "", exp: 10 };
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(atob(base64).split("").map((c) => `%${`00${c.charCodeAt(0).toString(16)}`.slice(-2)}`).join(""));
    return JSON.parse(jsonPayload);
  }
  isJWTExpired(token) {
    const now = new Date().getTime() / 1e3;
    if (token.exp && now > token.exp) {
      return true;
    }
    return false;
  }
  saveToken(token, type = "access") {
    if (type === "access") {
      this.isAuth = true;
    }
    return this.tokenManager.saveToken(token, type);
  }
  removeToken(type = "access") {
    if (type === "access") {
      this.isAuth = false;
    }
    this.tokenManager.deleteToken(type);
  }
  getToken(type = "access") {
    const token = this.tokenManager.getToken(type);
    if (!token)
      return null;
    const jwtPayload = this.parseJWT(token);
    return { token, payload: jwtPayload };
  }
}

class UrlBuilder {
  constructor(config) {
    this.config = config;
  }
  getUrl(flow) {
    var _a, _b;
    const scopes = ((_a = this.config.scopes) == null ? void 0 : _a.join()) || "";
    const redirectUri = (_b = this.config.redirectUri) != null ? _b : window.location.href;
    const baseParams = `appId=${this.config.appId}&scopes=${scopes}`;
    const urlParams = `${baseParams}&callbackUrl=${encodeURIComponent(redirectUri)}`;
    const postLogoutRedirectUri = this.config.postLogoutRedirectUri ? `${this.config.postLogoutRedirectUri}` : `${redirectUri}&redirectUri=${this.config.url}/web/login?${encodeURIComponent(baseParams)}`;
    const urls = {
      signup: `${this.config.url}/web/register?${urlParams}`,
      signin: `${this.config.url}/web/login?${urlParams}`,
      logout: `${this.config.url}/web/logout?${baseParams}&callbackUrl=${encodeURIComponent(postLogoutRedirectUri)}`,
      renew: `${this.config.url}/web/token/renew?${baseParams}&redirectUri=${encodeURIComponent(redirectUri)}`,
      default: "default"
    };
    return urls[flow] || urls.default;
  }
  createSignupUrl() {
    return this.getUrl("signup");
  }
  createSigninUrl() {
    return this.getUrl("signin");
  }
  createLogoutUrl() {
    return this.getUrl("logout");
  }
  createRenewSessionUrl() {
    return this.getUrl("renew");
  }
}

var __defProp$1 = Object.defineProperty;
var __defProps$1 = Object.defineProperties;
var __getOwnPropDescs$1 = Object.getOwnPropertyDescriptors;
var __getOwnPropSymbols$1 = Object.getOwnPropertySymbols;
var __hasOwnProp$1 = Object.prototype.hasOwnProperty;
var __propIsEnum$1 = Object.prototype.propertyIsEnumerable;
var __defNormalProp$1 = (obj, key, value) => key in obj ? __defProp$1(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __spreadValues$1 = (a, b) => {
  for (var prop in b || (b = {}))
    if (__hasOwnProp$1.call(b, prop))
      __defNormalProp$1(a, prop, b[prop]);
  if (__getOwnPropSymbols$1)
    for (var prop of __getOwnPropSymbols$1(b)) {
      if (__propIsEnum$1.call(b, prop))
        __defNormalProp$1(a, prop, b[prop]);
    }
  return a;
};
var __spreadProps$1 = (a, b) => __defProps$1(a, __getOwnPropDescs$1(b));
var __async$1 = (__this, __arguments, generator) => {
  return new Promise((resolve, reject) => {
    var fulfilled = (value) => {
      try {
        step(generator.next(value));
      } catch (e) {
        reject(e);
      }
    };
    var rejected = (value) => {
      try {
        step(generator.throw(value));
      } catch (e) {
        reject(e);
      }
    };
    var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
    step((generator = generator.apply(__this, __arguments)).next());
  });
};
class IdentifoAuth {
  constructor(config) {
    this.token = null;
    if (config) {
      this.configure(config);
    }
  }
  get isAuth() {
    var _a;
    return !!((_a = this.tokenService) == null ? void 0 : _a.isAuth);
  }
  configure(config) {
    var _a, _b;
    this.config = __spreadProps$1(__spreadValues$1({}, config), { autoRenew: (_a = config.autoRenew) != null ? _a : true });
    this.tokenService = new TokenService(config.tokenManager);
    this.urlBuilder = new UrlBuilder(this.config);
    this.api = new API(config, this.tokenService);
    this.handleToken(((_b = this.tokenService.getToken()) == null ? void 0 : _b.token) || "", "access");
  }
  handleToken(token, tokenType) {
    if (token) {
      if (tokenType === "access") {
        const payload = this.tokenService.parseJWT(token);
        this.token = { token, payload };
        this.tokenService.saveToken(token);
      } else {
        this.tokenService.saveToken(token, "refresh");
      }
    }
  }
  resetAuthValues() {
    this.token = null;
    this.tokenService.removeToken();
    this.tokenService.removeToken("refresh");
  }
  signup() {
    window.location.href = this.urlBuilder.createSignupUrl();
  }
  signin() {
    window.location.href = this.urlBuilder.createSigninUrl();
  }
  logout() {
    this.resetAuthValues();
    window.location.href = this.urlBuilder.createLogoutUrl();
  }
  handleAuthentication() {
    return __async$1(this, null, function* () {
      const { access, refresh } = this.getTokenFromUrl();
      if (!access) {
        this.resetAuthValues();
        return Promise.reject();
      }
      try {
        yield this.tokenService.handleVerification(access, this.config.appId, this.config.issuer);
        this.handleToken(access, "access");
        if (refresh) {
          this.handleToken(refresh, "refresh");
        }
        return yield Promise.resolve(true);
      } catch (err) {
        this.resetAuthValues();
        return yield Promise.reject();
      } finally {
        window.history.pushState({}, document.title, window.location.pathname);
      }
    });
  }
  getTokenFromUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    const tokens = { access: "", refresh: "" };
    const accessToken = urlParams.get(TOKEN_QUERY_KEY);
    const refreshToken = urlParams.get(REFRESH_TOKEN_QUERY_KEY);
    if (refreshToken && jwtRegex.test(refreshToken)) {
      tokens.refresh = refreshToken;
    }
    if (accessToken && jwtRegex.test(accessToken)) {
      tokens.access = accessToken;
    }
    return tokens;
  }
  getToken() {
    return __async$1(this, null, function* () {
      const token = this.tokenService.getToken();
      const refreshToken = this.tokenService.getToken("refresh");
      if (token) {
        const isExpired = this.tokenService.isJWTExpired(token.payload);
        if (isExpired && refreshToken) {
          try {
            yield this.renewSession();
            return yield Promise.resolve(this.token);
          } catch (err) {
            this.resetAuthValues();
            throw new Error("No token");
          }
        }
        return Promise.resolve(token);
      }
      return Promise.resolve(null);
    });
  }
  renewSession() {
    return __async$1(this, null, function* () {
      try {
        const { access, refresh } = yield this.renewSessionWithToken();
        this.handleToken(access, "access");
        this.handleToken(refresh, "refresh");
        return yield Promise.resolve(access);
      } catch (err) {
        return Promise.reject();
      }
    });
  }
  renewSessionWithToken() {
    return __async$1(this, null, function* () {
      try {
        const tokens = yield this.api.renewToken().then((l) => ({ access: l.access_token || "", refresh: l.refresh_token || "" }));
        return tokens;
      } catch (err) {
        return Promise.reject(err);
      }
    });
  }
}

exports.Routes = void 0;
(function(Routes2) {
  Routes2["LOGIN"] = "login";
  Routes2["REGISTER"] = "register";
  Routes2["TFA_VERIFY_SMS"] = "tfa/verify/sms";
  Routes2["TFA_VERIFY_EMAIL"] = "tfa/verify/email";
  Routes2["TFA_VERIFY_APP"] = "tfa/verify/app";
  Routes2["TFA_VERIFY_SELECT"] = "tfa/verify/select";
  Routes2["TFA_SETUP_SMS"] = "tfa/setup/sms";
  Routes2["TFA_SETUP_EMAIL"] = "tfa/setup/email";
  Routes2["TFA_SETUP_APP"] = "tfa/setup/app";
  Routes2["TFA_SETUP_SELECT"] = "tfa/setup/select";
  Routes2["PASSWORD_RESET"] = "password/reset";
  Routes2["PASSWORD_FORGOT"] = "password/forgot";
  Routes2["PASSWORD_FORGOT_TFA_SMS"] = "password/forgot/tfa/sms";
  Routes2["PASSWORD_FORGOT_TFA_EMAIL"] = "password/forgot/tfa/email";
  Routes2["PASSWORD_FORGOT_TFA_APP"] = "password/forgot/tfa/app";
  Routes2["PASSWORD_FORGOT_TFA_SELECT"] = "password/forgot/tfa/select";
  Routes2["CALLBACK"] = "callback";
  Routes2["LOGIN_PHONE"] = "login_phone";
  Routes2["LOGIN_PHONE_VERIFY"] = "login_phone_verify";
  Routes2["ERROR"] = "error";
  Routes2["PASSWORD_FORGOT_SUCCESS"] = "password/forgot/success";
  Routes2["LOGOUT"] = "logout";
  Routes2["LOADING"] = "loading";
})(exports.Routes || (exports.Routes = {}));
const typeToSetupRoute = {
  [exports.TFAType.TFATypeApp]: exports.Routes.TFA_SETUP_APP,
  [exports.TFAType.TFATypeEmail]: exports.Routes.TFA_SETUP_EMAIL,
  [exports.TFAType.TFATypeSMS]: exports.Routes.TFA_SETUP_SMS
};
const typeToTFAVerifyRoute = {
  [exports.TFAType.TFATypeApp]: exports.Routes.TFA_VERIFY_APP,
  [exports.TFAType.TFATypeEmail]: exports.Routes.TFA_VERIFY_EMAIL,
  [exports.TFAType.TFATypeSMS]: exports.Routes.TFA_VERIFY_SMS
};
const typeToPasswordForgotTFAVerifyRoute = {
  [exports.TFAType.TFATypeApp]: exports.Routes.PASSWORD_FORGOT_TFA_APP,
  [exports.TFAType.TFATypeEmail]: exports.Routes.PASSWORD_FORGOT_TFA_EMAIL,
  [exports.TFAType.TFATypeSMS]: exports.Routes.PASSWORD_FORGOT_TFA_SMS
};

var __defProp = Object.defineProperty;
var __defProps = Object.defineProperties;
var __getOwnPropDescs = Object.getOwnPropertyDescriptors;
var __getOwnPropSymbols = Object.getOwnPropertySymbols;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __propIsEnum = Object.prototype.propertyIsEnumerable;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __spreadValues = (a, b) => {
  for (var prop in b || (b = {}))
    if (__hasOwnProp.call(b, prop))
      __defNormalProp(a, prop, b[prop]);
  if (__getOwnPropSymbols)
    for (var prop of __getOwnPropSymbols(b)) {
      if (__propIsEnum.call(b, prop))
        __defNormalProp(a, prop, b[prop]);
    }
  return a;
};
var __spreadProps = (a, b) => __defProps(a, __getOwnPropDescs(b));
var __async = (__this, __arguments, generator) => {
  return new Promise((resolve, reject) => {
    var fulfilled = (value) => {
      try {
        step(generator.next(value));
      } catch (e) {
        reject(e);
      }
    };
    var rejected = (value) => {
      try {
        step(generator.throw(value));
      } catch (e) {
        reject(e);
      }
    };
    var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
    step((generator = generator.apply(__this, __arguments)).next());
  });
};
const emailRegex = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
const phoneRegex = /^[\+][0-9]{9,15}$/;
class CDK {
  constructor() {
    this.scopes = new Set();
    this.state = new rxjs.BehaviorSubject({ route: exports.Routes.LOADING });
    this.afterLoginRedirect = (loginResponse) => __async(this, null, function* () {
      if (loginResponse.require_2fa) {
        if (!loginResponse.enabled_2fa) {
          yield this.redirectTfaSetup(loginResponse);
          return;
        }
        if (loginResponse.enabled_2fa) {
          yield this.redirectTfaVerify(loginResponse);
          return;
        }
      }
      if (this.settings.tfaStatus === exports.TFAStatus.OPTIONAL && [exports.Routes.LOGIN, exports.Routes.REGISTER].includes(this.state.getValue().route)) {
        this.tfaSetupSelect(loginResponse);
        return;
      }
      if (loginResponse.access_token && loginResponse.refresh_token) {
        this.callback(loginResponse);
        return;
      }
      if (loginResponse.access_token && !loginResponse.refresh_token) {
        this.callback(loginResponse);
        return;
      }
      this.login();
    });
    this.loginCatchRedirect = (data) => {
      if (data.id === exports.APIErrorCodes.PleaseEnableTFA) {
        return;
      }
      throw data;
    };
    this.auth = new IdentifoAuth();
  }
  configure(authConfig, callbackUrl) {
    return __async(this, null, function* () {
      var _a;
      this.state.next({ route: exports.Routes.LOADING });
      this.callbackUrl = callbackUrl;
      this.scopes = new Set((_a = authConfig.scopes) != null ? _a : []);
      this.postLogoutRedirectUri = window.location.origin + window.location.pathname;
      if (!authConfig.appId) {
        this.state.next({
          route: exports.Routes.ERROR,
          error: { message: "app-id param is empty", name: "app-id empty" }
        });
        return;
      }
      if (!authConfig.url) {
        this.state.next({
          route: exports.Routes.ERROR,
          error: { message: "url param is empty", name: "url empty" }
        });
        return;
      }
      this.auth.configure(authConfig);
      try {
        this.settings = yield this.auth.api.getAppSettings(callbackUrl);
      } catch (err) {
        this.state.next({
          route: exports.Routes.ERROR,
          error: err
        });
        return;
      }
      this.settings.tfaType = Array.isArray(this.settings.tfaType) ? this.settings.tfaType : [this.settings.tfaType];
      const href = new URL(window.location.href);
      if (!!href.searchParams.get("provider") && !!href.searchParams.get("state")) {
        const u = new URL(window.location.href);
        const sp = new URLSearchParams();
        const appId = href.searchParams.get("appId");
        if (appId === null) {
          this.state.next({
            route: exports.Routes.ERROR,
            error: { message: "app-id param is empty", name: "app-id empty" }
          });
          return;
        }
        sp.set("appId", appId);
        window.history.replaceState({}, document.title, `${u.pathname}?${sp.toString()}`);
        this.auth.api.federatedLoginComplete(u.searchParams).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).catch((e) => this.processError(e));
      }
    });
  }
  login() {
    switch (true) {
      case (!this.auth.config.loginWith && this.settings.loginWith["phone"] || this.auth.config.loginWith === "phone" && this.settings.loginWith["phone"]):
        return this.loginWithPhone();
      case (!this.auth.config.loginWith && this.settings.loginWith["email"] || this.auth.config.loginWith === "email" && this.settings.loginWith["email"]):
        return this.loginWithPassword();
      default:
        throw "Unsupported login way";
    }
  }
  loginWithPhone() {
    var _a, _b;
    this.state.next({
      route: exports.Routes.LOGIN_PHONE,
      registrationForbidden: (_a = this.settings) == null ? void 0 : _a.registrationForbidden,
      error: this.lastError,
      federatedProviders: (_b = this.settings) == null ? void 0 : _b.federatedProviders,
      loginTypes: this.getLoginTypes("phone"),
      requestCode: (phone, remember) => __async(this, null, function* () {
        if (!this.validatePhone(phone)) {
          return;
        }
        const scopes = new Set(this.scopes);
        if (remember) {
          scopes.add("offline");
        }
        yield this.auth.api.requestPhoneCode(phone).then(() => this.loginWithPhoneVerify(phone, remember)).catch((e) => this.processError(e));
      }),
      socialLogin: (provider) => __async(this, null, function* () {
        this.state.next({ route: exports.Routes.LOADING });
        const federatedRedirectUrl = window.location.origin + window.location.pathname;
        return this.auth.api.federatedLogin(provider, [...this.scopes], federatedRedirectUrl, this.callbackUrl);
      })
    });
  }
  loginWithPhoneVerify(phone, remember) {
    this.state.next({
      route: exports.Routes.LOGIN_PHONE_VERIFY,
      error: this.lastError,
      phone,
      resendTimeout: this.settings.tfaResendTimeout * 1e3,
      resendCode: () => __async(this, null, function* () {
        yield this.auth.api.requestPhoneCode(phone);
      }),
      login: (code) => __async(this, null, function* () {
        const scopes = new Set(this.scopes);
        if (remember) {
          scopes.add("offline");
        }
        yield this.auth.api.phoneLogin(phone, code, [...this.scopes]).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).catch((e) => this.processError(e));
      }),
      goback: () => __async(this, null, function* () {
        this.login();
      })
    });
  }
  loginWithPassword() {
    var _a, _b;
    this.state.next({
      route: exports.Routes.LOGIN,
      registrationForbidden: (_a = this.settings) == null ? void 0 : _a.registrationForbidden,
      error: this.lastError,
      federatedProviders: (_b = this.settings) == null ? void 0 : _b.federatedProviders,
      loginTypes: this.getLoginTypes("email"),
      signup: () => __async(this, null, function* () {
        this.register();
      }),
      signin: (email, password, remember) => __async(this, null, function* () {
        if (!this.validateEmail(email)) {
          return;
        }
        const scopes = new Set(this.scopes);
        if (remember) {
          scopes.add("offline");
        }
        yield this.auth.api.login(email, password, "", [...Array.from(scopes)]).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).catch((e) => this.processError(e));
      }),
      socialLogin: (provider) => __async(this, null, function* () {
        this.state.next({ route: exports.Routes.LOADING });
        const federatedRedirectUrl = window.location.origin + window.location.pathname;
        return this.auth.api.federatedLogin(provider, [...this.scopes], federatedRedirectUrl, this.callbackUrl);
      }),
      passwordForgot: () => __async(this, null, function* () {
        this.forgotPassword();
      })
    });
  }
  register() {
    this.state.next({
      route: exports.Routes.REGISTER,
      signup: (email, password, token) => __async(this, null, function* () {
        if (!this.validateEmail(email)) {
          return;
        }
        yield this.auth.api.register(email, password, [...this.scopes], token).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).catch((e) => this.processError(e));
      }),
      goback: () => __async(this, null, function* () {
        this.login();
      })
    });
  }
  forgotPassword() {
    this.state.next({
      route: exports.Routes.PASSWORD_FORGOT,
      restorePassword: (email) => __async(this, null, function* () {
        return this.auth.api.requestResetPassword(email).then((response) => __async(this, null, function* () {
          if (response.result === "tfa-required") {
            yield this.redirectTfaForgot(email);
            return;
          }
          if (response.result === "ok") {
            this.forgotPasswordSuccess();
          }
        })).catch((e) => this.processError(e));
      }),
      goback: () => __async(this, null, function* () {
        this.login();
      })
    });
  }
  forgotPasswordSuccess() {
    this.state.next({
      route: exports.Routes.PASSWORD_FORGOT_SUCCESS,
      goback: () => __async(this, null, function* () {
        this.login();
      })
    });
  }
  passwordReset() {
    this.state.next({
      route: exports.Routes.PASSWORD_RESET,
      setNewPassword: (password) => __async(this, null, function* () {
        this.auth.api.resetPassword(password).then(() => {
          this.login();
        }).catch((e) => this.processError(e));
      })
    });
  }
  callback(result) {
    this.state.next({
      route: exports.Routes.CALLBACK,
      callbackUrl: this.callbackUrl,
      result
    });
    if (this.callbackUrl) {
      const url = new URL(this.callbackUrl);
      if (result.access_token) {
        url.searchParams.set("token", result.access_token);
      }
      if (result.refresh_token) {
        url.searchParams.set("refresh_token", result.refresh_token);
      }
      window.location.href = url.toString();
    }
  }
  validateEmail(email) {
    if (!emailRegex.test(email)) {
      this.processError({
        detailedMessage: "Email address is not valid",
        name: "Validation error",
        message: "Email address is not valid"
      });
      return false;
    }
    return true;
  }
  validatePhone(email) {
    if (!phoneRegex.test(email)) {
      this.processError({
        detailedMessage: "Phone is not valid",
        name: "Validation error",
        message: "Phone is not valid"
      });
      return false;
    }
    return true;
  }
  tfaSetup(loginResponse, type) {
    return __async(this, null, function* () {
      switch (type) {
        case exports.TFAType.TFATypeApp: {
          this.state.next({
            route: exports.Routes.TFA_SETUP_APP,
            provisioningURI: "",
            provisioningQR: "",
            setupTFA: () => __async(this, null, function* () {
            })
          });
          const tfa = yield this.auth.api.enableTFA({});
          if (tfa.provisioning_uri) {
            this.state.next({
              route: exports.Routes.TFA_SETUP_APP,
              provisioningURI: tfa.provisioning_uri,
              provisioningQR: tfa.provisioning_qr || "",
              setupTFA: () => __async(this, null, function* () {
                return this.tfaVerify(loginResponse, type);
              })
            });
          }
          break;
        }
        case exports.TFAType.TFATypeEmail: {
          this.state.next({
            route: exports.Routes.TFA_SETUP_EMAIL,
            email: loginResponse.user.email || "",
            setupTFA: (email) => __async(this, null, function* () {
              yield this.auth.api.enableTFA({ email });
              return this.tfaVerify(__spreadProps(__spreadValues({}, loginResponse), { user: __spreadProps(__spreadValues({}, loginResponse.user), { email }) }), type);
            })
          });
          break;
        }
        case exports.TFAType.TFATypeSMS: {
          this.state.next({
            route: exports.Routes.TFA_SETUP_SMS,
            phone: loginResponse.user.phone || "",
            setupTFA: (phone) => __async(this, null, function* () {
              yield this.auth.api.enableTFA({ phone });
              return this.tfaVerify(__spreadProps(__spreadValues({}, loginResponse), { user: __spreadProps(__spreadValues({}, loginResponse.user), { phone }) }), type);
            })
          });
          break;
        }
      }
    });
  }
  tfaVerify(loginResponse, type) {
    return __async(this, null, function* () {
      const state = {
        route: typeToTFAVerifyRoute[type],
        email: loginResponse.user.email,
        phone: loginResponse.user.phone,
        verifyTFA: (code) => __async(this, null, function* () {
          yield this.auth.api.verifyTFA(code, [...this.scopes]).then(this.afterLoginRedirect).catch(this.loginCatchRedirect).catch((e) => this.processError(e));
        })
      };
      switch (type) {
        case exports.TFAType.TFATypeApp: {
          this.state.next(__spreadValues({}, state));
          break;
        }
        case exports.TFAType.TFATypeEmail:
        case exports.TFAType.TFATypeSMS: {
          this.state.next(__spreadProps(__spreadValues({}, state), {
            resendTimeout: this.settings.tfaResendTimeout * 1e3,
            resendTFA: () => __async(this, null, function* () {
              yield this.auth.api.resendTFA();
            })
          }));
          break;
        }
      }
    });
  }
  passwordForgotTFAVerify(email, type) {
    return __async(this, null, function* () {
      this.state.next({
        route: typeToPasswordForgotTFAVerifyRoute[type],
        verifyTFA: (code) => __async(this, null, function* () {
          this.auth.api.requestResetPassword(email, code).then(() => {
            this.forgotPasswordSuccess();
          }).catch((e) => this.processError(e));
        })
      });
    });
  }
  logout() {
    return __async(this, null, function* () {
      this.state.next({
        route: exports.Routes.LOGOUT,
        logout: () => __async(this, null, function* () {
          return this.auth.api.logout();
        })
      });
    });
  }
  processError(e) {
    var _a, _b;
    e.detailedMessage = (_a = e.detailedMessage) == null ? void 0 : _a.trim();
    e.message = (_b = e.message) == null ? void 0 : _b.trim();
    this.state.next(__spreadProps(__spreadValues({}, this.state.getValue()), {
      error: e
    }));
  }
  redirectTfaSetup(loginResponse) {
    return __async(this, null, function* () {
      if (this.settings.tfaType.length === 1) {
        yield this.tfaSetup(loginResponse, this.settings.tfaType[0]);
        return;
      }
      this.tfaSetupSelect(loginResponse);
    });
  }
  tfaSetupSelect(loginResponse) {
    this.state.next({
      route: exports.Routes.TFA_SETUP_SELECT,
      tfaStatus: this.settings.tfaStatus,
      tfaTypes: this.settings.tfaType,
      select: (type) => __async(this, null, function* () {
        yield this.tfaSetup(loginResponse, type);
      }),
      setupNextTime: () => {
        this.callback(loginResponse);
      }
    });
  }
  redirectTfaVerify(e) {
    return __async(this, null, function* () {
      if (this.settings.tfaType.length === 1) {
        yield this.tfaVerify(e, this.settings.tfaType[0]);
        return;
      }
      this.state.next({
        route: exports.Routes.TFA_VERIFY_SELECT,
        tfaStatus: this.settings.tfaStatus,
        tfaTypes: this.settings.tfaType,
        select: (type) => __async(this, null, function* () {
          yield this.tfaVerify(e, type);
        })
      });
    });
  }
  redirectTfaForgot(email) {
    return __async(this, null, function* () {
      if (this.settings.tfaType.length === 1) {
        yield this.passwordForgotTFAVerify(email, this.settings.tfaType[0]);
        return;
      }
      this.state.next({
        route: exports.Routes.PASSWORD_FORGOT_TFA_SELECT,
        tfaStatus: this.settings.tfaStatus,
        tfaTypes: this.settings.tfaType,
        select: (type) => __async(this, null, function* () {
          yield this.passwordForgotTFAVerify(email, type);
        })
      });
    });
  }
  getLoginTypes(current) {
    const result = {};
    Object.entries(this.settings.loginWith).filter((v) => v[1] && v[0] !== current).forEach((v) => {
      result[v[0]] = {
        type: v[0],
        click: () => {
          this.auth.config.loginWith = v[0];
          this.login();
        }
      };
    });
    return result;
  }
}

exports.ApiError = ApiError;
exports.CDK = CDK;
exports.CookieStorageManager = CookieStorage;
exports.IdentifoAuth = IdentifoAuth;
exports.LocalStorageManager = LocalStorage;
exports.SessionStorageManager = SessionStorage;
exports.typeToPasswordForgotTFAVerifyRoute = typeToPasswordForgotTFAVerifyRoute;
exports.typeToSetupRoute = typeToSetupRoute;
exports.typeToTFAVerifyRoute = typeToTFAVerifyRoute;
//# sourceMappingURL=identifo.js.map

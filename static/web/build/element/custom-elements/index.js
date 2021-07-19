import { attachShadow, createEvent, h, getAssetPath, proxyCustomElement } from '@stencil/core/internal/client';
export { setAssetPath, setPlatformOptions } from '@stencil/core/internal/client';

var APIErrorCodes;
(function(APIErrorCodes2) {
  APIErrorCodes2["PleaseEnableTFA"] = "error.api.request.2fa.please_enable";
  APIErrorCodes2["NetworkError"] = "error.network";
})(APIErrorCodes || (APIErrorCodes = {}));
var TFAType;
(function(TFAType2) {
  TFAType2["TFATypeApp"] = "app";
  TFAType2["TFATypeSMS"] = "sms";
  TFAType2["TFATypeEmail"] = "email";
})(TFAType || (TFAType = {}));
class ApiError extends Error {
  constructor(error) {
    super((error == null ? void 0 : error.message) || "Unknown API error");
    this.detailedMessage = error == null ? void 0 : error.detailed_message;
    this.id = error == null ? void 0 : error.id;
    this.status = error == null ? void 0 : error.status;
  }
}

var __defProp$1 = Object.defineProperty;
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
const APP_ID_HEADER_KEY = "X-Identifo-Clientid";
const AUTHORIZATION_HEADER_KEY = "Authorization";
class Api {
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
          id: APIErrorCodes.NetworkError,
          status: 0,
          message: "Configuration error",
          detailed_message: `Please check Identifo URL and add "${window.location.protocol}//${window.location.host}" to "REDIRECT URLS" in Identifo app settings.`
        });
      }
      throw e;
    };
    this.checkStatusCodeAndGetJSON = (r) => __async$2(this, null, function* () {
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
    return this.send(path, __spreadValues$1({ method: "GET" }, options));
  }
  put(path, data, options) {
    return this.send(path, __spreadValues$1({ method: "PUT", body: JSON.stringify(data) }, options));
  }
  post(path, data, options) {
    return this.send(path, __spreadValues$1({ method: "POST", body: JSON.stringify(data) }, options));
  }
  send(path, options) {
    const init = __spreadValues$1({}, options);
    init.headers = __spreadValues$1(__spreadValues$1({}, init.headers), this.defaultHeaders);
    return fetch(`${this.baseUrl}${path}`, init).catch(this.catchNetworkErrorHandler).then(this.checkStatusCodeAndGetJSON).then((value) => value);
  }
  getUser() {
    return __async$2(this, null, function* () {
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
    return __async$2(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken("refresh")) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.post("/auth/token", {}, {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${(_b = this.tokenService.getToken("refresh")) == null ? void 0 : _b.token}`
        }
      }).then((r) => this.storeToken(r));
    });
  }
  updateUser(user) {
    return __async$2(this, null, function* () {
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
    return __async$2(this, null, function* () {
      const data = {
        email,
        password,
        device_token: deviceToken,
        scopes
      };
      return this.post("/auth/login", data).then((r) => this.storeToken(r));
    });
  }
  register(email, password, scopes) {
    return __async$2(this, null, function* () {
      const data = {
        email,
        password,
        scopes
      };
      return this.post("/auth/register", data).then((r) => this.storeToken(r));
    });
  }
  requestResetPassword(email) {
    return __async$2(this, null, function* () {
      const data = {
        email
      };
      return this.post("/auth/request_reset_password", data);
    });
  }
  resetPassword(password) {
    return __async$2(this, null, function* () {
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
  getAppSettings() {
    return __async$2(this, null, function* () {
      return this.get("/auth/app_settings");
    });
  }
  enableTFA() {
    return __async$2(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.put("/auth/tfa/enable", {}, {
        headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}` }
      });
    });
  }
  verifyTFA(code, scopes) {
    return __async$2(this, null, function* () {
      var _a, _b;
      if (!((_a = this.tokenService.getToken()) == null ? void 0 : _a.token)) {
        throw new Error("No token in token service.");
      }
      return this.post("/auth/tfa/login", { tfa_code: code, scopes }, { headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${(_b = this.tokenService.getToken()) == null ? void 0 : _b.token}` } }).then((r) => this.storeToken(r));
    });
  }
  logout() {
    return __async$2(this, null, function* () {
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

const Iframe = {
  create() {
    const iframe = document.createElement("iframe");
    iframe.style.display = "none";
    document.body.appendChild(iframe);
    return iframe;
  },
  remove(iframe) {
    setTimeout(() => {
      if (document.body.contains(iframe)) {
        document.body.removeChild(iframe);
      }
    }, 0);
  },
  captureMessage(iframe, src) {
    return new Promise((resolve, reject) => {
      const handleMessage = (event) => {
        if (event.data.error)
          reject(event.data.error);
        resolve(event.data.accessToken);
        window.removeEventListener("message", handleMessage);
      };
      window.addEventListener("message", handleMessage, false);
      iframe.src = src;
    });
  }
};

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
class TokenService {
  constructor(tokenManager) {
    this.tokenManager = tokenManager || new LocalStorage();
  }
  handleVerification(token, audience, issuer) {
    return __async$1(this, null, function* () {
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
    return __async$1(this, null, function* () {
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
  isAuthenticated(audience, issuer) {
    if (!this.tokenManager.isAccessible)
      return Promise.resolve(true);
    const token = this.tokenManager.getToken("access");
    return this.validateToken(token, audience, issuer);
  }
  saveToken(token, type = "access") {
    return this.tokenManager.saveToken(token, type);
  }
  removeToken(type = "access") {
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
    const scopes = JSON.stringify((_a = this.config.scopes) != null ? _a : []);
    const redirectUri = encodeURIComponent((_b = this.config.redirectUri) != null ? _b : window.location.href);
    const baseParams = `appId=${this.config.appId}&scopes=${scopes}`;
    const urlParams = `${baseParams}&callbackUrl=${redirectUri}`;
    const postLogoutRedirectUri = this.config.postLogoutRedirectUri ? `&callbackUrl=${encodeURIComponent(this.config.postLogoutRedirectUri)}` : `&callbackUrl=${redirectUri}&redirectUri=${this.config.url}/web/login?${encodeURIComponent(baseParams)}`;
    const urls = {
      signup: `${this.config.url}/web/register?${urlParams}`,
      signin: `${this.config.url}/web/login?${urlParams}`,
      logout: `${this.config.url}/web/logout?${baseParams}${postLogoutRedirectUri}`,
      renew: `${this.config.url}/web/token/renew?${baseParams}&redirectUri=${redirectUri}`,
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
class IdentifoAuth {
  constructor(config) {
    this.token = null;
    this.refreshToken = null;
    this.isAuth = false;
    var _a;
    this.config = __spreadProps(__spreadValues({}, config), { autoRenew: (_a = config.autoRenew) != null ? _a : true });
    this.tokenService = new TokenService(config.tokenManager);
    this.urlBuilder = new UrlBuilder(this.config);
    this.api = new Api(config, this.tokenService);
  }
  init() {
    const token = this.tokenService.getToken();
    if (token) {
      const isExpired = this.tokenService.isJWTExpired(token.payload);
      if (isExpired) {
        this.renewSession().then((t) => this.handleToken(t)).catch(() => this.resetAuthValues());
      } else {
        this.handleToken(token.token);
      }
    }
  }
  handleToken(token) {
    const payload = this.tokenService.parseJWT(token);
    this.token = { token, payload };
    this.isAuth = true;
    this.tokenService.saveToken(token);
    if (this.renewSessionId) {
      window.clearTimeout(this.renewSessionId);
    }
    if (payload.exp) {
      this.renewSessionId = window.setTimeout(() => {
        if (this.config.autoRenew) {
          this.renewSession().then((t) => this.handleToken(t)).catch(() => this.resetAuthValues());
        } else {
          this.resetAuthValues();
        }
      }, (payload.exp - new Date().getTime() / 1e3 - 6e4) * 1e3);
    }
  }
  resetAuthValues() {
    this.token = null;
    this.isAuth = false;
    this.tokenService.removeToken();
  }
  signup() {
    window.location.href = this.urlBuilder.createSignupUrl();
  }
  signin() {
    window.location.href = this.urlBuilder.createSigninUrl();
  }
  logout() {
    this.tokenService.removeToken("access");
    window.location.href = this.urlBuilder.createLogoutUrl();
  }
  handleAuthentication() {
    return __async(this, null, function* () {
      const { access } = this.getTokenFromUrl();
      if (!access) {
        return Promise.reject();
      }
      try {
        yield this.tokenService.handleVerification(access, this.config.appId, this.config.issuer);
        this.handleToken(access);
        return yield Promise.resolve(true);
      } catch (err) {
        return yield Promise.reject();
      } finally {
        window.location.hash = "";
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
    const token = this.tokenService.getToken();
    if (token) {
      return token;
    }
    return { token: "", payload: {} };
  }
  renewSession() {
    return __async(this, null, function* () {
      try {
        const token = yield this.renewSessionWithIframe();
        this.handleToken(token);
        return yield Promise.resolve(token);
      } catch (err) {
        return Promise.reject();
      }
    });
  }
  renewSessionWithToken() {
    return __async(this, null, function* () {
      try {
        const r = yield this.api.renewToken().then((l) => l.access_token || "");
        return r;
      } catch (err) {
        return Promise.resolve("");
      }
    });
  }
  renewSessionWithIframe() {
    return __async(this, null, function* () {
      const iframe = Iframe.create();
      const timeout = setTimeout(() => {
        Iframe.remove(iframe);
        throw new Error("Timeout expired");
      }, 3e4);
      try {
        const token = yield Iframe.captureMessage(iframe, this.urlBuilder.createRenewSessionUrl());
        yield this.tokenService.handleVerification(token, this.config.appId, this.config.issuer);
        return token;
      } catch (err) {
        return Promise.reject(err);
      } finally {
        clearTimeout(timeout);
        Iframe.remove(iframe);
      }
    });
  }
}

const mainCss = ".wrapper,.wrapper-dark{--content-width:416px}.wrapper{--main-background:#f7f7f7;--blue-text:#6163f6;--field-background:#fff;--gray-line:#e0e0e0;--social-button:#1b1b1b;--text:#1b1b1b;--upload-photo:#e0e0e0;--content-width:416px}.wrapper-dark{--main-background:#1b1b1b;--blue-text:#8b8dfa;--field-background:#423f3f;--gray-line:#423f3f;--social-button:#423f3f;--text:#fff;--upload-photo:#423f3f;--content-width:416px}*{margin:0;padding:0;-webkit-box-sizing:border-box;box-sizing:border-box;font-family:inherit}.wrapper,.wrapper-dark{display:-ms-flexbox;display:flex;-ms-flex-pack:center;justify-content:center;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.social-buttons{width:100%;position:relative}.social-buttons__text{font-size:14px;line-height:21px;color:#828282;padding:4px 8px;margin-bottom:39px;text-align:center;position:static}.social-buttons__text::before{content:\"\";position:absolute;height:1px;width:36%;left:0;top:14px;background-color:var(--gray-line)}.social-buttons__text::after{content:\"\";position:absolute;height:1px;width:36%;right:0;top:14px;background-color:var(--gray-line)}.social-buttons__social-medias{display:-ms-flexbox;display:flex;-ms-flex-pack:center;justify-content:center;-ms-flex-align:center;align-items:center}.social-buttons__media{width:56px;height:56px;border-radius:50%;background-color:var(--social-button);display:-ms-flexbox;display:flex;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;cursor:pointer}.social-buttons__media:not(:last-of-type){margin-right:24px}@media (max-width: 599px){.social-buttons__media{width:44px;height:44px}.social-buttons__text{margin-bottom:36px}.social-buttons__text::before{width:26%}.social-buttons__text::after{width:26%}.social-buttons__image{width:16px;height:16px}}.primary-button{background-color:#6163f6;border:none;outline:none;display:-ms-flexbox;display:flex;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;width:192px;height:64px;border-radius:8px;cursor:pointer;color:#fff;font-size:18px;line-height:26px;-webkit-transition:all 0.4s;transition:all 0.4s}.primary-button:active{-webkit-transform:translateY(-4px);transform:translateY(-4px)}.primary-button:hover{opacity:0.8}.primary-button:disabled{cursor:initial;opacity:0.3}@media (max-width: 599px){.primary-button{width:100%}}.info-card{border:1px solid var(--gray-line);border-radius:8px;padding:24px}.info-card__controls{display:-ms-flexbox;display:flex;-ms-flex-pack:justify;justify-content:space-between}.info-card__title{color:var(--text);font-size:18px;line-height:26px;font-weight:700}.info-card__button{color:var(--blue-text);background:none;border:none;cursor:pointer;font-size:18px;line-height:26px}.info-card__text{color:#828282;font-size:16px;line-height:24px;margin-top:8px}.info-card__subtitle{color:var(--text);font-size:16px;line-height:24px;margin:4px 0 12px}@media (max-width: 599px){.info-card__text{font-size:14px;line-height:21px}}.form-control{width:100%;max-width:var(--content-width);height:72px;background-color:var(--field-background);-webkit-box-shadow:0px 11px 15px rgba(0, 0, 0, 0.04);box-shadow:0px 11px 15px rgba(0, 0, 0, 0.04);border-radius:8px;border:none;outline:none;font-size:18px;line-height:26px;color:var(--text);padding:23px 24px}.form-control::-webkit-inner-spin-button{-webkit-appearance:none;margin:0}.form-control::-webkit-input-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control::-moz-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control:-ms-input-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control::-ms-input-placeholder{font-size:18px;line-height:26px;color:#828282}.form-control::placeholder{font-size:18px;line-height:26px;color:#828282}.form-control-danger{border:1px solid #F66161}@media (max-width: 599px){.form-control{height:64px}}.upload-photo{display:-ms-flexbox;display:flex;-ms-flex-pack:center;justify-content:center;-ms-flex-align:center;align-items:center;margin-bottom:48px}.upload-photo__field{display:none}.upload-photo__label{cursor:pointer;color:var(--blue-text);font-size:16px;line-height:24px}.upload-photo__label:first-of-type{margin-right:16px}.upload-photo__avatar{height:64px;width:64px;border-radius:50%;background-color:var(--upload-photo)}@media (max-width: 599px){.upload-photo{margin-bottom:32px}}.error{min-height:21px;width:100%;font-size:14px;line-height:21px;color:#FF5160}.login-form{display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;width:var(--content-width)}.login-form__register-text{margin-bottom:32px;font-weight:400;font-size:14px;line-height:24px;color:#828282}.login-form__register-link{color:var(--blue-text);cursor:pointer}.login-form .form-control:first-of-type{margin-bottom:32px}.login-form__buttons{margin-top:48px;display:-ms-flexbox;display:flex;width:100%;-ms-flex-align:center;align-items:center;margin-bottom:36px}.login-form__buttons_mt-32{margin-top:32px}.login-form__forgot-pass{color:var(--blue-text);font-size:16px;line-height:24px;cursor:pointer}.login-form .primary-button{margin-right:32px}.login-form .error{margin-top:12px}@media (max-width: 599px){.login-form{width:100%;max-width:var(--content-width);padding:0 24px}.login-form__register-text{font-size:16px}.login-form .form-control:first-of-type{margin-bottom:24px}.login-form__buttons{margin-top:32px;-ms-flex-direction:column;flex-direction:column}.login-form .primary-button{margin-right:0;margin-bottom:36px}}.register-form{width:var(--content-width);padding:64px 0 44px}.register-form .form-control:not(:last-of-type){margin-bottom:32px}.register-form__buttons{display:-ms-flexbox;display:flex;width:100%;-ms-flex-align:center;align-items:center;margin-top:48px}.register-form__buttons_mt-32{margin-top:32px}.register-form__log-in{color:var(--blue-text);font-size:16px;line-height:24px;cursor:pointer}.register-form .primary-button{margin-right:32px}.register-form .error{margin-top:12px}@media (max-width: 599px){.register-form{width:100%;max-width:var(--content-width);padding:48px 24px 32px}.register-form .form-control{margin-bottom:24px}.register-form .primary-button{margin-right:0;margin-bottom:36px}.register-form__buttons{-ms-flex-direction:column;flex-direction:column;margin-top:32px}}.tfa-setup{padding:48px 0 80px;width:var(--content-width);display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.tfa-setup__text{font-size:24px;line-height:32px;color:var(--text);text-align:center;max-width:260px;width:100%}.tfa-setup .error{margin-top:12px}.tfa-setup__form{width:100%}.tfa-setup__form .primary-button{margin:48px auto 0}.tfa-setup__form .primary-button-mt-32{margin-top:32px}.tfa-setup__subtitle{text-align:center;font-size:16px;line-height:24px;color:#828282;max-width:270px;text-align:center;margin:16px auto 48px}.tfa-setup .info-card{margin-top:48px}@media (max-width: 599px){.tfa-setup{width:100%;max-width:var(--content-width);padding:38px 24px 36px}}.tfa-verify{padding:48px 0 52px;width:var(--content-width);display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.tfa-verify .error{margin-top:12px}.tfa-verify__title-wrapper{width:100%;display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.tfa-verify__title,.tfa-verify__title_mb-40{color:var(--text);font-size:24px;line-height:32px;max-width:280px;text-align:center;font-weight:400;margin-bottom:16px}.tfa-verify__title_mb-40{margin-bottom:40px}.tfa-verify__app-button{font-size:18px;line-height:26px;background:none;border:none;color:var(--blue-text);margin-bottom:40px}.tfa-verify .primary-button{margin-top:48px}.tfa-verify .primary-button-mt-32{margin-top:32px}.tfa-verify__back{font-size:16px;line-height:24px;color:var(--blue-text)}.tfa-verify__subtitle{font-size:16px;line-height:24px;margin-bottom:48px;color:#828282;max-width:189px;text-align:center}.tfa-verify__qr-code{width:160px;height:160px;margin-bottom:64px}@media (max-width: 599px){.tfa-verify{padding:102px 24px 41px;width:100%;max-width:var(--content-width)}.tfa-verify .primary-button{margin-top:32px}.tfa-verify__qr-code{margin-bottom:48px}}.otp-login{display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center;-ms-flex-pack:center;justify-content:center;padding:44px 24px 102px;width:100%;max-width:var(--content-width)}.otp-login__register-text{margin-bottom:32px;font-weight:400;font-size:14px;line-height:24px;color:#828282}.otp-login .form-control{margin-bottom:48px}.otp-login .primary-button{margin-bottom:36px}.error-view{width:100%;max-width:var(--content-width);padding:0 24px;display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center}.error-view__message{color:var(--text);font-size:24px;line-height:32px;text-align:center;margin-bottom:16px}.error-view__details{color:#828282;font-size:16px;line-height:24px;text-align:center}.error-view .primary-button{margin-top:64px}.forgot-password{width:var(--content-width)}.forgot-password__title{text-align:center;margin-bottom:16px;font-size:24px;line-height:32px;font-weight:400;color:var(--text)}.forgot-password__subtitle{text-align:center;font-size:16px;line-height:24px;color:#828282;max-width:189px;text-align:center;margin:0 auto 48px}.forgot-password .error{margin-top:12px}.forgot-password .primary-button{margin:48px auto 0}.forgot-password .primary-button-mt-32{margin-top:32px}@media (max-width: 599px){.forgot-password{width:100%;max-width:var(--content-width);padding:0 24px}.forgot-password .primary-button{margin-top:32px}}.forgot-password-success{display:-ms-flexbox;display:flex;-ms-flex-direction:column;flex-direction:column;-ms-flex-align:center;align-items:center;padding:0 16px}.forgot-password-success__text{width:100%;max-width:367px;font-size:24px;line-height:32px;text-align:center;color:var(--text)}.forgot-password-success__image{margin-bottom:56px}.reset-password{width:var(--content-width)}.reset-password__title{text-align:center;font-size:24px;line-height:32px;font-weight:400;color:var(--text);max-width:270px;margin:0 auto 16px}.reset-password__subtitle{text-align:center;font-size:16px;line-height:24px;color:#828282;max-width:189px;text-align:center;margin:0 auto 48px}.reset-password .error{margin-top:12px}.reset-password .primary-button{margin:48px auto 0}.reset-password .primary-button-mt-32{margin-top:32px}@media (max-width: 599px){.reset-password{width:100%;max-width:var(--content-width);padding:0 24px}.reset-password .primary-button{margin-top:32px}}";

const emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
const IdentifoForm$1 = class extends HTMLElement {
  constructor() {
    super();
    this.__registerHost();
    attachShadow(this);
    this.complete = createEvent(this, "complete", 7);
    this.error = createEvent(this, "error", 7);
    this.afterLoginRedirect = (e) => {
      this.phone = e.user.phone || '';
      this.email = e.user.email || '';
      this.lastResponse = e;
      if (e.require_2fa) {
        if (!e.enabled_2fa) {
          return 'tfa/setup';
        }
        if (e.enabled_2fa) {
          return 'tfa/verify';
        }
      }
      if (e.access_token && e.refresh_token) {
        return 'callback';
      }
      if (e.access_token && !e.refresh_token) {
        return 'callback';
      }
    };
    this.loginCatchRedirect = (data) => {
      if (data.id === APIErrorCodes.PleaseEnableTFA) {
        return 'tfa/setup';
      }
      throw data;
    };
  }
  // /**
  //  * The last name
  //  */
  // @Prop() last: string;
  // private getText(): string {
  //   return format(this.first, this.middle, this.last);
  // }
  processError(e) {
    this.lastError = e;
    this.error.emit(e);
  }
  async signIn() {
    await this.auth.api
      .login(this.email, this.password, '', this.scopes.split(','))
      .then(this.afterLoginRedirect)
      .catch(this.loginCatchRedirect)
      .then(route => this.openRoute(route))
      .catch(e => this.processError(e));
  }
  async signUp() {
    if (!this.validateEmail(this.email)) {
      return;
    }
    await this.auth.api
      .register(this.email, this.password, this.scopes.split(','))
      .then(this.afterLoginRedirect)
      .catch(this.loginCatchRedirect)
      .then(route => this.openRoute(route))
      .catch(e => this.processError(e));
  }
  async verifyTFA() {
    this.auth.api
      .verifyTFA(this.tfaCode, [])
      .then(() => this.openRoute('callback'))
      .catch(e => this.processError(e));
  }
  async setupTFA() {
    if (this.tfaType == TFAType.TFATypeSMS) {
      try {
        await this.auth.api.updateUser({ new_phone: this.phone });
      }
      catch (e) {
        this.processError(e);
        return;
      }
    }
    await this.auth.api.enableTFA().then(r => {
      if (!r.provisioning_uri) {
        this.openRoute('tfa/verify');
      }
      if (r.provisioning_uri) {
        this.provisioningURI = r.provisioning_uri;
        this.provisioningQR = r.provisioning_qr;
        this.openRoute('tfa/verify');
      }
    });
  }
  restorePassword() {
    this.auth.api
      .requestResetPassword(this.email)
      .then(() => {
      this.success = true;
      this.openRoute('password/forgot/success');
    })
      .catch(e => this.processError(e));
  }
  setNewPassword() {
    if (this.token) {
      this.auth.tokenService.saveToken(this.token, 'access');
    }
    this.auth.api
      .resetPassword(this.password)
      .then(() => {
      this.success = true;
      this.openRoute('login');
      this.password = '';
    })
      .catch(e => this.processError(e));
  }
  openRoute(route) {
    this.lastError = undefined;
    this.route = route;
  }
  usernameChange(event) {
    this.username = event.target.value;
  }
  passwordChange(event) {
    this.password = event.target.value;
  }
  emailChange(event) {
    this.email = event.target.value;
  }
  phoneChange(event) {
    this.phone = event.target.value;
  }
  tfaCodeChange(event) {
    this.tfaCode = event.target.value;
  }
  validateEmail(email) {
    if (!emailRegex.test(email)) {
      this.processError({ detailedMessage: 'Email address is not valid', name: 'Validation error', message: 'Email address is not valid' });
      return false;
    }
    return true;
  }
  renderRoute(route) {
    var _a, _b, _c, _d, _e, _f;
    switch (route) {
      case 'login':
        return (h("div", { class: "login-form" }, h("p", { class: "login-form__register-text" }, "Don't have an account?", h("a", { onClick: () => this.openRoute('register'), class: "login-form__register-link" }, ' ', "Sign Up")), h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingInput", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn() }), h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingPassword", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn() }), !!this.lastError && (h("div", { class: "error", role: "alert" }, (_a = this.lastError) === null || _a === void 0 ? void 0 : _a.detailedMessage)), h("div", { class: `login-form__buttons ${!!this.lastError ? 'login-form__buttons_mt-32' : ''}` }, h("button", { onClick: () => this.signIn(), class: "primary-button", disabled: !this.email || !this.password }, "Login"), h("a", { onClick: () => this.openRoute('password/forgot'), class: "login-form__forgot-pass" }, "Forgot password")), h("div", { class: "social-buttons" }, h("p", { class: "social-buttons__text" }, "or continue with"), h("div", { class: "social-buttons__social-medias" }, h("div", { class: "social-buttons__media" }, h("img", { src: getAssetPath(`assets/images/${'apple.svg'}`), class: "social-buttons__image", alt: "login via apple" })), h("div", { class: "social-buttons__media" }, h("img", { src: getAssetPath(`assets/images/${'google.svg'}`), class: "social-buttons__image", alt: "login via google" })), h("div", { class: "social-buttons__media" }, h("img", { src: getAssetPath(`assets/images/${'fb.svg'}`), class: "social-buttons__image", alt: "login via facebook" }))))));
      case 'register':
        return (h("div", { class: "register-form" }, h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingInput", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp() }), h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingPassword", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp() }), !!this.lastError && (h("div", { class: "error", role: "alert" }, (_b = this.lastError) === null || _b === void 0 ? void 0 : _b.detailedMessage)), h("div", { class: `register-form__buttons ${!!this.lastError ? 'register-form__buttons_mt-32' : ''}` }, h("button", { onClick: () => this.signUp(), class: "primary-button", disabled: !this.email || !this.password }, "Continue"), h("a", { onClick: () => this.openRoute('login'), class: "register-form__log-in" }, "Go back to login"))));
      case 'otp/login':
        return (h("div", { class: "otp-login" }, h("p", { class: "otp-login__register-text" }, "Don't have an account?", h("a", { onClick: () => this.openRoute('register'), class: "login-form__register-link" }, ' ', "Sign Up")), h("input", { type: "phone", class: "form-control", id: "floatingInput", value: this.phone, placeholder: "Phone number", onInput: event => this.phoneChange(event) }), h("button", { onClick: () => this.openRoute('tfa/verify'), class: "primary-button", disabled: !this.phone }, "Continue"), h("div", { class: "social-buttons" }, h("p", { class: "social-buttons__text" }, "or continue with"), h("div", { class: "social-buttons__social-medias" }, h("div", { class: "social-buttons__media" }, h("img", { src: getAssetPath(`./assets/images/${'apple.svg'}`), class: "social-buttons__image", alt: "login via apple" })), h("div", { class: "social-buttons__media" }, h("img", { src: getAssetPath(`./assets/images/${'google.svg'}`), class: "social-buttons__image", alt: "login via google" })), h("div", { class: "social-buttons__media" }, h("img", { src: getAssetPath(`./assets/images/${'fb.svg'}`), class: "social-buttons__image", alt: "login via facebook" }))))));
      case 'tfa/setup':
        return (h("div", { class: "tfa-setup" }, h("p", { class: "tfa-setup__text" }, "Protect your account with 2-step verification"), this.tfaType === TFAType.TFATypeApp && (h("div", { class: "info-card" }, h("div", { class: "info-card__controls" }, h("p", { class: "info-card__title" }, "Authenticator app"), h("button", { type: "button", class: "info-card__button", onClick: () => this.setupTFA() }, "Setup")), h("p", { class: "info-card__text" }, "Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone."))), this.tfaType === TFAType.TFATypeEmail && (h("div", { class: "info-card" }, h("div", { class: "info-card__controls" }, h("p", { class: "info-card__title" }, "Email"), h("button", { type: "button", class: "info-card__button", onClick: () => this.setupTFA() }, "Setup")), h("p", { class: "info-card__subtitle" }, this.email), h("p", { class: "info-card__text" }, " Use email as 2fa, please check your email, we will send confirmation code to this email."))), this.tfaType === TFAType.TFATypeSMS && (h("div", { class: "tfa-setup__form" }, h("p", { class: "tfa-setup__subtitle" }, " Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone"), h("input", { type: "phone", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingInput", value: this.phone, placeholder: "Phone", onInput: event => this.phoneChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.phone) && this.setupTFA() }), !!this.lastError && (h("div", { class: "error", role: "alert" }, (_c = this.lastError) === null || _c === void 0 ? void 0 : _c.detailedMessage)), h("button", { onClick: () => this.setupTFA(), class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.phone }, "Setup phone")))));
      case 'tfa/verify':
        return (h("div", { class: "tfa-verify" }, !!(this.tfaType === TFAType.TFATypeApp) && (h("div", { class: "tfa-verify__title-wrapper" }, h("h2", { class: this.provisioningURI ? 'tfa-verify__title' : 'tfa-verify__title_mb-40' }, !!this.provisioningURI ? 'Please scan QR-code with the app' : 'Use GoogleAuth as 2fa'), !!this.provisioningURI && h("img", { src: `data:image/png;base64, ${this.provisioningQR}`, alt: this.provisioningURI, class: "tfa-verify__qr-code" }))), !!(this.tfaType === TFAType.TFATypeSMS) && (h("div", { class: "tfa-verify__title-wrapper" }, h("h2", { class: "tfa-verify__title" }, "Enter the code sent to your phone number"), h("p", { class: "tfa-verify__subtitle" }, "The code has been sent to ", this.phone))), !!(this.tfaType === TFAType.TFATypeEmail) && (h("div", { class: "tfa-verify__title-wrapper" }, h("h2", { class: "tfa-verify__title" }, "Enter the code sent to your email address"), h("p", { class: "tfa-verify__subtitle" }, "The email has been sent to ", this.email))), h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingCode", value: this.tfaCode, placeholder: "Verify code", onInput: event => this.tfaCodeChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.tfaCode) && this.verifyTFA() }), !!this.lastError && (h("div", { class: "error", role: "alert" }, (_d = this.lastError) === null || _d === void 0 ? void 0 : _d.detailedMessage)), h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.tfaCode, onClick: () => this.verifyTFA() }, "Confirm")));
      case 'password/forgot':
        return (h("div", { class: "forgot-password" }, h("h2", { class: "forgot-password__title" }, "Enter the email you gave when you registered"), h("p", { class: "forgot-password__subtitle" }, "We will send you a link to create a new password on email"), h("input", { type: "email", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingEmail", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email) && this.restorePassword() }), !!this.lastError && (h("div", { class: "error", role: "alert" }, (_e = this.lastError) === null || _e === void 0 ? void 0 : _e.detailedMessage)), h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.email, onClick: () => this.restorePassword() }, "Send the link")));
      case 'password/forgot/success':
        return (h("div", { class: "forgot-password-success" }, this.theme === 'dark' && h("img", { src: getAssetPath(`./assets/images/${'email-dark.svg'}`), alt: "email", class: "forgot-password-success__image" }), this.theme === 'light' && h("img", { src: getAssetPath(`./assets/images/${'email.svg'}`), alt: "email", class: "forgot-password-success__image" }), h("p", { class: "forgot-password-success__text" }, "We sent you an email with a link to create a new password")));
      case 'password/reset':
        return (h("div", { class: "reset-password" }, h("h2", { class: "reset-password__title" }, "Set up a new password to log in to the website"), h("p", { class: "reset-password__subtitle" }, "Memorize your password and do not give it to anyone."), h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingPassword", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password) && this.setNewPassword() }), !!this.lastError && (h("div", { class: "error", role: "alert" }, (_f = this.lastError) === null || _f === void 0 ? void 0 : _f.detailedMessage)), h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.password, onClick: () => this.setNewPassword() }, "Save password")));
      case 'error':
        return (h("div", { class: "error-view" }, h("div", { class: "error-view__message" }, this.lastError.message), h("div", { class: "error-view__details" }, this.lastError.detailedMessage)));
    }
  }
  async componentWillLoad() {
    this.auth = new IdentifoAuth({ appId: this.appId, url: this.url });
    try {
      const settings = await this.auth.api.getAppSettings();
      this.registrationForbidden = settings.registrationForbidden;
      this.tfaType = settings.tfaType;
    }
    catch (err) {
      this.route = 'error';
      this.lastError = err;
    }
  }
  componentWillRender() {
    if (this.route === 'callback') {
      this.complete.emit(this.lastResponse);
    }
    if (this.route === 'logout') {
      this.auth.api.logout().then(() => this.complete.emit());
    }
  }
  render() {
    return h("div", { class: { 'wrapper': this.theme === 'light', 'wrapper-dark': this.theme === 'dark' } }, this.renderRoute(this.route));
  }
  static get assetsDirs() { return ["assets"]; }
  static get style() { return mainCss; }
};

const IdentifoForm = /*@__PURE__*/proxyCustomElement(IdentifoForm$1, [1,"identifo-form",{"route":[1537],"token":[1],"appId":[513,"app-id"],"url":[513],"theme":[1],"scopes":[1],"auth":[32],"username":[32],"password":[32],"phone":[32],"email":[32],"registrationForbidden":[32],"tfaCode":[32],"tfaType":[32],"tfaMandatory":[32],"provisioningURI":[32],"provisioningQR":[32],"success":[32],"lastError":[32],"lastResponse":[32]}]);
const defineCustomElements = (opts) => {
  if (typeof customElements !== 'undefined') {
    [
      IdentifoForm
    ].forEach(cmp => {
      if (!customElements.get(cmp.is)) {
        customElements.define(cmp.is, cmp, opts);
      }
    });
  }
};

export { IdentifoForm, defineCustomElements };

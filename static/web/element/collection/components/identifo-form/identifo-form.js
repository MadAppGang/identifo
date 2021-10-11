import { APIErrorCodes, IdentifoAuth, TFAType, TFAStatus } from '@identifo/identifo-auth-js';
import { Component, Event, getAssetPath, h, Host, Prop, State } from '@stencil/core';
const emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
export class IdentifoForm {
  constructor() {
    this.route = 'login';
    this.theme = 'auto';
    this.scopes = '';
    this.selectedTheme = 'light';
    this.federatedProviders = [];
    this.afterLoginRedirect = (e) => {
      this.phone = e.user.phone || '';
      this.email = e.user.email || '';
      this.lastResponse = e;
      if (e.require_2fa) {
        if (!e.enabled_2fa) {
          return this.redirectTfa('tfa/setup');
        }
        if (e.enabled_2fa) {
          return this.redirectTfa('tfa/verify');
        }
      }
      if (this.tfaStatus === TFAStatus.OPTIONAL) {
        return `tfa/setup/select`;
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
        return this.redirectTfa('tfa/setup');
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
    var _a, _b;
    e.detailedMessage = (_a = e.detailedMessage) === null || _a === void 0 ? void 0 : _a.trim();
    e.message = (_b = e.message) === null || _b === void 0 ? void 0 : _b.trim();
    this.lastError = e;
    this.error.emit(e);
  }
  redirectTfa(prefix) {
    if (this.tfaTypes.length === 1) {
      return `${prefix}/${this.tfaTypes[0]}`;
    }
    else {
      return `${prefix}/select`;
    }
  }
  async signIn() {
    if (!this.validateEmail(this.email)) {
      return;
    }
    await this.auth.api
      .login(this.email, this.password, '', this.scopes.split(','))
      .then(this.afterLoginRedirect)
      .catch(this.loginCatchRedirect)
      .then(route => this.openRoute(route))
      .catch(e => this.processError(e));
  }
  async loginWith(provider) {
    this.route = 'loading';
    const federatedRedirectUrl = this.federatedRedirectUrl || window.location.origin + window.location.pathname;
    this.auth.api.federatedLogin(provider, this.scopes.split(','), federatedRedirectUrl, this.callbackUrl);
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
    if (this.route.indexOf('password/forgot/tfa') === 0) {
      this.auth.api
        .requestResetPassword(this.email, this.tfaCode)
        .then(() => {
        this.success = true;
        this.openRoute('password/forgot/success');
      })
        .catch(e => this.processError(e));
    }
    else {
      this.auth.api
        .verifyTFA(this.tfaCode, [])
        .then(() => this.openRoute('callback'))
        .catch(e => this.processError(e));
    }
  }
  async selectTFA(type) {
    this.openRoute(`tfa/setup/${type}`);
  }
  async setupTFA(type) {
    switch (type) {
      case TFAType.TFATypeApp:
        break;
      case TFAType.TFATypeEmail:
        await this.auth.api.enableTFA();
        break;
      case TFAType.TFATypeSMS:
        try {
          await this.auth.api.updateUser({ new_phone: this.phone });
        }
        catch (e) {
          this.processError(e);
          return;
        }
        await this.auth.api.enableTFA();
        break;
    }
    this.openRoute(`tfa/verify/${type}`);
  }
  restorePassword() {
    this.auth.api
      .requestResetPassword(this.email)
      .then(response => {
      if (response.result === 'tfa-required') {
        this.openRoute(this.redirectTfa('password/forgot/tfa'));
      }
      if (response.result === 'ok') {
        this.success = true;
        this.openRoute('password/forgot/success');
      }
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
      this.processError({ detailedMessage: 'Email address is not valid.', name: 'Validation error', message: 'Email address is not valid.' });
      return false;
    }
    return true;
  }
  renderBackToLogin() {
    return (h("a", { onClick: () => this.openRoute('login'), class: "forgot-password__login" }, "Go back to login"));
  }
  renderRoute(route) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p;
    switch (route) {
      case 'login':
        return (h("div", { class: "login-form" },
          !this.registrationForbidden && (h("p", { class: "login-form__register-text" },
            "Don't have an account?\u00A0",
            h("a", { onClick: () => this.openRoute('register'), class: "login-form__register-link" }, "Sign Up"))),
          h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "login", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn() }),
          h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "password", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, ((_a = this.lastError) === null || _a === void 0 ? void 0 : _a.detailedMessage) || ((_b = this.lastError) === null || _b === void 0 ? void 0 : _b.message))),
          h("div", { class: `login-form__buttons ${!!this.lastError ? 'login-form__buttons_mt-32' : ''}` },
            h("button", { onClick: () => this.signIn(), class: "primary-button", disabled: !this.email || !this.password }, "Login"),
            h("a", { onClick: () => this.openRoute('password/forgot'), class: "login-form__forgot-pass" }, "Forgot password")),
          this.federatedProviders.length > 0 && (h("div", { class: "social-buttons" },
            h("p", { class: "social-buttons__text" }, "or continue with"),
            h("div", { class: "social-buttons__social-medias" },
              this.federatedProviders.indexOf('apple') > -1 && (h("div", { class: "social-buttons__media social-buttons__apple", onClick: () => this.loginWith('apple') },
                h("img", { src: getAssetPath(`assets/images/${'apple.svg'}`), class: "social-buttons__image", alt: "login via apple" }))),
              this.federatedProviders.indexOf('google') > -1 && (h("div", { class: "social-buttons__media social-buttons__google", onClick: () => this.loginWith('google') },
                h("img", { src: getAssetPath(`assets/images/${'google.svg'}`), class: "social-buttons__image", alt: "login via google" }))),
              this.federatedProviders.indexOf('facebook') > -1 && (h("div", { class: "social-buttons__media social-buttons__facebook", onClick: () => this.loginWith('facebook') },
                h("img", { src: getAssetPath(`assets/images/${'fb.svg'}`), class: "social-buttons__image", alt: "login via facebook" }))))))));
      case 'register':
        return (h("div", { class: "register-form" },
          h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "login", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp() }),
          h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "password", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, ((_c = this.lastError) === null || _c === void 0 ? void 0 : _c.detailedMessage) || ((_d = this.lastError) === null || _d === void 0 ? void 0 : _d.message))),
          h("div", { class: `register-form__buttons ${!!this.lastError ? 'register-form__buttons_mt-32' : ''}` },
            h("button", { onClick: () => this.signUp(), class: "primary-button", disabled: !this.email || !this.password }, "Continue"),
            this.renderBackToLogin())));
      case 'otp/login':
        return (h("div", { class: "otp-login" },
          !this.registrationForbidden && (h("p", { class: "otp-login__register-text" },
            "Don't have an account?\u00A0",
            h("a", { onClick: () => this.openRoute('register'), class: "login-form__register-link" }, "Sign Up"))),
          h("input", { type: "phone", class: "form-control", id: "login", value: this.phone, placeholder: "Phone number", onInput: event => this.phoneChange(event) }),
          h("button", { onClick: () => this.openRoute(this.redirectTfa('tfa/verify')), class: "primary-button", disabled: !this.phone }, "Continue"),
          this.federatedProviders.length > 0 && (h("div", { class: "social-buttons" },
            h("p", { class: "social-buttons__text" }, "or continue with"),
            h("div", { class: "social-buttons__social-medias" },
              this.federatedProviders.indexOf('apple') > -1 && (h("div", { class: "social-buttons__media social-buttons__apple", onClick: () => this.loginWith('apple') },
                h("img", { src: getAssetPath(`assets/images/${'apple.svg'}`), class: "social-buttons__image", alt: "login via apple" }))),
              this.federatedProviders.indexOf('google') > -1 && (h("div", { class: "social-buttons__media social-buttons__google", onClick: () => this.loginWith('google') },
                h("img", { src: getAssetPath(`assets/images/${'google.svg'}`), class: "social-buttons__image", alt: "login via google" }))),
              this.federatedProviders.indexOf('facebook') > -1 && (h("div", { class: "social-buttons__media social-buttons__facebook", onClick: () => this.loginWith('facebook') },
                h("img", { src: getAssetPath(`assets/images/${'fb.svg'}`), class: "social-buttons__image", alt: "login via facebook" }))))))));
      case 'tfa/verify/select':
      case 'tfa/setup/select':
      case 'password/forgot/tfa/select':
        return (h("div", { class: "tfa-setup" },
          this.route === 'tfa/verify/select' && h("p", { class: "tfa-setup__text" }, "Select 2-step verification method"),
          this.route === 'tfa/setup/select' && h("p", { class: "tfa-setup__text" }, "Protect your account with 2-step verification"),
          this.tfaTypes.includes(TFAType.TFATypeApp) && (h("div", { class: "info-card info-card-app" },
            h("div", { class: "info-card__controls" },
              h("p", { class: "info-card__title" }, "Authenticator app"),
              h("button", { type: "button", class: "info-card__button", onClick: () => this.selectTFA(TFAType.TFATypeApp) }, "Setup")),
            h("p", { class: "info-card__text" }, "Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone."))),
          this.tfaTypes.includes(TFAType.TFATypeEmail) && (h("div", { class: "info-card info-card-email" },
            h("div", { class: "info-card__controls" },
              h("p", { class: "info-card__title" }, "Email"),
              h("button", { type: "button", class: "info-card__button", onClick: () => this.selectTFA(TFAType.TFATypeEmail) }, "Setup")),
            h("p", { class: "info-card__subtitle" }, this.email),
            h("p", { class: "info-card__text" }, " Use email as 2fa, please check your email, we will send confirmation code to this email."))),
          this.tfaTypes.includes(TFAType.TFATypeSMS) && (h("div", { class: "info-card info-card-sms" },
            h("div", { class: "info-card__controls" },
              h("p", { class: "info-card__title" }, "SMS"),
              h("button", { type: "button", class: "info-card__button", onClick: () => this.selectTFA(TFAType.TFATypeSMS) }, "Setup")),
            h("p", { class: "info-card__subtitle" }, this.phone),
            h("p", { class: "info-card__text" }, " Use phone as 2fa, please check your phone, we will send confirmation code to this phone"))),
          this.route === 'tfa/setup/select' && this.tfaStatus === TFAStatus.OPTIONAL && (h("a", { onClick: () => this.openRoute('callback'), class: "forgot-password__login" }, "Setup next time")),
          this.tfaStatus !== TFAStatus.OPTIONAL && this.renderBackToLogin()));
      case 'tfa/setup/email':
      case 'tfa/setup/sms':
      case 'tfa/setup/app':
        return (h("div", { class: "tfa-setup" },
          h("p", { class: "tfa-setup__text" }, "Protect your account with 2-step verification"),
          this.route === 'tfa/setup/app' && (h("div", { class: "tfa-setup__form" },
            h("p", { class: "tfa-setup__subtitle" }, "Please scan QR-code with the app and click Continue"),
            h("div", { class: "tfa-setup__qr-wrapper" }, !!this.provisioningURI && h("img", { src: `data:image/png;base64, ${this.provisioningQR}`, alt: this.provisioningURI, class: "tfa-setup__qr-code" })),
            h("button", { onClick: () => this.setupTFA(TFAType.TFATypeApp), class: `primary-button ${this.lastError && 'primary-button-mt-32'}` }, "Continue"))),
          this.route === 'tfa/setup/email' && (h("div", { class: "tfa-setup__form" },
            h("p", { class: "tfa-setup__subtitle" }, " Use email as 2fa, please check your enail bellow, we will send confirmation code to this email"),
            h("input", { type: "email", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "email", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email) && this.setupTFA(TFAType.TFATypeEmail) }),
            !!this.lastError && (h("div", { class: "error", role: "alert" }, ((_e = this.lastError) === null || _e === void 0 ? void 0 : _e.detailedMessage) || ((_f = this.lastError) === null || _f === void 0 ? void 0 : _f.message))),
            h("button", { onClick: () => this.setupTFA(TFAType.TFATypeEmail), class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.email }, "Setup email"))),
          this.route === 'tfa/setup/sms' && (h("div", { class: "tfa-setup__form" },
            h("p", { class: "tfa-setup__subtitle" }, " Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone"),
            h("input", { type: "phone", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "phone", value: this.phone, placeholder: "Phone", onInput: event => this.phoneChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.phone) && this.setupTFA(TFAType.TFATypeSMS) }),
            !!this.lastError && (h("div", { class: "error", role: "alert" }, ((_g = this.lastError) === null || _g === void 0 ? void 0 : _g.detailedMessage) || ((_h = this.lastError) === null || _h === void 0 ? void 0 : _h.message))),
            h("button", { onClick: () => this.setupTFA(TFAType.TFATypeSMS), class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.phone }, "Setup phone"))),
          this.renderBackToLogin()));
      case 'tfa/verify/app':
      case 'tfa/verify/email':
      case 'tfa/verify/sms':
      case 'password/forgot/tfa/app':
      case 'password/forgot/tfa/email':
      case 'password/forgot/tfa/sms':
        return (h("div", { class: "tfa-verify" },
          this.route.indexOf('app') > 0 && (h("div", { class: "tfa-verify__title-wrapper" },
            h("h2", { class: "tfa-verify__title" }, "Enter the code from authenticator app"),
            h("p", { class: "tfa-verify__subtitle" }, "Code will be generated by app"))),
          this.route.indexOf('sms') > 0 && (h("div", { class: "tfa-verify__title-wrapper" },
            h("h2", { class: "tfa-verify__title" }, "Enter the code sent to your phone number"),
            h("p", { class: "tfa-verify__subtitle" },
              "The code has been sent to ",
              this.phone))),
          this.route.indexOf('email') > 0 && (h("div", { class: "tfa-verify__title-wrapper" },
            h("h2", { class: "tfa-verify__title" }, "Enter the code sent to your email address"),
            h("p", { class: "tfa-verify__subtitle" },
              "The email has been sent to ",
              this.email))),
          h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "tfaCode", value: this.tfaCode, placeholder: "Verify code", onInput: event => this.tfaCodeChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.tfaCode) && this.verifyTFA() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, ((_j = this.lastError) === null || _j === void 0 ? void 0 : _j.detailedMessage) || ((_k = this.lastError) === null || _k === void 0 ? void 0 : _k.message))),
          h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.tfaCode, onClick: () => this.verifyTFA() }, "Confirm"),
          this.renderBackToLogin()));
      case 'password/forgot':
        return (h("div", { class: "forgot-password" },
          h("h2", { class: "forgot-password__title" }, "Enter the email you gave when you registered"),
          h("p", { class: "forgot-password__subtitle" }, "We will send you a link to create a new password on email"),
          h("input", { type: "email", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "email", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email) && this.restorePassword() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, ((_l = this.lastError) === null || _l === void 0 ? void 0 : _l.detailedMessage) || ((_m = this.lastError) === null || _m === void 0 ? void 0 : _m.message))),
          h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.email, onClick: () => this.restorePassword() }, "Send the link"),
          this.renderBackToLogin()));
      case 'password/forgot/success':
        return (h("div", { class: "forgot-password-success" },
          this.selectedTheme === 'dark' && h("img", { src: getAssetPath(`./assets/images/${'email-dark.svg'}`), alt: "email", class: "forgot-password-success__image" }),
          this.selectedTheme === 'light' && h("img", { src: getAssetPath(`./assets/images/${'email.svg'}`), alt: "email", class: "forgot-password-success__image" }),
          h("p", { class: "forgot-password-success__text" }, "We sent you an email with a link to create a new password"),
          this.renderBackToLogin()));
      case 'password/reset':
        return (h("div", { class: "reset-password" },
          h("h2", { class: "reset-password__title" }, "Set up a new password to log in to the website"),
          h("p", { class: "reset-password__subtitle" }, "Memorize your password and do not give it to anyone."),
          h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "password", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password) && this.setNewPassword() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, ((_o = this.lastError) === null || _o === void 0 ? void 0 : _o.detailedMessage) || ((_p = this.lastError) === null || _p === void 0 ? void 0 : _p.message))),
          h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.password, onClick: () => this.setNewPassword() }, "Save password")));
      case 'error':
        return (h("div", { class: "error-view" },
          h("div", { class: "error-view__message" }, this.lastError.message),
          h("div", { class: "error-view__details" }, this.lastError.detailedMessage)));
      case 'callback':
        return (h("div", { class: "error-view" },
          h("div", null, "Success"),
          this.debug && (h("div", null,
            h("div", null,
              "Access token: ",
              this.lastResponse.access_token),
            h("div", null,
              "Refresh token: ",
              this.lastResponse.refresh_token),
            h("div", null,
              "User: ",
              JSON.stringify(this.lastResponse.user))))));
      case 'loading':
        return (h("div", { class: "error-view" },
          h("div", null, "Loading ...")));
    }
  }
  async componentWillLoad() {
    // const base = (document.querySelector('base') || {}).href;
    // const path = window.location.href.split('?')[0];
    // this.route = path.replace(base, '').replace(/^\/|\/$/g, '') as Routes;
    this.token = new URLSearchParams(window.location.search).get('token');
    const postLogoutRedirectUri = this.postLogoutRedirectUri || window.location.origin + window.location.pathname;
    if (!this.appId) {
      this.lastError = { message: 'app-id param is empty', name: 'app-id empty' };
      this.error.emit(this.lastError);
      this.route = 'error';
      return;
    }
    if (!this.url) {
      this.lastError = { message: 'url param is empty', name: 'url empty' };
      this.error.emit(this.lastError);
      this.route = 'error';
      return;
    }
    try {
      this.auth = new IdentifoAuth({ appId: this.appId, url: this.url, postLogoutRedirectUri });
      const settings = await this.auth.api.getAppSettings();
      this.registrationForbidden = settings.registrationForbidden;
      this.tfaTypes = Array.isArray(settings.tfaType) ? settings.tfaType : [settings.tfaType];
      this.tfaStatus = settings.tfaStatus;
      this.federatedProviders = settings.federatedProviders;
    }
    catch (err) {
      this.route = 'error';
      this.lastError = err;
    }
    // If we have provider and state then we need to complete federated login
    const href = new URL(window.location.href);
    if (!!href.searchParams.get('provider') && !!href.searchParams.get('state')) {
      // Also we clear all url params after parsing
      const u = new URL(window.location.href);
      const sp = new URLSearchParams();
      const appId = href.searchParams.get('appId');
      sp.set('appId', appId);
      window.history.replaceState({}, document.title, `${u.pathname}?${sp.toString()}`);
      this.route = 'loading';
      this.auth.api
        .federatedLoginComplete(u.searchParams)
        .then(this.afterLoginRedirect)
        .catch(this.loginCatchRedirect)
        .then(route => this.openRoute(route))
        .catch(e => this.processError(e));
    }
    // Auto theme select
    this.selectedTheme = 'light';
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      if (this.theme === 'auto') {
        this.selectedTheme = 'dark';
      }
    }
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
      if (this.theme === 'auto') {
        this.selectedTheme = e.matches ? 'dark' : 'light';
      }
    });
  }
  componentWillRender() {
    if (this.route === 'callback') {
      const u = new URL(window.location.href);
      u.searchParams.set('callbackUrl', this.lastResponse.callbackUrl);
      window.history.replaceState({}, document.title, `${u.pathname}?${u.searchParams.toString()}`);
      this.complete.emit(this.lastResponse);
    }
    if (this.route === 'logout') {
      this.complete.emit();
    }
    if (this.route === 'tfa/setup/app') {
      this.auth.api.enableTFA().then(r => {
        if (r.provisioning_uri) {
          this.provisioningURI = r.provisioning_uri;
          this.provisioningQR = r.provisioning_qr;
        }
      });
    }
  }
  render() {
    return (h(Host, null,
      h("div", { class: { 'wrapper': this.selectedTheme === 'light', 'wrapper-dark': this.selectedTheme === 'dark' } }, this.renderRoute(this.route)),
      h("div", { class: "error-view" }, this.debug && (h("div", null,
        h("br", null),
        this.appId)))));
  }
  static get is() { return "identifo-form"; }
  static get originalStyleUrls() { return {
    "$": ["../../styles/identifo-form/main.scss"]
  }; }
  static get styleUrls() { return {
    "$": ["../../styles/identifo-form/main.css"]
  }; }
  static get assetsDirs() { return ["assets"]; }
  static get properties() { return {
    "route": {
      "type": "string",
      "mutable": true,
      "complexType": {
        "original": "Routes",
        "resolved": "\"login\" | \"register\" | \"tfa/verify/sms\" | \"tfa/verify/email\" | \"tfa/verify/app\" | \"tfa/verify/select\" | \"tfa/setup/sms\" | \"tfa/setup/email\" | \"tfa/setup/app\" | \"tfa/setup/select\" | \"password/reset\" | \"password/forgot\" | \"password/forgot/tfa/sms\" | \"password/forgot/tfa/email\" | \"password/forgot/tfa/app\" | \"password/forgot/tfa/select\" | \"callback\" | \"otp/login\" | \"error\" | \"password/forgot/success\" | \"logout\" | \"loading\"",
        "references": {
          "Routes": {
            "location": "local"
          }
        }
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "route",
      "reflect": true,
      "defaultValue": "'login'"
    },
    "token": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "string",
        "resolved": "string",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "token",
      "reflect": false
    },
    "appId": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "string",
        "resolved": "string",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "app-id",
      "reflect": true
    },
    "url": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "string",
        "resolved": "string",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "url",
      "reflect": true
    },
    "theme": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "'dark' | 'light' | 'auto'",
        "resolved": "\"auto\" | \"dark\" | \"light\"",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "theme",
      "reflect": false,
      "defaultValue": "'auto'"
    },
    "scopes": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "string",
        "resolved": "string",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "scopes",
      "reflect": false,
      "defaultValue": "''"
    },
    "callbackUrl": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "string",
        "resolved": "string",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "callback-url",
      "reflect": false
    },
    "federatedRedirectUrl": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "string",
        "resolved": "string",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "federated-redirect-url",
      "reflect": false
    },
    "postLogoutRedirectUri": {
      "type": "string",
      "mutable": false,
      "complexType": {
        "original": "string",
        "resolved": "string",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "post-logout-redirect-uri",
      "reflect": false
    },
    "debug": {
      "type": "boolean",
      "mutable": false,
      "complexType": {
        "original": "boolean",
        "resolved": "boolean",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "debug",
      "reflect": false
    }
  }; }
  static get states() { return {
    "selectedTheme": {},
    "auth": {},
    "username": {},
    "password": {},
    "phone": {},
    "email": {},
    "registrationForbidden": {},
    "tfaCode": {},
    "tfaTypes": {},
    "federatedProviders": {},
    "tfaStatus": {},
    "provisioningURI": {},
    "provisioningQR": {},
    "success": {},
    "lastError": {},
    "lastResponse": {}
  }; }
  static get events() { return [{
      "method": "complete",
      "name": "complete",
      "bubbles": true,
      "cancelable": true,
      "composed": true,
      "docs": {
        "tags": [],
        "text": ""
      },
      "complexType": {
        "original": "LoginResponse",
        "resolved": "LoginResponse",
        "references": {
          "LoginResponse": {
            "location": "import",
            "path": "@identifo/identifo-auth-js"
          }
        }
      }
    }, {
      "method": "error",
      "name": "error",
      "bubbles": true,
      "cancelable": true,
      "composed": true,
      "docs": {
        "tags": [],
        "text": ""
      },
      "complexType": {
        "original": "ApiError",
        "resolved": "ApiError",
        "references": {
          "ApiError": {
            "location": "import",
            "path": "@identifo/identifo-auth-js"
          }
        }
      }
    }]; }
}

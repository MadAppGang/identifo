import { APIErrorCodes, IdentifoAuth, TFAType } from '@identifo/identifo-auth-js';
import { Component, Event, getAssetPath, h, Prop, State } from '@stencil/core';
const emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
export class IdentifoForm {
  constructor() {
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
        return (h("div", { class: "login-form" },
          h("p", { class: "login-form__register-text" },
            "Don't have an account?",
            h("a", { onClick: () => this.openRoute('register'), class: "login-form__register-link" },
              ' ',
              "Sign Up")),
          h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingInput", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn() }),
          h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingPassword", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, (_a = this.lastError) === null || _a === void 0 ? void 0 : _a.detailedMessage)),
          h("div", { class: `login-form__buttons ${!!this.lastError ? 'login-form__buttons_mt-32' : ''}` },
            h("button", { onClick: () => this.signIn(), class: "primary-button", disabled: !this.email || !this.password }, "Login"),
            h("a", { onClick: () => this.openRoute('password/forgot'), class: "login-form__forgot-pass" }, "Forgot password")),
          h("div", { class: "social-buttons" },
            h("p", { class: "social-buttons__text" }, "or continue with"),
            h("div", { class: "social-buttons__social-medias" },
              h("div", { class: "social-buttons__media" },
                h("img", { src: getAssetPath(`assets/images/${'apple.svg'}`), class: "social-buttons__image", alt: "login via apple" })),
              h("div", { class: "social-buttons__media" },
                h("img", { src: getAssetPath(`assets/images/${'google.svg'}`), class: "social-buttons__image", alt: "login via google" })),
              h("div", { class: "social-buttons__media" },
                h("img", { src: getAssetPath(`assets/images/${'fb.svg'}`), class: "social-buttons__image", alt: "login via facebook" }))))));
      case 'register':
        return (h("div", { class: "register-form" },
          h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingInput", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp() }),
          h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingPassword", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, (_b = this.lastError) === null || _b === void 0 ? void 0 : _b.detailedMessage)),
          h("div", { class: `register-form__buttons ${!!this.lastError ? 'register-form__buttons_mt-32' : ''}` },
            h("button", { onClick: () => this.signUp(), class: "primary-button", disabled: !this.email || !this.password }, "Continue"),
            h("a", { onClick: () => this.openRoute('login'), class: "register-form__log-in" }, "Go back to login"))));
      case 'otp/login':
        return (h("div", { class: "otp-login" },
          h("p", { class: "otp-login__register-text" },
            "Don't have an account?",
            h("a", { onClick: () => this.openRoute('register'), class: "login-form__register-link" },
              ' ',
              "Sign Up")),
          h("input", { type: "phone", class: "form-control", id: "floatingInput", value: this.phone, placeholder: "Phone number", onInput: event => this.phoneChange(event) }),
          h("button", { onClick: () => this.openRoute('tfa/verify'), class: "primary-button", disabled: !this.phone }, "Continue"),
          h("div", { class: "social-buttons" },
            h("p", { class: "social-buttons__text" }, "or continue with"),
            h("div", { class: "social-buttons__social-medias" },
              h("div", { class: "social-buttons__media" },
                h("img", { src: getAssetPath(`./assets/images/${'apple.svg'}`), class: "social-buttons__image", alt: "login via apple" })),
              h("div", { class: "social-buttons__media" },
                h("img", { src: getAssetPath(`./assets/images/${'google.svg'}`), class: "social-buttons__image", alt: "login via google" })),
              h("div", { class: "social-buttons__media" },
                h("img", { src: getAssetPath(`./assets/images/${'fb.svg'}`), class: "social-buttons__image", alt: "login via facebook" }))))));
      case 'tfa/setup':
        return (h("div", { class: "tfa-setup" },
          h("p", { class: "tfa-setup__text" }, "Protect your account with 2-step verification"),
          this.tfaType === TFAType.TFATypeApp && (h("div", { class: "info-card" },
            h("div", { class: "info-card__controls" },
              h("p", { class: "info-card__title" }, "Authenticator app"),
              h("button", { type: "button", class: "info-card__button", onClick: () => this.setupTFA() }, "Setup")),
            h("p", { class: "info-card__text" }, "Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone."))),
          this.tfaType === TFAType.TFATypeEmail && (h("div", { class: "info-card" },
            h("div", { class: "info-card__controls" },
              h("p", { class: "info-card__title" }, "Email"),
              h("button", { type: "button", class: "info-card__button", onClick: () => this.setupTFA() }, "Setup")),
            h("p", { class: "info-card__subtitle" }, this.email),
            h("p", { class: "info-card__text" }, " Use email as 2fa, please check your email, we will send confirmation code to this email."))),
          this.tfaType === TFAType.TFATypeSMS && (h("div", { class: "tfa-setup__form" },
            h("p", { class: "tfa-setup__subtitle" }, " Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone"),
            h("input", { type: "phone", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingInput", value: this.phone, placeholder: "Phone", onInput: event => this.phoneChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.phone) && this.setupTFA() }),
            !!this.lastError && (h("div", { class: "error", role: "alert" }, (_c = this.lastError) === null || _c === void 0 ? void 0 : _c.detailedMessage)),
            h("button", { onClick: () => this.setupTFA(), class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.phone }, "Setup phone")))));
      case 'tfa/verify':
        return (h("div", { class: "tfa-verify" },
          !!(this.tfaType === TFAType.TFATypeApp) && (h("div", { class: "tfa-verify__title-wrapper" },
            h("h2", { class: this.provisioningURI ? 'tfa-verify__title' : 'tfa-verify__title_mb-40' }, !!this.provisioningURI ? 'Please scan QR-code with the app' : 'Use GoogleAuth as 2fa'),
            !!this.provisioningURI && h("img", { src: `data:image/png;base64, ${this.provisioningQR}`, alt: this.provisioningURI, class: "tfa-verify__qr-code" }))),
          !!(this.tfaType === TFAType.TFATypeSMS) && (h("div", { class: "tfa-verify__title-wrapper" },
            h("h2", { class: "tfa-verify__title" }, "Enter the code sent to your phone number"),
            h("p", { class: "tfa-verify__subtitle" },
              "The code has been sent to ",
              this.phone))),
          !!(this.tfaType === TFAType.TFATypeEmail) && (h("div", { class: "tfa-verify__title-wrapper" },
            h("h2", { class: "tfa-verify__title" }, "Enter the code sent to your email address"),
            h("p", { class: "tfa-verify__subtitle" },
              "The email has been sent to ",
              this.email))),
          h("input", { type: "text", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingCode", value: this.tfaCode, placeholder: "Verify code", onInput: event => this.tfaCodeChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.tfaCode) && this.verifyTFA() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, (_d = this.lastError) === null || _d === void 0 ? void 0 : _d.detailedMessage)),
          h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.tfaCode, onClick: () => this.verifyTFA() }, "Confirm")));
      case 'password/forgot':
        return (h("div", { class: "forgot-password" },
          h("h2", { class: "forgot-password__title" }, "Enter the email you gave when you registered"),
          h("p", { class: "forgot-password__subtitle" }, "We will send you a link to create a new password on email"),
          h("input", { type: "email", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingEmail", value: this.email, placeholder: "Email", onInput: event => this.emailChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.email) && this.restorePassword() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, (_e = this.lastError) === null || _e === void 0 ? void 0 : _e.detailedMessage)),
          h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.email, onClick: () => this.restorePassword() }, "Send the link")));
      case 'password/forgot/success':
        return (h("div", { class: "forgot-password-success" },
          this.theme === 'dark' && h("img", { src: getAssetPath(`./assets/images/${'email-dark.svg'}`), alt: "email", class: "forgot-password-success__image" }),
          this.theme === 'light' && h("img", { src: getAssetPath(`./assets/images/${'email.svg'}`), alt: "email", class: "forgot-password-success__image" }),
          h("p", { class: "forgot-password-success__text" }, "We sent you an email with a link to create a new password")));
      case 'password/reset':
        return (h("div", { class: "reset-password" },
          h("h2", { class: "reset-password__title" }, "Set up a new password to log in to the website"),
          h("p", { class: "reset-password__subtitle" }, "Memorize your password and do not give it to anyone."),
          h("input", { type: "password", class: `form-control ${this.lastError && 'form-control-danger'}`, id: "floatingPassword", value: this.password, placeholder: "Password", onInput: event => this.passwordChange(event), onKeyPress: e => !!(e.key === 'Enter' && this.password) && this.setNewPassword() }),
          !!this.lastError && (h("div", { class: "error", role: "alert" }, (_f = this.lastError) === null || _f === void 0 ? void 0 : _f.detailedMessage)),
          h("button", { type: "button", class: `primary-button ${this.lastError && 'primary-button-mt-32'}`, disabled: !this.password, onClick: () => this.setNewPassword() }, "Save password")));
      case 'error':
        return (h("div", { class: "error-view" },
          h("div", { class: "error-view__message" }, this.lastError.message),
          h("div", { class: "error-view__details" }, this.lastError.detailedMessage)));
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
  static get is() { return "identifo-form"; }
  static get encapsulation() { return "shadow"; }
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
        "resolved": "\"callback\" | \"error\" | \"login\" | \"logout\" | \"otp/login\" | \"password/forgot\" | \"password/forgot/success\" | \"password/reset\" | \"register\" | \"tfa/setup\" | \"tfa/verify\"",
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
      "reflect": true
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
        "original": "'dark' | 'light'",
        "resolved": "\"dark\" | \"light\"",
        "references": {}
      },
      "required": false,
      "optional": false,
      "docs": {
        "tags": [],
        "text": ""
      },
      "attribute": "theme",
      "reflect": false
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
      "reflect": false
    }
  }; }
  static get states() { return {
    "auth": {},
    "username": {},
    "password": {},
    "phone": {},
    "email": {},
    "registrationForbidden": {},
    "tfaCode": {},
    "tfaType": {},
    "tfaMandatory": {},
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

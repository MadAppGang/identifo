import { ApiError, APIErrorCodes, IdentifoAuth, LoginResponse, TFAType, FederatedLoginProvider } from '@identifo/identifo-auth-js';
import { Component, Event, EventEmitter, getAssetPath, h, Host, Prop, State } from '@stencil/core';

export type Routes =
  | 'login'
  | 'register'
  | 'tfa/verify'
  | 'tfa/setup'
  | 'password/reset'
  | 'password/forgot'
  | 'callback'
  | 'otp/login'
  | 'error'
  | 'password/forgot/success'
  | 'logout'
  | 'loading';

const emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

@Component({
  tag: 'identifo-form',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: true,
})
export class IdentifoForm {
  @Prop({ mutable: true, reflect: true }) route: Routes = 'login';
  @Prop() token: string;
  @Prop({ reflect: true }) appId: string;
  @Prop({ reflect: true }) url: string;
  @Prop() theme: 'dark' | 'light' = 'light';
  @Prop() scopes: string = '';

  // This url will be preserved when federated login will be completed
  @Prop() callbackUrl: string;
  // Used for redirect on federated login flow
  // default:window.location.origin + window.location.pathname
  @Prop() federatedRedirectUrl: string;
  // Url used to redirect after logout
  @Prop() postLogoutRedirectUri: string;

  @Prop() debug: boolean;

  @State() auth: IdentifoAuth;

  @State() username: string;
  @State() password: string;
  @State() phone: string;
  @State() email: string;
  @State() registrationForbidden: boolean;
  @State() tfaCode: string;
  @State() tfaType: TFAType;
  @State() federatedProviders: string[];
  @State() tfaMandatory: boolean;
  @State() provisioningURI: string;
  @State() provisioningQR: string;
  @State() success: boolean;

  @State() lastError: ApiError;
  @State() lastResponse: LoginResponse;

  @Event() complete: EventEmitter<LoginResponse>;
  @Event() error: EventEmitter<ApiError>;

  // /**
  //  * The last name
  //  */
  // @Prop() last: string;

  // private getText(): string {
  //   return format(this.first, this.middle, this.last);
  // }
  processError(e: ApiError) {
    this.lastError = e;
    this.error.emit(e);
  }
  afterLoginRedirect = (e: LoginResponse) => {
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
  loginCatchRedirect = (data: ApiError): 'tfa/setup' => {
    if (data.id === APIErrorCodes.PleaseEnableTFA) {
      return 'tfa/setup';
    }
    throw data;
  };
  async signIn() {
    await this.auth.api
      .login(this.email, this.password, '', this.scopes.split(','))
      .then(this.afterLoginRedirect)
      .catch(this.loginCatchRedirect)
      .then(route => this.openRoute(route))
      .catch(e => this.processError(e));
  }
  async loginWith(provider: FederatedLoginProvider) {
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
    this.auth.api
      .verifyTFA(this.tfaCode, [])
      .then(() => this.openRoute('callback'))
      .catch(e => this.processError(e));
  }
  async setupTFA() {
    if (this.tfaType == TFAType.TFATypeSMS) {
      try {
        await this.auth.api.updateUser({ new_phone: this.phone });
      } catch (e) {
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
  openRoute(route: Routes) {
    this.lastError = undefined;
    this.route = route;
  }
  usernameChange(event: InputEvent) {
    this.username = (event.target as HTMLInputElement).value;
  }
  passwordChange(event: InputEvent) {
    this.password = (event.target as HTMLInputElement).value;
  }
  emailChange(event: InputEvent) {
    this.email = (event.target as HTMLInputElement).value;
  }
  phoneChange(event: InputEvent) {
    this.phone = (event.target as HTMLInputElement).value;
  }
  tfaCodeChange(event: InputEvent) {
    this.tfaCode = (event.target as HTMLInputElement).value;
  }
  validateEmail(email: string) {
    if (!emailRegex.test(email)) {
      this.processError({ detailedMessage: 'Email address is not valid', name: 'Validation error', message: 'Email address is not valid' });
      return false;
    }
    return true;
  }
  renderRoute(route: Routes) {
    switch (route) {
      case 'login':
        return (
          <div class="login-form">
            {!this.registrationForbidden && (
              <p class="login-form__register-text">
                Don't have an account?
                <a onClick={() => this.openRoute('register')} class="login-form__register-link">
                  {' '}
                  Sign Up
                </a>
              </p>
            )}
            <input
              type="text"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="floatingInput"
              value={this.email}
              placeholder="Email"
              onInput={event => this.emailChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn()}
            />
            <input
              type="password"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="floatingPassword"
              value={this.password}
              placeholder="Password"
              onInput={event => this.passwordChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage}
              </div>
            )}

            <div class={`login-form__buttons ${!!this.lastError ? 'login-form__buttons_mt-32' : ''}`}>
              <button onClick={() => this.signIn()} class="primary-button" disabled={!this.email || !this.password}>
                Login
              </button>
              <a onClick={() => this.openRoute('password/forgot')} class="login-form__forgot-pass">
                Forgot password
              </a>
            </div>
            {this.federatedProviders.length > 0 && (
              <div class="social-buttons">
                <p class="social-buttons__text">or continue with</p>
                <div class="social-buttons__social-medias">
                  {this.federatedProviders.indexOf('apple') > -1 && (
                    <div class="social-buttons__media" onClick={() => this.loginWith('apple')}>
                      <img src={getAssetPath(`assets/images/${'apple.svg'}`)} class="social-buttons__image" alt="login via apple" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('google') > -1 && (
                    <div class="social-buttons__media" onClick={() => this.loginWith('google')}>
                      <img src={getAssetPath(`assets/images/${'google.svg'}`)} class="social-buttons__image" alt="login via google" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('facebook') > -1 && (
                    <div class="social-buttons__media" onClick={() => this.loginWith('facebook')}>
                      <img src={getAssetPath(`assets/images/${'fb.svg'}`)} class="social-buttons__image" alt="login via facebook" />
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        );
      case 'register':
        return (
          <div class="register-form">
            <input
              type="text"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="floatingInput"
              value={this.email}
              placeholder="Email"
              onInput={event => this.emailChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp()}
            />
            <input
              type="password"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="floatingPassword"
              value={this.password}
              placeholder="Password"
              onInput={event => this.passwordChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage}
              </div>
            )}

            <div class={`register-form__buttons ${!!this.lastError ? 'register-form__buttons_mt-32' : ''}`}>
              <button onClick={() => this.signUp()} class="primary-button" disabled={!this.email || !this.password}>
                Continue
              </button>
              <a onClick={() => this.openRoute('login')} class="register-form__log-in">
                Go back to login
              </a>
            </div>
          </div>
        );
      case 'otp/login':
        return (
          <div class="otp-login">
            registrationForbidden
            {this.registrationForbidden === false && (
              <p class="otp-login__register-text">
                Don't have an account?
                <a onClick={() => this.openRoute('register')} class="login-form__register-link">
                  {' '}
                  Sign Up
                </a>
              </p>
            )}
            <input type="phone" class="form-control" id="floatingInput" value={this.phone} placeholder="Phone number" onInput={event => this.phoneChange(event as InputEvent)} />
            <button onClick={() => this.openRoute('tfa/verify')} class="primary-button" disabled={!this.phone}>
              Continue
            </button>
            {this.federatedProviders.length > 0 && (
              <div class="social-buttons">
                <p class="social-buttons__text">or continue with</p>
                <div class="social-buttons__social-medias">
                  {this.federatedProviders.indexOf('apple') > -1 && (
                    <div class="social-buttons__media" onClick={() => this.loginWith('apple')}>
                      <img src={getAssetPath(`assets/images/${'apple.svg'}`)} class="social-buttons__image" alt="login via apple" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('google') > -1 && (
                    <div class="social-buttons__media" onClick={() => this.loginWith('google')}>
                      <img src={getAssetPath(`assets/images/${'google.svg'}`)} class="social-buttons__image" alt="login via google" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('facebook') > -1 && (
                    <div class="social-buttons__media" onClick={() => this.loginWith('facebook')}>
                      <img src={getAssetPath(`assets/images/${'fb.svg'}`)} class="social-buttons__image" alt="login via facebook" />
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        );
      case 'tfa/setup':
        return (
          <div class="tfa-setup">
            <p class="tfa-setup__text">Protect your account with 2-step verification</p>
            {this.tfaType === TFAType.TFATypeApp && (
              <div class="info-card">
                <div class="info-card__controls">
                  <p class="info-card__title">Authenticator app</p>
                  <button type="button" class="info-card__button" onClick={() => this.setupTFA()}>
                    Setup
                  </button>
                </div>
                <p class="info-card__text">Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone.</p>
              </div>
            )}
            {this.tfaType === TFAType.TFATypeEmail && (
              <div class="info-card">
                <div class="info-card__controls">
                  <p class="info-card__title">Email</p>
                  <button type="button" class="info-card__button" onClick={() => this.setupTFA()}>
                    Setup
                  </button>
                </div>
                <p class="info-card__subtitle">{this.email}</p>
                <p class="info-card__text"> Use email as 2fa, please check your email, we will send confirmation code to this email.</p>
              </div>
            )}
            {this.tfaType === TFAType.TFATypeSMS && (
              <div class="tfa-setup__form">
                <p class="tfa-setup__subtitle"> Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone</p>
                <input
                  type="phone"
                  class={`form-control ${this.lastError && 'form-control-danger'}`}
                  id="floatingInput"
                  value={this.phone}
                  placeholder="Phone"
                  onInput={event => this.phoneChange(event as InputEvent)}
                  onKeyPress={e => !!(e.key === 'Enter' && this.phone) && this.setupTFA()}
                />

                {!!this.lastError && (
                  <div class="error" role="alert">
                    {this.lastError?.detailedMessage}
                  </div>
                )}

                <button onClick={() => this.setupTFA()} class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.phone}>
                  Setup phone
                </button>
              </div>
            )}
          </div>
        );
      case 'tfa/verify':
        return (
          <div class="tfa-verify">
            {!!(this.tfaType === TFAType.TFATypeApp) && (
              <div class="tfa-verify__title-wrapper">
                <h2 class={this.provisioningURI ? 'tfa-verify__title' : 'tfa-verify__title_mb-40'}>
                  {!!this.provisioningURI ? 'Please scan QR-code with the app' : 'Use GoogleAuth as 2fa'}
                </h2>
                {!!this.provisioningURI && <img src={`data:image/png;base64, ${this.provisioningQR}`} alt={this.provisioningURI} class="tfa-verify__qr-code" />}
              </div>
            )}
            {!!(this.tfaType === TFAType.TFATypeSMS) && (
              <div class="tfa-verify__title-wrapper">
                <h2 class="tfa-verify__title">Enter the code sent to your phone number</h2>
                <p class="tfa-verify__subtitle">The code has been sent to {this.phone}</p>
              </div>
            )}
            {!!(this.tfaType === TFAType.TFATypeEmail) && (
              <div class="tfa-verify__title-wrapper">
                <h2 class="tfa-verify__title">Enter the code sent to your email address</h2>
                <p class="tfa-verify__subtitle">The email has been sent to {this.email}</p>
              </div>
            )}
            <input
              type="text"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="floatingCode"
              value={this.tfaCode}
              placeholder="Verify code"
              onInput={event => this.tfaCodeChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.tfaCode) && this.verifyTFA()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage}
              </div>
            )}

            <button type="button" class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.tfaCode} onClick={() => this.verifyTFA()}>
              Confirm
            </button>
          </div>
        );
      case 'password/forgot':
        return (
          <div class="forgot-password">
            <h2 class="forgot-password__title">Enter the email you gave when you registered</h2>
            <p class="forgot-password__subtitle">We will send you a link to create a new password on email</p>
            <input
              type="email"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="floatingEmail"
              value={this.email}
              placeholder="Email"
              onInput={event => this.emailChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.email) && this.restorePassword()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage}
              </div>
            )}

            <button type="button" class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.email} onClick={() => this.restorePassword()}>
              Send the link
            </button>
          </div>
        );
      case 'password/forgot/success':
        return (
          <div class="forgot-password-success">
            {this.theme === 'dark' && <img src={getAssetPath(`./assets/images/${'email-dark.svg'}`)} alt="email" class="forgot-password-success__image" />}
            {this.theme === 'light' && <img src={getAssetPath(`./assets/images/${'email.svg'}`)} alt="email" class="forgot-password-success__image" />}
            <p class="forgot-password-success__text">We sent you an email with a link to create a new password</p>
          </div>
        );
      case 'password/reset':
        return (
          <div class="reset-password">
            <h2 class="reset-password__title">Set up a new password to log in to the website</h2>
            <p class="reset-password__subtitle">Memorize your password and do not give it to anyone.</p>
            <input
              type="password"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="floatingPassword"
              value={this.password}
              placeholder="Password"
              onInput={event => this.passwordChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.password) && this.setNewPassword()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage}
              </div>
            )}

            <button type="button" class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.password} onClick={() => this.setNewPassword()}>
              Save password
            </button>
          </div>
        );
      case 'error':
        return (
          <div class="error-view">
            <div class="error-view__message">{this.lastError.message}</div>
            <div class="error-view__details">{this.lastError.detailedMessage}</div>
          </div>
        );
      case 'callback':
        return (
          <div class="error-view">
            <div>Success</div>
            {this.debug && (
              <div>
                <div>Access token: {this.lastResponse.access_token}</div>
                <div>Refresh token: {this.lastResponse.refresh_token}</div>
                <div>User: {JSON.stringify(this.lastResponse.user)}</div>
              </div>
            )}
          </div>
        );
      case 'loading':
        return (
          <div class="error-view">
            <div>Loading ...</div>
          </div>
        );
    }
  }

  async componentWillLoad() {
    const postLogoutRedirectUri = this.postLogoutRedirectUri || window.location.origin + window.location.pathname;

    this.auth = new IdentifoAuth({ appId: this.appId, url: this.url, postLogoutRedirectUri });

    try {
      const settings = await this.auth.api.getAppSettings();
      this.registrationForbidden = settings.registrationForbidden;
      this.tfaType = settings.tfaType;
      this.federatedProviders = settings.federatedProviders;
    } catch (err) {
      this.route = 'error';
      this.lastError = err as ApiError;
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
  }

  render() {
    return (
      <Host>
        <div class={{ 'wrapper': this.theme === 'light', 'wrapper-dark': this.theme === 'dark' }}>{this.renderRoute(this.route)}</div>
        <div class="error-view">
          {this.debug && (
            <div>
              <br />
              {this.appId}
            </div>
          )}
        </div>
      </Host>
    );
  }
}

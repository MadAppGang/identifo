import { ApiError, APIErrorCodes, IdentifoAuth, LoginResponse, TFAType, FederatedLoginProvider, TFAStatus } from '@identifo/identifo-auth-js';
import { Component, Event, EventEmitter, getAssetPath, h, Host, Prop, State } from '@stencil/core';

const routes = [
  'login',
  'register',
  'tfa/verify/sms',
  'tfa/verify/email',
  'tfa/verify/app',
  'tfa/verify/select',
  'tfa/setup/sms',
  'tfa/setup/email',
  'tfa/setup/app',
  'tfa/setup/select',
  'password/reset',
  'password/forgot',
  'password/forgot/tfa/sms',
  'password/forgot/tfa/email',
  'password/forgot/tfa/app',
  'password/forgot/tfa/select',
  'callback',
  'otp/login',
  'error',
  'password/forgot/success',
  'logout',
  'loading',
] as const;

export type Routes = typeof routes[number];

export type TFASetupRoutes = 'tfa/setup/select' | 'tfa/setup/sms' | 'tfa/setup/email' | 'tfa/setup/app';
export type TFALoginVerifyRoutes = 'tfa/verify/select' | 'tfa/verify/sms' | 'tfa/verify/email' | 'tfa/verify/app';
export type TFAResetVerifyRoutes = 'password/forgot/tfa/select' | 'password/forgot/tfa/sms' | 'password/forgot/tfa/email' | 'password/forgot/tfa/app';
const emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

@Component({
  tag: 'identifo-form',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoForm {
  @Prop({ mutable: true, reflect: true }) route: Routes = 'login';
  @Prop() token: string;
  @Prop({ reflect: true }) appId: string;
  @Prop({ reflect: true }) url: string;
  @Prop() theme: 'dark' | 'light' | 'auto' = 'auto';
  @Prop() scopes: string = '';

  // This url will be preserved when federated login will be completed
  @Prop() callbackUrl: string;
  // Used for redirect on federated login flow
  // default:window.location.origin + window.location.pathname
  @Prop() federatedRedirectUrl: string;
  // Url used to redirect after logout
  @Prop() postLogoutRedirectUri: string;

  @Prop() debug: boolean;

  @State() selectedTheme: 'dark' | 'light' = 'light';

  @State() auth: IdentifoAuth;

  @State() username: string;
  @State() password: string;
  @State() phone: string;
  @State() email: string;
  @State() registrationForbidden: boolean;
  @State() tfaCode: string;
  @State() tfaTypes: TFAType[];
  @State() federatedProviders: string[] = [];
  @State() tfaStatus: TFAStatus;
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
    e.detailedMessage = e.detailedMessage?.trim();
    e.message = e.message?.trim();
    this.lastError = e;
    this.error.emit(e);
  }
  redirectTfa(prefix: string): string {
    if (this.tfaTypes.length === 1) {
      return `${prefix}/${this.tfaTypes[0]}`;
    } else {
      return `${prefix}/select`;
    }
  }
  afterLoginRedirect = (e: LoginResponse) => {
    this.phone = e.user.phone || '';
    this.email = e.user.email || '';
    this.lastResponse = e;
    if (e.require_2fa) {
      if (!e.enabled_2fa) {
        return this.redirectTfa('tfa/setup') as TFASetupRoutes;
      }
      if (e.enabled_2fa) {
        return this.redirectTfa('tfa/verify') as TFALoginVerifyRoutes;
      }
    }
    // Ask about tfa on login only
    if (this.tfaStatus === TFAStatus.OPTIONAL && ['login', 'login/otp', 'register'].includes(this.route)) {
      return `tfa/setup/select`;
    }
    if (e.access_token && e.refresh_token) {
      return 'callback';
    }
    if (e.access_token && !e.refresh_token) {
      return 'callback';
    }
  };
  loginCatchRedirect = (data: ApiError): TFASetupRoutes => {
    if (data.id === APIErrorCodes.PleaseEnableTFA) {
      return this.redirectTfa('tfa/setup') as TFASetupRoutes;
    }
    throw data;
  };
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
    if (this.route.indexOf('password/forgot/tfa') === 0) {
      this.auth.api
        .requestResetPassword(this.email, this.tfaCode)
        .then(() => {
          this.success = true;
          this.openRoute('password/forgot/success');
        })
        .catch(e => this.processError(e));
    } else {
      this.auth.api
        .verifyTFA(this.tfaCode, [])
        .then(this.afterLoginRedirect)
        .catch(this.loginCatchRedirect)
        .then(route => this.openRoute(route))
        .catch(e => this.processError(e));
    }
  }
  async selectTFA(type: TFAType) {
    this.openRoute(`tfa/setup/${type}` as TFASetupRoutes);
  }
  async setupTFA(type: TFAType) {
    switch (type) {
      case TFAType.TFATypeApp:
        break;
      case TFAType.TFATypeEmail:
        await this.auth.api.enableTFA();
        break;
      case TFAType.TFATypeSMS:
        try {
          await this.auth.api.updateUser({ new_phone: this.phone });
        } catch (e) {
          this.processError(e);
          return;
        }
        await this.auth.api.enableTFA();

        break;
    }
    this.openRoute(`tfa/verify/${type}` as TFALoginVerifyRoutes);
  }
  restorePassword() {
    this.auth.api
      .requestResetPassword(this.email)
      .then(response => {
        if (response.result === 'tfa-required') {
          this.openRoute(this.redirectTfa('password/forgot/tfa') as TFALoginVerifyRoutes);
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
  renderBackToLogin() {
    return (
      <a onClick={() => this.openRoute('login')} class="forgot-password__login">
        Go back to login
      </a>
    );
  }
  renderRoute(route: Routes) {
    switch (route) {
      case 'login':
        return (
          <div class="login-form">
            {!this.registrationForbidden && (
              <p class="login-form__register-text">
                Don't have an account?&nbsp;
                <a onClick={() => this.openRoute('register')} class="login-form__register-link">
                  Sign Up
                </a>
              </p>
            )}
            <input
              type="text"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="login"
              value={this.email}
              placeholder="Email"
              onInput={event => this.emailChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn()}
            />
            <input
              type="password"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="password"
              value={this.password}
              placeholder="Password"
              onInput={event => this.passwordChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.email && this.password) && this.signIn()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.message || this.lastError?.detailedMessage}
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
                    <div class="social-buttons__media social-buttons__apple" onClick={() => this.loginWith('apple')}>
                      <img src={getAssetPath(`assets/images/${'apple.svg'}`)} class="social-buttons__image" alt="login via apple" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('google') > -1 && (
                    <div class="social-buttons__media social-buttons__google" onClick={() => this.loginWith('google')}>
                      <img src={getAssetPath(`assets/images/${'google.svg'}`)} class="social-buttons__image" alt="login via google" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('facebook') > -1 && (
                    <div class="social-buttons__media social-buttons__facebook" onClick={() => this.loginWith('facebook')}>
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
              id="login"
              value={this.email}
              placeholder="Email"
              onInput={event => this.emailChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp()}
            />
            <input
              type="password"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="password"
              value={this.password}
              placeholder="Password"
              onInput={event => this.passwordChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage || this.lastError?.message}
              </div>
            )}

            <div class={`register-form__buttons ${!!this.lastError ? 'register-form__buttons_mt-32' : ''}`}>
              <button onClick={() => this.signUp()} class="primary-button" disabled={!this.email || !this.password}>
                Continue
              </button>
              {this.renderBackToLogin()}
            </div>
          </div>
        );
      case 'otp/login':
        return (
          <div class="otp-login">
            {!this.registrationForbidden && (
              <p class="otp-login__register-text">
                Don't have an account?&nbsp;
                <a onClick={() => this.openRoute('register')} class="login-form__register-link">
                  Sign Up
                </a>
              </p>
            )}
            <input type="phone" class="form-control" id="login" value={this.phone} placeholder="Phone number" onInput={event => this.phoneChange(event as InputEvent)} />
            <button onClick={() => this.openRoute(this.redirectTfa('tfa/verify') as TFALoginVerifyRoutes)} class="primary-button" disabled={!this.phone}>
              Continue
            </button>
            {this.federatedProviders.length > 0 && (
              <div class="social-buttons">
                <p class="social-buttons__text">or continue with</p>
                <div class="social-buttons__social-medias">
                  {this.federatedProviders.indexOf('apple') > -1 && (
                    <div class="social-buttons__media social-buttons__apple" onClick={() => this.loginWith('apple')}>
                      <img src={getAssetPath(`assets/images/${'apple.svg'}`)} class="social-buttons__image" alt="login via apple" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('google') > -1 && (
                    <div class="social-buttons__media social-buttons__google" onClick={() => this.loginWith('google')}>
                      <img src={getAssetPath(`assets/images/${'google.svg'}`)} class="social-buttons__image" alt="login via google" />
                    </div>
                  )}
                  {this.federatedProviders.indexOf('facebook') > -1 && (
                    <div class="social-buttons__media social-buttons__facebook" onClick={() => this.loginWith('facebook')}>
                      <img src={getAssetPath(`assets/images/${'fb.svg'}`)} class="social-buttons__image" alt="login via facebook" />
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        );
      case 'tfa/verify/select':
      case 'tfa/setup/select':
      case 'password/forgot/tfa/select':
        return (
          <div class="tfa-setup">
            {this.route === 'tfa/verify/select' && <p class="tfa-setup__text">Select 2-step verification method</p>}
            {this.route === 'tfa/setup/select' && <p class="tfa-setup__text">Protect your account with 2-step verification</p>}

            {this.tfaTypes.includes(TFAType.TFATypeApp) && (
              <div class="info-card info-card-app">
                <div class="info-card__controls">
                  <p class="info-card__title">Authenticator app</p>
                  <button type="button" class="info-card__button" onClick={() => this.selectTFA(TFAType.TFATypeApp)}>
                    Setup
                  </button>
                </div>
                <p class="info-card__text">Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone.</p>
              </div>
            )}
            {this.tfaTypes.includes(TFAType.TFATypeEmail) && (
              <div class="info-card info-card-email">
                <div class="info-card__controls">
                  <p class="info-card__title">Email</p>
                  <button type="button" class="info-card__button" onClick={() => this.selectTFA(TFAType.TFATypeEmail)}>
                    Setup
                  </button>
                </div>
                <p class="info-card__subtitle">{this.email}</p>
                <p class="info-card__text"> Use email as 2fa, please check your email, we will send confirmation code to this email.</p>
              </div>
            )}
            {this.tfaTypes.includes(TFAType.TFATypeSMS) && (
              <div class="info-card info-card-sms">
                <div class="info-card__controls">
                  <p class="info-card__title">SMS</p>
                  <button type="button" class="info-card__button" onClick={() => this.selectTFA(TFAType.TFATypeSMS)}>
                    Setup
                  </button>
                </div>
                <p class="info-card__subtitle">{this.phone}</p>
                <p class="info-card__text"> Use phone as 2fa, please check your phone, we will send confirmation code to this phone</p>
              </div>
            )}
            {this.route === 'tfa/setup/select' && this.tfaStatus === TFAStatus.OPTIONAL && (
              <a onClick={() => this.openRoute('callback')} class="forgot-password__login">
                Setup next time
              </a>
            )}
            {this.tfaStatus !== TFAStatus.OPTIONAL && this.renderBackToLogin()}
          </div>
        );
      case 'tfa/setup/email':
      case 'tfa/setup/sms':
      case 'tfa/setup/app':
        return (
          <div class="tfa-setup">
            <p class="tfa-setup__text">Protect your account with 2-step verification</p>
            {this.route === 'tfa/setup/app' && (
              <div class="tfa-setup__form">
                <p class="tfa-setup__subtitle">Please scan QR-code with the app and click Continue</p>
                <div class="tfa-setup__qr-wrapper">
                  {!!this.provisioningURI && <img src={`data:image/png;base64, ${this.provisioningQR}`} alt={this.provisioningURI} class="tfa-setup__qr-code" />}
                </div>
                <button onClick={() => this.setupTFA(TFAType.TFATypeApp)} class={`primary-button ${this.lastError && 'primary-button-mt-32'}`}>
                  Continue
                </button>
              </div>
            )}
            {this.route === 'tfa/setup/email' && (
              <div class="tfa-setup__form">
                <p class="tfa-setup__subtitle"> Use email as 2fa, please check your email bellow, we will send confirmation code to this email</p>
                <input
                  type="email"
                  class={`form-control ${this.lastError && 'form-control-danger'}`}
                  id="email"
                  value={this.email}
                  placeholder="Email"
                  onInput={event => this.emailChange(event as InputEvent)}
                  onKeyPress={e => !!(e.key === 'Enter' && this.email) && this.setupTFA(TFAType.TFATypeEmail)}
                />

                {!!this.lastError && (
                  <div class="error" role="alert">
                    {this.lastError?.detailedMessage || this.lastError?.message}
                  </div>
                )}

                <button onClick={() => this.setupTFA(TFAType.TFATypeEmail)} class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.email}>
                  Setup email
                </button>
              </div>
            )}
            {this.route === 'tfa/setup/sms' && (
              <div class="tfa-setup__form">
                <p class="tfa-setup__subtitle"> Use phone as 2fa, please check your phone bellow, we will send confirmation code to this phone</p>
                <input
                  type="phone"
                  class={`form-control ${this.lastError && 'form-control-danger'}`}
                  id="phone"
                  value={this.phone}
                  placeholder="Phone"
                  onInput={event => this.phoneChange(event as InputEvent)}
                  onKeyPress={e => !!(e.key === 'Enter' && this.phone) && this.setupTFA(TFAType.TFATypeSMS)}
                />

                {!!this.lastError && (
                  <div class="error" role="alert">
                    {this.lastError?.detailedMessage || this.lastError?.message}
                  </div>
                )}

                <button onClick={() => this.setupTFA(TFAType.TFATypeSMS)} class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.phone}>
                  Setup phone
                </button>
              </div>
            )}
            {this.renderBackToLogin()}
          </div>
        );
      case 'tfa/verify/app':
      case 'tfa/verify/email':
      case 'tfa/verify/sms':
      case 'password/forgot/tfa/app':
      case 'password/forgot/tfa/email':
      case 'password/forgot/tfa/sms':
        return (
          <div class="tfa-verify">
            {this.route.indexOf('app') > 0 && (
              <div class="tfa-verify__title-wrapper">
                <h2 class="tfa-verify__title">Enter the code from authenticator app</h2>
                <p class="tfa-verify__subtitle">Code will be generated by app</p>
              </div>
            )}
            {this.route.indexOf('sms') > 0 && (
              <div class="tfa-verify__title-wrapper">
                <h2 class="tfa-verify__title">Enter the code sent to your phone number</h2>
                <p class="tfa-verify__subtitle">The code has been sent to {this.phone}</p>
              </div>
            )}
            {this.route.indexOf('email') > 0 && (
              <div class="tfa-verify__title-wrapper">
                <h2 class="tfa-verify__title">Enter the code sent to your email address</h2>
                <p class="tfa-verify__subtitle">The email has been sent to {this.email}</p>
              </div>
            )}
            <input
              type="text"
              class={`form-control ${this.lastError && 'form-control-danger'}`}
              id="tfaCode"
              value={this.tfaCode}
              placeholder="Verify code"
              onInput={event => this.tfaCodeChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.tfaCode) && this.verifyTFA()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage || this.lastError?.message}
              </div>
            )}

            <button type="button" class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.tfaCode} onClick={() => this.verifyTFA()}>
              Confirm
            </button>
            {this.renderBackToLogin()}
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
              id="email"
              value={this.email}
              placeholder="Email"
              onInput={event => this.emailChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.email) && this.restorePassword()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage || this.lastError?.message}
              </div>
            )}

            <button type="button" class={`primary-button ${this.lastError && 'primary-button-mt-32'}`} disabled={!this.email} onClick={() => this.restorePassword()}>
              Send the link
            </button>
            {this.renderBackToLogin()}
          </div>
        );
      case 'password/forgot/success':
        return (
          <div class="forgot-password-success">
            {this.selectedTheme === 'dark' && <img src={getAssetPath(`./assets/images/${'email-dark.svg'}`)} alt="email" class="forgot-password-success__image" />}
            {this.selectedTheme === 'light' && <img src={getAssetPath(`./assets/images/${'email.svg'}`)} alt="email" class="forgot-password-success__image" />}
            <p class="forgot-password-success__text">We sent you an email with a link to create a new password</p>

            {this.renderBackToLogin()}
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
              id="password"
              value={this.password}
              placeholder="Password"
              onInput={event => this.passwordChange(event as InputEvent)}
              onKeyPress={e => !!(e.key === 'Enter' && this.password) && this.setNewPassword()}
            />

            {!!this.lastError && (
              <div class="error" role="alert">
                {this.lastError?.detailedMessage || this.lastError?.message}
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
                <div>
                  Access token: <div id="access_token">{this.lastResponse.access_token}</div>
                </div>
                <div>
                  Refresh token: <div id="refresh_token">{this.lastResponse.refresh_token}</div>
                </div>
                <div>
                  User: <div id="user_data">{JSON.stringify(this.lastResponse.user)}</div>
                </div>
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
    const params = new URLSearchParams(window.location.search);
    this.callbackUrl = params.get('callback-url') || params.get('callbackUrl') || params.get('callback_url') || '';
    this.scopes = params.get('scopes') || '';
    this.token = params.get('token');
    const paramRoute = params.get('route');
    if (routes.includes(paramRoute as Routes)) {
      this.route = paramRoute as Routes;
    } else {
      this.route = 'login';
    }

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
      const settings = await this.auth.api.getAppSettings(this.callbackUrl);
      this.registrationForbidden = settings.registrationForbidden;
      this.tfaTypes = Array.isArray(settings.tfaType) ? settings.tfaType : [settings.tfaType];
      this.tfaStatus = settings.tfaStatus;
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
    return (
      <Host>
        <div class={{ 'wrapper': this.selectedTheme === 'light', 'wrapper-dark': this.selectedTheme === 'dark' }}>{this.renderRoute(this.route)}</div>
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

import { StateLogin } from '@identifo/identifo-auth-js';
import { Component, getAssetPath, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-login',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormLogin {
  @State() email: string;
  @State() password: string;
  @State() state: StateLogin;
  @State() remember: boolean = false;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateLogin));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  passwordChange(event: InputEvent) {
    this.password = (event.target as HTMLInputElement).value;
  }
  emailChange(event: InputEvent) {
    this.email = (event.target as HTMLInputElement).value;
  }
  rememberChange(event: InputEvent) {
    this.remember = (event.target as HTMLInputElement).checked;
  }

  signin() {
    this.state.signin(this.email, this.password, this.remember);
  }

  render() {
    return (
      <div class="login-form">
        {!this.state.registrationForbidden && (
          <p class="login-form__register-text">
            Don't have an account?&nbsp;
            <a onClick={() => this.state.signup()} class="login-form__register-link">
              Sign Up
            </a>
          </p>
        )}
        <input
          type="text"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="login"
          value={this.email}
          placeholder="Email"
          onInput={event => this.emailChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.email && this.password) && this.signin()}
        />
        <input
          type="password"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="password"
          value={this.password}
          placeholder="Password"
          onInput={event => this.passwordChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.email && this.password) && this.signin()}
        />
        <label class="login-form__remember-me form-checkbox" htmlFor="remember">
          <input type="checkbox" class="form-control" id="remember" checked={this.remember} onInput={event => this.rememberChange(event as InputEvent)} />
          <span>Remember me</span>
        </label>

        {!!this.state.error && (
          <div class="error" role="alert">
            {this.state.error?.message || this.state.error?.detailedMessage}
          </div>
        )}

        <div class={`login-form__buttons ${!!this.state.error ? 'login-form__buttons_mt-32' : ''}`}>
          <button onClick={() => this.signin()} class="primary-button" disabled={!this.email || !this.password}>
            Login
          </button>
          <a onClick={() => this.state.passwordForgot()} class="login-form__forgot-pass">
            Forgot password
          </a>
        </div>
        {this.state.federatedProviders?.length > 0 && (
          <div class="social-buttons">
            <p class="social-buttons__text">or continue with</p>
            <div class="social-buttons__social-medias">
              {this.state.federatedProviders.indexOf('apple') > -1 && (
                <div class="social-buttons__media social-buttons__apple" onClick={() => this.state.socialLogin('apple')}>
                  <img src={getAssetPath(`assets/images/${'apple.svg'}`)} class="social-buttons__image" alt="login via apple" />
                </div>
              )}
              {this.state.federatedProviders.indexOf('google') > -1 && (
                <div class="social-buttons__media social-buttons__google" onClick={() => this.state.socialLogin('google')}>
                  <img src={getAssetPath(`assets/images/${'google.svg'}`)} class="social-buttons__image" alt="login via google" />
                </div>
              )}
              {this.state.federatedProviders.indexOf('facebook') > -1 && (
                <div class="social-buttons__media social-buttons__facebook" onClick={() => this.state.socialLogin('facebook')}>
                  <img src={getAssetPath(`assets/images/${'fb.svg'}`)} class="social-buttons__image" alt="login via facebook" />
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    );
  }
}

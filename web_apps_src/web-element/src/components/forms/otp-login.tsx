import { StateOTPLogin } from '@identifo/identifo-auth-js';
import { Component, getAssetPath, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-otp-login',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormOtpLogin {
  @State() state: StateOTPLogin;
  @State() phone: string;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateOTPLogin));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  phoneChange(event: InputEvent) {
    this.phone = (event.target as HTMLInputElement).value;
  }

  signin() {
    this.state.signin(this.phone);
  }

  render() {
    return (
      <div class="otp-login">
        {!this.state.registrationForbidden && (
          <p class="otp-login__register-text">
            Don't have an account?&nbsp;
            <a onClick={() => this.state.signup()} class="login-form__register-link">
              Sign Up
            </a>
          </p>
        )}
        <input type="phone" class="form-control" id="login" value={this.phone} placeholder="Phone number" onInput={event => this.phoneChange(event as InputEvent)} />
        <button onClick={() => this.signin()} class="primary-button" disabled={!this.phone}>
          Continue
        </button>
        {this.state.federatedProviders.length > 0 && (
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

import { StateLogin, StateLoginPhone } from '@identifo/identifo-auth-js';
import { Component, getAssetPath, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-login-ways',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormLoginWays {
  @State() state: StateLogin | StateLoginPhone;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateLogin | StateLoginPhone));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    return (
      <div>
        {this.state.loginTypes && Object.values(this.state.loginTypes).length > 0 && (
          <div class="login-types">
            {Object.values(this.state.loginTypes).map(t => (
              <a onClick={() => t.click()} class="login-type">
                {t.type === 'phone' && 'login with phone'}
                {t.type === 'email' && 'login with password'}
              </a>
            ))}
          </div>
        )}
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

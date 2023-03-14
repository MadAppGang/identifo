import { StateLoginOidc } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';

import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-login-oidc',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormLoginOidc {
  @State() state: StateLoginOidc;
  @State() isVerify = false;

  subscription: Subscription;

  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateLoginOidc));
    const search = new URLSearchParams(window.location.search);
    const state = search.get('state');
    const code = search.get('code');
    if (state && code) {
      this.state.verify(state, code);
      this.isVerify = true;
    }
  }

  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  login = () => {
    const url = this.state.oidcLink;
    window.location.href = url.toString();
  };

  render() {
    return (
      <div class="oidc-login form-modal">
        <h2 class="form-modal__title">{this.isVerify ? 'Good, you’re almost in the system' : 'Sing In'}</h2>
        <p class="form-modal__subtitle">{this.isVerify ? 'After a moment you will be redirected to the dashboard' : 'To sign in, please tap on the button below'}</p>
        {!this.isVerify && (
          <div>
            <button onClick={this.login} class="primary-button">
              SIGN IN
            </button>
          </div>
        )}
      </div>
    );
  }
}

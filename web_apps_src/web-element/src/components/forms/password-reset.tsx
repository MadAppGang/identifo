import { StatePasswordReset } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-password-reset',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormError {
  @State() password: string;
  @State() state: StatePasswordReset;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StatePasswordReset));
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');
    if (token) {
      CDKService.cdk.auth.tokenService.saveToken(token, 'access');
    }
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  passwordChange(event: InputEvent) {
    this.password = (event.target as HTMLInputElement).value;
  }
  setNewPassword() {
    this.state.setNewPassword(this.password);
  }

  render() {
    return (
      <div class="reset-password">
        <h2 class="reset-password__title">Set up a new password to log in to the website</h2>
        <p class="reset-password__subtitle">Memorize your password and do not give it to anyone.</p>
        <input
          type="password"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="password"
          value={this.password}
          placeholder="Password"
          onInput={event => this.passwordChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.password) && this.setNewPassword()}
        />
        <identifo-form-error-alert></identifo-form-error-alert>
        <button type="button" class={`primary-button ${this.state.error && 'primary-button-mt-32'}`} disabled={!this.password} onClick={() => this.setNewPassword()}>
          Save password
        </button>
      </div>
    );
  }
}

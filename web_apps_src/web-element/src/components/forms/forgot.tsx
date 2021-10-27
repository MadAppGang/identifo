import { StatePasswordForgot } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-forgot',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormLogin {
  @State() email: string;
  @State() password: string;
  @State() state: StatePasswordForgot;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StatePasswordForgot));
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
  restorePassword() {
    this.state.restorePassword(this.email);
  }

  render() {
    return (
      <div class="forgot-password">
        <h2 class="forgot-password__title">Enter the email you gave when you registered</h2>
        <p class="forgot-password__subtitle">We will send you a link to create a new password on email</p>
        <input
          type="email"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="email"
          value={this.email}
          placeholder="Email"
          onInput={event => this.emailChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.email) && this.restorePassword()}
        />
        <identifo-form-error-alert></identifo-form-error-alert>
        <button type="button" class={`primary-button ${this.state.error && 'primary-button-mt-32'}`} disabled={!this.email} onClick={() => this.restorePassword()}>
          Send the link
        </button>
        <identifo-form-goback></identifo-form-goback>
      </div>
    );
  }
}

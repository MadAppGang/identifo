import { StateRegister } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-register',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormRegister {
  @State() email: string;
  @State() password: string;
  @State() state: StateRegister;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateRegister));
    const params = new URLSearchParams(window.location.search);
    this.email = params.get('email') || '';
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
  signUp() {
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');

    this.state.signup(this.email, this.password, token);
  }

  render() {
    return (
      <div class="register-form">
        <input
          type="text"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="login"
          value={this.email}
          placeholder="Email"
          onInput={event => this.emailChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp()}
        />
        <input
          type="password"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="password"
          value={this.password}
          placeholder="Password"
          onInput={event => this.passwordChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.password && this.email) && this.signUp()}
        />
        <identifo-form-error-alert></identifo-form-error-alert>
        <div class={`register-form__buttons ${!!this.state.error ? 'register-form__buttons_mt-32' : ''}`}>
          <button onClick={() => this.signUp()} class="primary-button" disabled={!this.email || !this.password}>
            Continue
          </button>
          <identifo-form-goback></identifo-form-goback>
        </div>
      </div>
    );
  }
}

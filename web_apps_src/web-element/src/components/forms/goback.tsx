import { StateRegister, StatePasswordForgot, StatePasswordForgotSuccess } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-goback',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormGoback {
  @State() state: StateRegister | StatePasswordForgot | StatePasswordForgotSuccess;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateRegister | StatePasswordForgot | StatePasswordForgotSuccess));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    return (
      <a onClick={() => this.state.goback()} class="forgot-password__login">
        Go back to login
      </a>
    );
  }
}

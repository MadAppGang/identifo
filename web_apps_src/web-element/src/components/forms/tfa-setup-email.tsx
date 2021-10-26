import { StateTFASetupEmail } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-tfa-setup-email',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormTFASetupEmail {
  @State() state: StateTFASetupEmail;
  @State() email: string;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => {
      this.state = state as StateTFASetupEmail;
      this.email = this.state.email;
    });
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  emailChange(event: InputEvent) {
    this.email = (event.target as HTMLInputElement).value;
  }
  render() {
    return (
      <div class="tfa-setup__form">
        <p class="tfa-setup__subtitle"> Use email as 2fa, please check your email bellow, we will send confirmation code to this email</p>
        <input
          type="email"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="email"
          value={this.email}
          placeholder="Email"
          onInput={event => this.emailChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.email) && this.state.setupTFA(this.email)}
        />
        <identifo-form-error-alert></identifo-form-error-alert>
        <button onClick={() => this.state.setupTFA(this.email)} class={`primary-button ${this.state.error && 'primary-button-mt-32'}`} disabled={!this.email}>
          Setup email
        </button>
      </div>
    );
  }
}

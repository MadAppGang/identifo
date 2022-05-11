import { StateTFASetupSMS } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-tfa-setup-sms',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormTFASetupSMS {
  @State() state: StateTFASetupSMS;
  @State() phone: string;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => {
      this.state = state as StateTFASetupSMS;
      this.phone = this.state.phone;
    });
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  phoneChange(event: InputEvent) {
    this.phone = (event.target as HTMLInputElement).value;
  }

  render() {
    return (
      <div class="tfa-setup__form">
        <p class="tfa-setup__subtitle">
          Your phone will be used for 2-step verification. Please enter your phone number below and click 'Setup Phone'. We will send a confirmation code to the phone number you
          enter
        </p>
        <input
          type="phone"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="phone"
          value={this.phone}
          placeholder="Phone"
          onInput={event => this.phoneChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.phone) && this.state.setupTFA(this.phone)}
        />
        <identifo-form-error-alert></identifo-form-error-alert>
        <button onClick={() => this.state.setupTFA(this.phone)} class={`primary-button ${this.state.error && 'primary-button-mt-32'}`} disabled={!this.phone}>
          Setup phone
        </button>
      </div>
    );
  }
}

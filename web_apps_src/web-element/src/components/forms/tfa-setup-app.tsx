import { StateTFASetupApp } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-tfa-setup-app',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormTFASetupApp {
  @State() state: StateTFASetupApp;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateTFASetupApp));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    return (
      <div class="tfa-setup__form">
        <p class="tfa-setup__subtitle">Please scan QR-code with the app and click Continue</p>
        <div class="tfa-setup__qr-wrapper">
          {!!this.state.provisioningURI && <img src={`data:image/png;base64, ${this.state.provisioningQR}`} alt={this.state.provisioningURI} class="tfa-setup__qr-code" />}
        </div>
        <button onClick={() => this.state.setupTFA()} class={`primary-button ${this.state.error && 'primary-button-mt-32'}`}>
          Continue
        </button>
      </div>
    );
  }
}

import { Routes, StateTFASetupApp, StateTFASetupEmail, StateTFASetupSMS } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-tfa-setup',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormTFASetup {
  @State() state: StateTFASetupEmail | StateTFASetupApp | StateTFASetupSMS;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateTFASetupEmail | StateTFASetupApp | StateTFASetupSMS));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    return (
      <div class="tfa-setup">
        <p class="tfa-setup__text">Protect your account with 2-step verification</p>
        {this.state.route === Routes.TFA_SETUP_APP && <identifo-form-tfa-setup-app></identifo-form-tfa-setup-app>}
        {this.state.route === Routes.TFA_SETUP_EMAIL && <identifo-form-tfa-setup-email></identifo-form-tfa-setup-email>}
        {this.state.route === Routes.TFA_SETUP_SMS && <identifo-form-tfa-setup-sms></identifo-form-tfa-setup-sms>}
        <identifo-form-goback></identifo-form-goback>
      </div>
    );
  }
}

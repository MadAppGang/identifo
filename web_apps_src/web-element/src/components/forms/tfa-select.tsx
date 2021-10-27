import { Routes, StatePasswordForgotTFASelect, StateTFASetupSelect, StateTFAVerifySelect, TFAStatus, TFAType } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-tfa-select',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormTFASelect {
  @State() state: StateTFASetupSelect | StateTFAVerifySelect | StatePasswordForgotTFASelect;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateTFASetupSelect | StateTFAVerifySelect | StatePasswordForgotTFASelect));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    return (
      <div class="tfa-setup">
        {this.state.route === Routes.TFA_VERIFY_SELECT && <p class="tfa-setup__text">Select 2-step verification method</p>}
        {this.state.route === Routes.TFA_SETUP_SELECT && <p class="tfa-setup__text">Protect your account with 2-step verification</p>}

        {this.state.tfaTypes.includes(TFAType.TFATypeApp) && (
          <div class="info-card info-card-app">
            <div class="info-card__controls">
              <p class="info-card__title">Authenticator app</p>
              <button type="button" class="info-card__button" onClick={() => this.state.select(TFAType.TFATypeApp)}>
                Setup
              </button>
            </div>
            <p class="info-card__text">Use the Authenticator app to get free verification codes, even when your phone is offline. Available for Android and iPhone.</p>
          </div>
        )}
        {this.state.tfaTypes.includes(TFAType.TFATypeEmail) && (
          <div class="info-card info-card-email">
            <div class="info-card__controls">
              <p class="info-card__title">Email</p>
              <button type="button" class="info-card__button" onClick={() => this.state.select(TFAType.TFATypeEmail)}>
                Setup
              </button>
            </div>
            <p class="info-card__subtitle">{this.state.email}</p>
            <p class="info-card__text"> Use email as 2fa, please check your email, we will send confirmation code to this email.</p>
          </div>
        )}
        {this.state.tfaTypes.includes(TFAType.TFATypeSMS) && (
          <div class="info-card info-card-sms">
            <div class="info-card__controls">
              <p class="info-card__title">SMS</p>
              <button type="button" class="info-card__button" onClick={() => this.state.select(TFAType.TFATypeSMS)}>
                Setup
              </button>
            </div>
            <p class="info-card__subtitle">{this.state.phone}</p>
            <p class="info-card__text"> Use phone as 2fa, please check your phone, we will send confirmation code to this phone</p>
          </div>
        )}
        {this.state.route === Routes.TFA_SETUP_SELECT && (
          <div>
            {this.state.tfaStatus === TFAStatus.OPTIONAL && (
              <a onClick={() => (this.state as StateTFASetupSelect).setupNextTime()} class="forgot-password__login">
                Setup next time
              </a>
            )}
            {this.state.tfaStatus !== TFAStatus.OPTIONAL && <identifo-form-goback></identifo-form-goback>}
          </div>
        )}
      </div>
    );
  }
}

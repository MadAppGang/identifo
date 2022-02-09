import { StateLoginPhoneVerify } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-login-phone-verify',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormOtpLogin {
  @State() state: StateLoginPhoneVerify;
  @State() code: string;
  @State() resendLink: boolean;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => {
      this.state = state as StateLoginPhoneVerify;
      if (this.state.resendTimeout > 0) {
        window.setTimeout(() => (this.resendLink = true), (this.state as StateLoginPhoneVerify).resendTimeout);
      }
    });
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  codeChange(event: InputEvent) {
    this.code = (event.target as HTMLInputElement).value;
  }

  verify() {
    this.state.login(this.code);
  }

  resend() {
    this.state.resendCode();
  }

  render() {
    return (
      <div class="tfa-verify">
        <div class="tfa-verify__title-wrapper">
          <h2 class="tfa-verify__title">Enter the code sent to your phone number</h2>
          <p class="tfa-verify__subtitle">The code has been sent to {this.state.phone}</p>
        </div>

        <input
          type="text"
          class={`form-control ${this.state.error && 'form-control-danger'}`}
          id="code"
          value={this.code}
          placeholder="Verify code"
          onInput={event => this.codeChange(event as InputEvent)}
          onKeyPress={e => !!(e.key === 'Enter' && this.code) && this.verify()}
        />
        <identifo-form-error-alert></identifo-form-error-alert>
        <button type="button" class={`primary-button ${this.state.error && 'primary-button-mt-32'}`} disabled={!this.code} onClick={() => this.verify()}>
          Confirm
        </button>

        {this.resendLink && (
          <a onClick={() => this.resend()} class="forgot-password__login">
            Resend code
          </a>
        )}
        <identifo-form-goback></identifo-form-goback>
      </div>
    );
  }
}

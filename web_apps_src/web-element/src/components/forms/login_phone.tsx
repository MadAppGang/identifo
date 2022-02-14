import { StateLoginPhone } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-login-phone',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormOtpLogin {
  @State() state: StateLoginPhone;
  @State() phone: string;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateLoginPhone));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  phoneChange(event: InputEvent) {
    this.phone = (event.target as HTMLInputElement).value;
  }

  requestCode() {
    this.state.requestCode(this.phone);
  }

  render() {
    return (
      <div class="otp-login">
        <input type="phone" class="form-control" id="login" value={this.phone} placeholder="Phone number" onInput={event => this.phoneChange(event as InputEvent)} />
        <button onClick={() => this.requestCode()} class="primary-button" disabled={!this.phone}>
          Continue
        </button>
        <identifo-form-login-ways></identifo-form-login-ways>
      </div>
    );
  }
}

import { StatePasswordForgotSuccess } from '@identifo/identifo-auth-js';
import { Component, getAssetPath, h, Prop, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-forgot-success',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormLogin {
  @Prop() selectedTheme: 'dark' | 'light' = 'light';
  @State() email: string;
  @State() password: string;
  @State() state: StatePasswordForgotSuccess;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StatePasswordForgotSuccess));
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

  render() {
    return (
      <div class="forgot-password-success">
        {this.selectedTheme === 'dark' && <img src={getAssetPath(`./assets/images/${'email-dark.svg'}`)} alt="email" class="forgot-password-success__image" />}
        {this.selectedTheme === 'light' && <img src={getAssetPath(`./assets/images/${'email.svg'}`)} alt="email" class="forgot-password-success__image" />}
        <p class="forgot-password-success__text">We sent you an email with a link to create a new password</p>

        <identifo-form-goback></identifo-form-goback>
      </div>
    );
  }
}

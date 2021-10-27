import { StateWithError } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-error-alert',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormErrorAlert {
  @State() state: StateWithError;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateWithError));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    if (!!this.state.error) {
      return (
        <div class="error" role="alert">
          {this.state.error?.detailedMessage || this.state.error?.message}
        </div>
      );
    }
  }
}

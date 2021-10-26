import { StateError } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-error',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormError {
  @State() state: StateError;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateError));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    return (
      <div class="error-view">
        <div class="error-view__message">{this.state.error.message}</div>
        <div class="error-view__details">{this.state.error.detailedMessage}</div>
      </div>
    );
  }
}

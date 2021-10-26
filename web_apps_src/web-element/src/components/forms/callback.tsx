import { StateCallback } from '@identifo/identifo-auth-js';
import { Component, h, State } from '@stencil/core';
import { Subscription } from 'rxjs';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form-callback',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoFormCallback {
  @State() state: StateCallback;

  subscription: Subscription;
  connectedCallback() {
    this.subscription = CDKService.cdk.state.subscribe(state => (this.state = state as StateCallback));
  }
  disconnectedCallback() {
    this.subscription.unsubscribe();
  }

  render() {
    return (
      <div class="error-view">
        <div>Success</div>
        {CDKService.debug && (
          <div>
            <div>
              Access token: <div id="access_token">{this.state.result.access_token}</div>
            </div>
            <div>
              Refresh token: <div id="refresh_token">{this.state.result.refresh_token}</div>
            </div>
            <div>
              User: <div id="user_data">{JSON.stringify(this.state.result.user)}</div>
            </div>
          </div>
        )}
      </div>
    );
  }
}

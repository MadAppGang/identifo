import { StateCallback, States } from '@identifo/identifo-auth-js';
import { ApiError, CDK, LoginResponse, Routes } from '@identifo/identifo-auth-js';
import { Component, Event, EventEmitter, h, Host, Prop, State } from '@stencil/core';
import { CDKService } from '../../services/cdk.service';

@Component({
  tag: 'identifo-form',
  styleUrl: '../../styles/identifo-form/main.scss',
  assetsDirs: ['assets'],
  shadow: false,
})
export class IdentifoForm {
  @Prop({ mutable: true, reflect: true }) route: Routes;
  @Prop({ reflect: true }) appId: string;
  @Prop({ reflect: true }) url: string;
  @Prop() theme: 'dark' | 'light' | 'auto' = 'auto';

  // This url will be preserved when federated login will be completed
  @Prop() callbackUrl: string;

  // Used for redirect on federated login flow
  // default:window.location.origin + window.location.pathname
  @Prop() federatedRedirectUrl: string;

  // Url used to redirect after logout
  @Prop() postLogoutRedirectUri: string;

  @Prop() debug: boolean;

  @State() selectedTheme: 'dark' | 'light' = 'light';

  @State() identifoCdk: CDK;

  @Event() complete: EventEmitter<LoginResponse>;
  @Event() error: EventEmitter<ApiError>;

  renderRoute(route: Routes) {
    switch (route) {
      case Routes.LOGIN:
        return <identifo-form-login></identifo-form-login>;
      case Routes.LOGIN_PHONE:
        return <identifo-form-login-phone></identifo-form-login-phone>;
      case Routes.LOGIN_PHONE_VERIFY:
        return <identifo-form-login-phone-verify></identifo-form-login-phone-verify>;
      case Routes.ERROR:
        return <identifo-form-error></identifo-form-error>;
      case Routes.CALLBACK:
        return <identifo-form-callback></identifo-form-callback>;
      case Routes.REGISTER:
        return <identifo-form-register></identifo-form-register>;
      case Routes.PASSWORD_RESET:
        return <identifo-form-password-reset></identifo-form-password-reset>;
      case Routes.PASSWORD_FORGOT:
        return <identifo-form-forgot></identifo-form-forgot>;
      case Routes.PASSWORD_FORGOT_SUCCESS:
        return <identifo-form-forgot-success selectedTheme={this.selectedTheme}></identifo-form-forgot-success>;
      case Routes.TFA_SETUP_APP:
      case Routes.TFA_SETUP_EMAIL:
      case Routes.TFA_SETUP_SMS:
        return <identifo-form-tfa-setup></identifo-form-tfa-setup>;
      case Routes.TFA_SETUP_SELECT:
      case Routes.TFA_VERIFY_SELECT:
      case Routes.PASSWORD_FORGOT_TFA_SELECT:
        return <identifo-form-tfa-select></identifo-form-tfa-select>;
      case Routes.TFA_VERIFY_APP:
      case Routes.TFA_VERIFY_EMAIL:
      case Routes.TFA_VERIFY_SMS:
      case Routes.PASSWORD_FORGOT_TFA_APP:
      case Routes.PASSWORD_FORGOT_TFA_EMAIL:
      case Routes.PASSWORD_FORGOT_TFA_SMS:
        return <identifo-form-tfa-verify></identifo-form-tfa-verify>;
      case Routes.LOADING:
        return (
          <div class="error-view">
            <div>Loading ...</div>
          </div>
        );
    }
  }

  async componentWillLoad() {
    // Get url params and configure CDK serivice
    const params = new URLSearchParams(window.location.search);
    const callbackUrl = params.get('callback-url') || params.get('callbackUrl') || params.get('callback_url') || '';
    const scopes = (params.get('scopes') || '').split(',').map(s => s.trim());
    const postLogoutRedirectUri = this.postLogoutRedirectUri || window.location.origin + window.location.pathname;

    await CDKService.configure(
      {
        appId: this.appId,
        url: this.url,
        postLogoutRedirectUri,
        scopes,
      },
      callbackUrl,
    );

    CDKService.debug = this.debug;
    CDKService.cdk.state.subscribe((s: States) => {
      if (this.debug) {
        console.debug(s);
      }
      this.route = s.route;
      if (this.route === Routes.CALLBACK) {
        const u = new URL(window.location.href);
        u.searchParams.set('callbackUrl', (s as StateCallback).callbackUrl);
        window.history.replaceState({}, document.title, `${u.pathname}?${u.searchParams.toString()}`);
        this.complete.emit({ ...(s as StateCallback).result, callbackUrl: (s as StateCallback).callbackUrl });
      }
      if (this.route === Routes.LOGOUT) {
        this.complete.emit();
      }
    });

    if (CDKService.cdk.state.getValue().route !== Routes.ERROR) {
      switch (true) {
        case location.pathname.indexOf(Routes.REGISTER) > -1:
          CDKService.cdk.register();
          break;
        case location.pathname.indexOf(Routes.PASSWORD_RESET) > -1:
          CDKService.cdk.passwordReset();
          break;
        default:
          CDKService.cdk.login();
          break;
      }
    }

    // Auto theme select
    this.selectedTheme = 'light';
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      if (this.theme === 'auto') {
        this.selectedTheme = 'dark';
      }
    }
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
      if (this.theme === 'auto') {
        this.selectedTheme = e.matches ? 'dark' : 'light';
      }
    });
  }

  render() {
    return (
      <Host>
        <div class={{ 'wrapper': this.selectedTheme === 'light', 'wrapper-dark': this.selectedTheme === 'dark' }}>{this.renderRoute(this.route)}</div>
      </Host>
    );
  }
}

import { ApiError, FederatedLoginProvider, LoginResponse, TFAStatus, TFAType, SuccessResponse } from '../api/model';

export enum Routes {
  'LOGIN' = 'login',
  'REGISTER' = 'register',
  'TFA_VERIFY_SMS' = 'tfa/verify/sms',
  'TFA_VERIFY_EMAIL' = 'tfa/verify/email',
  'TFA_VERIFY_APP' = 'tfa/verify/app',
  'TFA_VERIFY_SELECT' = 'tfa/verify/select',
  'TFA_SETUP_SMS' = 'tfa/setup/sms',
  'TFA_SETUP_EMAIL' = 'tfa/setup/email',
  'TFA_SETUP_APP' = 'tfa/setup/app',
  'TFA_SETUP_SELECT' = 'tfa/setup/select',
  'PASSWORD_RESET' = 'password/reset',
  'PASSWORD_FORGOT' = 'password/forgot',
  'PASSWORD_FORGOT_TFA_SMS' = 'password/forgot/tfa/sms',
  'PASSWORD_FORGOT_TFA_EMAIL' = 'password/forgot/tfa/email',
  'PASSWORD_FORGOT_TFA_APP' = 'password/forgot/tfa/app',
  'PASSWORD_FORGOT_TFA_SELECT' = 'password/forgot/tfa/select',
  'CALLBACK' = 'callback',
  'OTP_LOGIN' = 'otp/login',
  'ERROR' = 'error',
  'PASSWORD_FORGOT_SUCCESS' = 'password/forgot/success',
  'LOGOUT' = 'logout',
  'LOADING' = 'loading',
}

export type TFASetupRoutes =
  | Routes.TFA_SETUP_SELECT
  | Routes.TFA_SETUP_SMS
  | Routes.TFA_SETUP_EMAIL
  | Routes.TFA_SETUP_APP;
export type TFALoginVerifyRoutes =
  | Routes.TFA_VERIFY_SELECT
  | Routes.TFA_VERIFY_SMS
  | Routes.TFA_VERIFY_EMAIL
  | Routes.TFA_VERIFY_APP;
export type TFAResetVerifyRoutes =
  | Routes.PASSWORD_FORGOT_TFA_SELECT
  | Routes.PASSWORD_FORGOT_TFA_SMS
  | Routes.PASSWORD_FORGOT_TFA_EMAIL
  | Routes.PASSWORD_FORGOT_TFA_APP;

export interface State {
  route: Routes;
}

export interface StateWithError {
  error: ApiError;
}

export interface StateLogin extends State, StateWithError {
  route: Routes.LOGIN;
  registrationForbidden: boolean;
  federatedProviders: FederatedLoginProvider[];
  signup: () => Promise<void>;
  signin: (email: string, password: string) => Promise<void>;
  socialLogin: (provider: FederatedLoginProvider) => Promise<void>;
  passwordForgot: () => Promise<void>;
}

export interface StateRegister extends State, StateWithError {
  route: Routes.REGISTER;
  signup: (email: string, password: string) => Promise<void>;
  goback: () => Promise<void>;
}

export interface StatePasswordForgot extends State, StateWithError {
  route: Routes.PASSWORD_FORGOT;
  restorePassword: (email: string) => Promise<void>;
  goback: () => Promise<void>;
}
export interface StatePasswordForgotSuccess extends State {
  route: Routes.PASSWORD_FORGOT_SUCCESS;
  goback: () => Promise<void>;
}

export interface StateError extends State, StateWithError {
  route: Routes.ERROR;
}

export interface StateCallback extends State {
  route: Routes.CALLBACK;
  callbackUrl: string;
  result: LoginResponse;
}

export interface StatePasswordReset extends State, StateWithError {
  route: Routes.PASSWORD_RESET;
  setNewPassword: (password: string) => Promise<void>;
}

export interface StateLoading extends State {
  route: Routes.LOADING;
}

export interface StateOTPLogin extends State {
  route: Routes.OTP_LOGIN;
  registrationForbidden: boolean;
  federatedProviders: FederatedLoginProvider[];
  signup: () => Promise<void>;
  signin: (phone: string) => Promise<void>;
  socialLogin: (provider: FederatedLoginProvider) => Promise<void>;
}

interface StateTFASetup extends State, StateWithError {}

export interface StateTFASetupApp extends StateTFASetup {
  route: Routes.TFA_SETUP_APP;
  provisioningURI: string;
  provisioningQR: string;
  setupTFA: () => Promise<void>;
}
export interface StateTFASetupEmail extends StateTFASetup {
  route: Routes.TFA_SETUP_EMAIL;
  email: string;
  setupTFA: (email: string) => Promise<void>;
}
export interface StateTFASetupSMS extends StateTFASetup {
  route: Routes.TFA_SETUP_SMS;
  phone: string;
  setupTFA: (phone: string) => Promise<void>;
}

interface StateTFASelect extends State {
  tfaTypes: TFAType[];
  select: (type: TFAType) => Promise<void>;
  email?: string;
  phone?: string;
}

export interface StateTFASetupSelect extends StateTFASelect {
  route: Routes.TFA_SETUP_SELECT;
  tfaStatus: TFAStatus;
  setupNextTime: () => Promise<void>;
}
export interface StateTFAVerifySelect extends StateTFASelect {
  route: Routes.TFA_VERIFY_SELECT;
}
export interface StatePasswordForgotTFASelect extends StateTFASelect {
  route: Routes.PASSWORD_FORGOT_TFA_SELECT;
}

export interface StateTFAVerifyApp extends State, StateWithError {
  route: Routes.TFA_VERIFY_APP;
  email?: string;
  phone?: string;
  verifyTFA: (code: string) => Promise<void>;
}

export interface StateTFAVerifyEmailSms extends State, StateWithError {
  route: Routes.TFA_VERIFY_EMAIL | Routes.TFA_VERIFY_SMS;
  email?: string;
  phone?: string;
  resendTimeout: number;
  verifyTFA: (code: string) => Promise<void>;
  resendTFA: () => Promise<void>;
}

export interface StatePasswordForgotTFAVerify extends State, StateWithError {
  route: Routes.PASSWORD_FORGOT_TFA_APP | Routes.PASSWORD_FORGOT_TFA_EMAIL | Routes.PASSWORD_FORGOT_TFA_SMS;
  email?: string;
  phone?: string;
  verifyTFA: (code: string) => Promise<void>;
}

export interface StateLogout extends State {
  route: Routes.LOGOUT;
  logout: () => Promise<SuccessResponse>;
}

export const typeToSetupRoute = {
  [TFAType.TFATypeApp]: Routes.TFA_SETUP_APP,
  [TFAType.TFATypeEmail]: Routes.TFA_SETUP_EMAIL,
  [TFAType.TFATypeSMS]: Routes.TFA_SETUP_SMS,
};

export const typeToTFAVerifyRoute = {
  [TFAType.TFATypeApp]: Routes.TFA_VERIFY_APP,
  [TFAType.TFATypeEmail]: Routes.TFA_VERIFY_EMAIL,
  [TFAType.TFATypeSMS]: Routes.TFA_VERIFY_SMS,
};
export const typeToPasswordForgotTFAVerifyRoute = {
  [TFAType.TFATypeApp]: Routes.PASSWORD_FORGOT_TFA_APP,
  [TFAType.TFATypeEmail]: Routes.PASSWORD_FORGOT_TFA_EMAIL,
  [TFAType.TFATypeSMS]: Routes.PASSWORD_FORGOT_TFA_SMS,
};

// TODO exclude generalState
export type States =
  | State
  | StateTFASetupApp
  | StateTFASetupEmail
  | StateTFASetupSMS
  | StatePasswordReset
  | StatePasswordForgot
  | StatePasswordForgotSuccess
  | StateLoading
  | StateCallback
  | StateLogin
  | StateRegister
  | StateError;

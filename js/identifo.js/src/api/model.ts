/* eslint-disable camelcase */
export enum APIErrorCodes {
  PleaseEnableTFA = 'error.api.request.2fa.please_enable',
  NetworkError = 'error.network',
}

export enum TFAType {
  TFATypeApp = 'app',
  TFATypeSMS = 'sms',
  TFATypeEmail = 'email',
}
export interface ApiRequestError {
  error: {
    detailed_message?: string;
    id?: APIErrorCodes;
    message?: string;
    status?: number;
  };
}
export class ApiError extends Error {
  detailedMessage?: string;

  id?: APIErrorCodes;

  status?: number;

  constructor(error?: ApiRequestError['error']) {
    super(error?.message || 'Unknown API error');
    this.detailedMessage = error?.detailed_message;
    this.id = error?.id;
    this.status = error?.status;
  }
}
export interface LoginResponse {
  access_token?: string;
  refresh_token?: string;
  require_2fa: boolean;
  enabled_2fa: boolean;
  user: {
    active: boolean;
    email?: string;
    id: string;
    latest_login_time: number;
    num_of_logins: number;
    username?: string;
    tfa_info: { hotp_expired_at: string };
    phone?: string;
  };
  scopes?: string[];
  callbackUrl?: string;
}
export interface EnableTFAResponse {
  provisioning_uri?: string;
  provisioning_qr?: string;
  access_token?: string;
}
export interface AppSettingsResponse {
  anonymousResitrationAllowed: boolean;
  active: boolean;
  description: string;
  id: string;
  newUserDefaultRole: string;
  offline: boolean;
  registrationForbidden: boolean;
  tfaType: TFAType;
  federatedProviders: string[];
}

export interface User {
  id: string;
  username: string;
  email: string;
  phone: string;
  active: boolean;
  tfa_info: {
    is_enabled: boolean;
  };
  num_of_logins: number;
  latest_login_time: number;
  access_role: string;
  anonymous: boolean;
  federated_ids: string[];
}
export interface UpdateUser {
  new_email?: string;
  new_phone?: string;
}
export interface SuccessResponse {
  result: 'ok';
}

export type FederatedLoginProvider = 'apple' | 'google' | 'facebook';

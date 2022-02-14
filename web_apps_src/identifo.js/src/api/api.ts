import TokenService from '../tokenService';
import { IdentifoConfig } from '../types/types';
import {
  AppSettingsResponse,
  EnableTFAResponse,
  LoginResponse,
  SuccessResponse,
  UpdateUser,
  User,
  ApiRequestError,
  ApiError,
  APIErrorCodes,
  FederatedLoginProvider,
  TokenResponse,
  TFARequiredRespopnse,
} from './model';

const APP_ID_HEADER_KEY = 'X-Identifo-Clientid';
const AUTHORIZATION_HEADER_KEY = 'Authorization';

export class API {
  baseUrl: string;

  appId: string;

  defaultHeaders = {
    [APP_ID_HEADER_KEY]: '',
    Accept: 'application/json',
    'Content-Type': 'application/json',
  };

  catchNetworkErrorHandler = (e: TypeError): never => {
    if (
      e.message === 'Network Error' ||
      e.message === 'Failed to fetch' ||
      e.message === 'Preflight response is not successful' ||
      e.message.indexOf('is not allowed by Access-Control-Allow-Origin') > -1
    ) {
      // eslint-disable-next-line no-console
      console.error(e.message);
      throw new ApiError({
        id: APIErrorCodes.NetworkError,
        status: 0,
        message: 'Configuration error',
        detailed_message:
          'Please check Identifo URL and add "' +
          `${window.location.protocol}//${window.location.host}" ` +
          'to "REDIRECT URLS" in Identifo app settings.',
      });
    }
    throw e;
  };

  checkStatusCodeAndGetJSON = async (r: Response): Promise<any> => {
    if (!r.ok) {
      const error = (await r.json()) as ApiRequestError;
      throw new ApiError(error?.error);
    }
    return r.json();
  };

  constructor(private config: IdentifoConfig, private tokenService: TokenService) {
    // remove trailing slash if exist
    this.baseUrl = config.url.replace(/\/$/, '');
    this.defaultHeaders[APP_ID_HEADER_KEY] = config.appId;
    this.appId = config.appId;
  }

  get<T>(path: string, options?: RequestInit): Promise<T> {
    return this.send(path, { method: 'GET', ...options });
  }

  put<T>(path: string, data: unknown, options?: RequestInit): Promise<T> {
    return this.send(path, { method: 'PUT', body: JSON.stringify(data), ...options });
  }

  post<T>(path: string, data: unknown, options?: RequestInit): Promise<T> {
    return this.send(path, { method: 'POST', body: JSON.stringify(data), ...options });
  }

  send<T>(path: string, options?: RequestInit): Promise<T> {
    const init = { ...options };
    init.credentials = 'include';
    init.headers = {
      ...init.headers,
      ...this.defaultHeaders,
    };
    return fetch(`${this.baseUrl}${path}`, init)
      .catch(this.catchNetworkErrorHandler)
      .then(this.checkStatusCodeAndGetJSON)
      .then((value) => value as T);
  }

  async getUser(): Promise<User> {
    if (!this.tokenService.getToken()?.token) {
      throw new Error('No token in token service.');
    }
    return this.get<User>('/me', {
      headers: {
        [AUTHORIZATION_HEADER_KEY]: `Bearer ${this.tokenService.getToken()?.token}`,
      },
    });
  }

  async renewToken(): Promise<LoginResponse> {
    if (!this.tokenService.getToken('refresh')?.token) {
      throw new Error('No token in token service.');
    }
    return this.post<LoginResponse>(
      '/auth/token',
      { scopes: this.config.scopes },
      {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${this.tokenService.getToken('refresh')?.token}`,
        },
      },
    ).then((r) => this.storeToken(r));
  }

  async updateUser(user: UpdateUser): Promise<User> {
    if (!this.tokenService.getToken()?.token) {
      throw new Error('No token in token service.');
    }
    return this.put<User>('/me', user, {
      headers: {
        [AUTHORIZATION_HEADER_KEY]: `Bearer ${this.tokenService.getToken('access')?.token}`,
      },
    });
  }

  async login(email: string, password: string, deviceToken: string, scopes: string[]): Promise<LoginResponse> {
    const data = {
      email,
      password,
      device_token: deviceToken,
      scopes,
    };

    return this.post<LoginResponse>('/auth/login', data).then((r) => this.storeToken(r));
  }

  async requestPhoneCode(phone: string): Promise<SuccessResponse> {
    const data = {
      phone_number: phone,
    };

    return this.post<SuccessResponse>('/auth/request_phone_code', data);
  }

  async phoneLogin(phone: string, code: string, scopes: string[]): Promise<LoginResponse> {
    const data = {
      phone_number: phone,
      code,
      scopes,
    };

    return this.post<LoginResponse>('/auth/phone_login', data).then((r) => this.storeToken(r));
  }

  // After complete login on provider browser will be redirected to redirectUrl
  // callbackUrl will be stored in sesson and returned after successfull login complete
  async federatedLogin(
    provider: FederatedLoginProvider,
    scopes: string[],
    redirectUrl: string,
    callbackUrl?: string,
    opts: { width?: number; height?: number; popUp?: boolean } = { width: 600, height: 800, popUp: false },
  ): Promise<void> {
    const dataForm = document.createElement('form');
    dataForm.style.display = 'none';
    if (opts.popUp) {
      dataForm.target = 'TargetWindow'; // Make sure the window name is same as this value
    }
    dataForm.method = 'POST';
    const params = new URLSearchParams();
    params.set('appId', this.config.appId);
    params.set('provider', provider);
    params.set('scopes', scopes.join(','));
    params.set('redirectUrl', redirectUrl);
    if (callbackUrl) {
      params.set('callbackUrl', callbackUrl);
    }
    dataForm.action = `${this.baseUrl}/auth/federated?${params.toString()}`;

    document.body.appendChild(dataForm);

    if (opts.popUp) {
      const left = window.screenX + window.outerWidth / 2 - (opts.width || 600) / 2;
      const top = window.screenY + window.outerHeight / 2 - (opts.height || 800) / 2;
      const postWindow = window.open(
        '',
        'TargetWindow',
        `status=0,title=0,height=${opts.height},width=${opts.width},top=${top},left=${left},scrollbars=1`,
      );
      if (postWindow) {
        dataForm.submit();
      }
    } else {
      window.location.assign(`${this.baseUrl}/auth/federated?${params.toString()}`);
      // dataForm.submit();
    }
  }

  async federatedLoginComplete(params: URLSearchParams): Promise<LoginResponse> {
    return this.get<LoginResponse>(`/auth/federated/complete?${params.toString()}`).then((r) => this.storeToken(r));
  }

  async register(email: string, password: string, scopes: string[]): Promise<LoginResponse> {
    const data = {
      email,
      password,
      scopes,
    };

    return this.post<LoginResponse>('/auth/register', data).then((r) => this.storeToken(r));
  }

  async requestResetPassword(email: string, tfaCode?: string): Promise<SuccessResponse | TFARequiredRespopnse> {
    const data = {
      email,
      tfa_code: tfaCode,
    };

    return this.post<SuccessResponse | TFARequiredRespopnse>('/auth/request_reset_password', data);
  }

  async resetPassword(password: string): Promise<SuccessResponse> {
    if (!this.tokenService.getToken()?.token) {
      throw new Error('No token in token service.');
    }
    const data = {
      password,
    };

    return this.post<SuccessResponse>('/auth/reset_password', data, {
      headers: {
        [AUTHORIZATION_HEADER_KEY]: `Bearer ${this.tokenService.getToken()?.token}`,
      },
    });
  }

  async getAppSettings(callbackUrl: string): Promise<AppSettingsResponse> {
    return this.get<AppSettingsResponse>(`/auth/app_settings?${new URLSearchParams({ callbackUrl }).toString()}`);
  }

  async enableTFA(data: { phone?: string; email?: string }): Promise<EnableTFAResponse> {
    if (!this.tokenService.getToken()?.token) {
      throw new Error('No token in token service.');
    }
    return this.put<EnableTFAResponse>('/auth/tfa/enable', data, {
      headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${this.tokenService.getToken()?.token}` },
    }).then((r) => this.storeToken(r));
  }

  async verifyTFA(code: string, scopes: string[]): Promise<LoginResponse> {
    if (!this.tokenService.getToken()?.token) {
      throw new Error('No token in token service.');
    }
    return this.post<LoginResponse>(
      '/auth/tfa/login',
      { tfa_code: code, scopes },
      { headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${this.tokenService.getToken()?.token}` } },
    ).then((r) => this.storeToken(r));
  }

  async resendTFA(): Promise<LoginResponse> {
    if (!this.tokenService.getToken()?.token) {
      throw new Error('No token in token service.');
    }
    return this.post<LoginResponse>('/auth/tfa/resend', null, {
      headers: { [AUTHORIZATION_HEADER_KEY]: `BEARER ${this.tokenService.getToken()?.token}` },
    }).then((r) => this.storeToken(r));
  }

  async logout(): Promise<SuccessResponse> {
    if (!this.tokenService.getToken()?.token) {
      throw new Error('No token in token service.');
    }
    return this.post<SuccessResponse>(
      '/me/logout',
      {
        refresh_token: this.tokenService.getToken('refresh')?.token,
      },
      {
        headers: {
          [AUTHORIZATION_HEADER_KEY]: `Bearer ${this.tokenService.getToken()?.token}`,
        },
      },
    ).then((r) => {
      this.tokenService.removeToken();
      this.tokenService.removeToken('refresh');
      return r;
    });
  }

  storeToken<T extends TokenResponse>(response: T): T {
    if (response.access_token) {
      this.tokenService.saveToken(response.access_token, 'access');
    }
    if (response.refresh_token) {
      this.tokenService.saveToken(response.refresh_token, 'refresh');
    }
    return response;
  }
}

import TokenService from '../tokenService';
import { IdentifoConfig } from '../types/types';
import { AppSettingsResponse, EnableTFAResponse, LoginResponse, SuccessResponse, UpdateUser, User, FederatedLoginProvider, TokenResponse, TFARequiredRespopnse } from './model';
export declare class API {
    private config;
    private tokenService;
    baseUrl: string;
    appId: string;
    defaultHeaders: {
        "X-Identifo-Clientid": string;
        Accept: string;
        'Content-Type': string;
    };
    catchNetworkErrorHandler: (e: TypeError) => never;
    checkStatusCodeAndGetJSON: (r: Response) => Promise<any>;
    constructor(config: IdentifoConfig, tokenService: TokenService);
    get<T>(path: string, options?: RequestInit): Promise<T>;
    put<T>(path: string, data: unknown, options?: RequestInit): Promise<T>;
    post<T>(path: string, data: unknown, options?: RequestInit): Promise<T>;
    send<T>(path: string, options?: RequestInit): Promise<T>;
    getUser(): Promise<User>;
    renewToken(): Promise<LoginResponse>;
    updateUser(user: UpdateUser): Promise<User>;
    login(email: string, password: string, deviceToken: string, scopes: string[]): Promise<LoginResponse>;
    federatedLogin(provider: FederatedLoginProvider, scopes: string[], redirectUrl: string, callbackUrl?: string, opts?: {
        width?: number;
        height?: number;
        popUp?: boolean;
    }): Promise<void>;
    federatedLoginComplete(params: URLSearchParams): Promise<LoginResponse>;
    register(email: string, password: string, scopes: string[]): Promise<LoginResponse>;
    requestResetPassword(email: string, tfaCode?: string): Promise<SuccessResponse | TFARequiredRespopnse>;
    resetPassword(password: string): Promise<SuccessResponse>;
    getAppSettings(callbackUrl: string): Promise<AppSettingsResponse>;
    enableTFA(): Promise<EnableTFAResponse>;
    verifyTFA(code: string, scopes: string[]): Promise<LoginResponse>;
    logout(): Promise<SuccessResponse>;
    storeToken<T extends TokenResponse>(response: T): T;
}

import { JWSHeaderParameters } from 'jose/webcrypto/types';

declare type TokenType = 'access' | 'refresh';
interface TokenManager {
    isAccessible: boolean;
    preffix: string;
    storageType: string;
    access: string;
    refresh: string;
    saveToken: (token: string, tokenType: TokenType) => boolean;
    getToken: (tokenType: TokenType) => string;
    deleteToken: (tokenType: TokenType) => void;
}
declare type IdentifoConfig = {
    issuer?: string;
    appId: string;
    url: string;
    scopes?: string[];
    redirectUri?: string;
    postLogoutRedirectUri?: string;
    tokenManager?: TokenManager;
    autoRenew?: boolean;
};
declare type UrlBuilderInit = {
    createSignupUrl: () => string;
    createSigninUrl: () => string;
    createLogoutUrl: () => string;
    createRenewSessionUrl: () => string;
};
interface JWTPayload {
    /**
     * JWT Issuer - [RFC7519#section-4.1.1](https://tools.ietf.org/html/rfc7519#section-4.1.1).
     */
    iss?: string;
    /**
     * JWT Subject - [RFC7519#section-4.1.2](https://tools.ietf.org/html/rfc7519#section-4.1.2).
     */
    sub?: string;
    /**
     * JWT Audience [RFC7519#section-4.1.3](https://tools.ietf.org/html/rfc7519#section-4.1.3).
     */
    aud?: string[];
    /**
     * JWT ID - [RFC7519#section-4.1.7](https://tools.ietf.org/html/rfc7519#section-4.1.7).
     */
    jti?: string;
    /**
     * JWT Not Before - [RFC7519#section-4.1.5](https://tools.ietf.org/html/rfc7519#section-4.1.5).
     */
    nbf?: number;
    /**
     * JWT Expiration Time - [RFC7519#section-4.1.4](https://tools.ietf.org/html/rfc7519#section-4.1.4).
     */
    exp?: number;
    /**
     * JWT Issued At - [RFC7519#section-4.1.6](https://tools.ietf.org/html/rfc7519#section-4.1.6).
     */
    iat?: number;
    /**
     * Any other JWT Claim Set member.
     */
    [propName: string]: unknown;
}
declare type ClientToken = {
    token: string;
    payload: JWTPayload;
    header?: JWSHeaderParameters;
};

declare class TokenService {
    private tokenManager;
    constructor(tokenManager?: TokenManager);
    handleVerification(token: string, audience: string, issuer?: string): Promise<boolean>;
    validateToken(token: string, audience: string, issuer?: string): Promise<boolean>;
    parseJWT(token: string): JWTPayload;
    isJWTExpired(token: JWTPayload): boolean;
    isAuthenticated(audience: string, issuer?: string): Promise<boolean>;
    saveToken(token: string, type?: TokenType): boolean;
    removeToken(type?: TokenType): void;
    getToken(type?: TokenType): ClientToken | null;
}

declare enum APIErrorCodes {
    PleaseEnableTFA = "error.api.request.2fa.please_enable",
    NetworkError = "error.network"
}
declare enum TFAType {
    TFATypeApp = "app",
    TFATypeSMS = "sms",
    TFATypeEmail = "email"
}
interface ApiRequestError {
    error: {
        detailed_message?: string;
        id?: APIErrorCodes;
        message?: string;
        status?: number;
    };
}
declare class ApiError extends Error {
    detailedMessage?: string;
    id?: APIErrorCodes;
    status?: number;
    constructor(error?: ApiRequestError['error']);
}
interface LoginResponse {
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
        tfa_info: {
            hotp_expired_at: string;
        };
        phone?: string;
    };
    scopes?: string[];
    callbackUrl?: string;
}
interface EnableTFAResponse {
    provisioning_uri?: string;
    provisioning_qr?: string;
    access_token?: string;
}
interface AppSettingsResponse {
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
interface User {
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
interface UpdateUser {
    new_email?: string;
    new_phone?: string;
}
interface SuccessResponse {
    result: 'ok';
}
declare type FederatedLoginProvider = 'apple' | 'google' | 'facebook';

declare class Api {
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
    requestResetPassword(email: string): Promise<SuccessResponse>;
    resetPassword(password: string): Promise<SuccessResponse>;
    getAppSettings(): Promise<AppSettingsResponse>;
    enableTFA(): Promise<EnableTFAResponse>;
    verifyTFA(code: string, scopes: string[]): Promise<LoginResponse>;
    logout(): Promise<SuccessResponse>;
    storeToken(response: LoginResponse): LoginResponse;
}

declare class IdentifoAuth {
    api: Api;
    tokenService: TokenService;
    config: IdentifoConfig;
    urlBuilder: UrlBuilderInit;
    private token;
    isAuth: boolean;
    constructor(config: IdentifoConfig);
    private handleToken;
    private resetAuthValues;
    signup(): void;
    signin(): void;
    logout(): void;
    handleAuthentication(): Promise<boolean>;
    private getTokenFromUrl;
    getToken(): Promise<ClientToken | null>;
    renewSession(): Promise<string>;
    private renewSessionWithToken;
}

declare class CookieStorage {
    isAccessible: boolean;
    saveToken(): boolean;
    getToken(): string;
    deleteToken(): void;
}

declare class StorageManager implements TokenManager {
    preffix: string;
    storageType: 'localStorage' | 'sessionStorage';
    access: string;
    refresh: string;
    isAccessible: boolean;
    constructor(storageType: 'localStorage' | 'sessionStorage', accessKey?: string, refreshKey?: string);
    saveToken(token: string, tokenType: TokenType): boolean;
    getToken(tokenType: TokenType): string;
    deleteToken(tokenType: TokenType): void;
}

declare class LocalStorage extends StorageManager {
    constructor(accessKey?: string, refreshKey?: string);
}

declare class SessionStorage extends StorageManager {
    constructor(accessKey?: string, refreshKey?: string);
}

export { APIErrorCodes, ApiError, ApiRequestError, AppSettingsResponse, CookieStorage as CookieStorageManager, EnableTFAResponse, FederatedLoginProvider, IdentifoAuth, LocalStorage as LocalStorageManager, LoginResponse, SessionStorage as SessionStorageManager, SuccessResponse, TFAType, UpdateUser, User };

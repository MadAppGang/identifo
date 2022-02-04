import { BehaviorSubject } from 'rxjs';

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
declare type UrlFlows = 'signin' | 'signup' | 'logout' | 'renew' | 'default';
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
};

declare class TokenService {
    isAuth: boolean;
    private tokenManager;
    constructor(tokenManager?: TokenManager);
    handleVerification(token: string, audience: string, issuer?: string): Promise<boolean>;
    validateToken(token: string, audience: string, issuer?: string): Promise<boolean>;
    parseJWT(token: string): JWTPayload;
    isJWTExpired(token: JWTPayload): boolean;
    saveToken(token: string, type?: TokenType): boolean;
    removeToken(type?: TokenType): void;
    getToken(type?: TokenType): ClientToken | null;
}

declare enum APIErrorCodes {
    PleaseEnableTFA = "error.api.request.2fa.please_enable",
    InvalidCallbackURL = "error.api.request.callbackurl.invalid",
    NetworkError = "error.network"
}
declare enum TFAType {
    TFATypeApp = "app",
    TFATypeSMS = "sms",
    TFATypeEmail = "email"
}
declare enum TFAStatus {
    DISABLED = "disabled",
    OPTIONAL = "optional",
    MANDATORY = "mandatory"
}
declare type FederatedLoginProvider = 'apple' | 'google' | 'facebook';
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
interface TokenResponse {
    access_token?: string;
    refresh_token?: string;
}
interface AppSettingsResponse {
    anonymousResitrationAllowed: boolean;
    active: boolean;
    description: string;
    id: string;
    newUserDefaultRole: string;
    offline: boolean;
    registrationForbidden: boolean;
    tfaType: TFAType[] | TFAType;
    tfaResendTimeout: number;
    tfaStatus: TFAStatus;
    federatedProviders: FederatedLoginProvider[];
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
interface TFARequiredRespopnse {
    result: 'tfa-required';
}

declare class API {
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
    enableTFA(data: {
        phone?: string;
        email?: string;
    }): Promise<EnableTFAResponse>;
    verifyTFA(code: string, scopes: string[]): Promise<LoginResponse>;
    resendTFA(): Promise<LoginResponse>;
    logout(): Promise<SuccessResponse>;
    storeToken<T extends TokenResponse>(response: T): T;
}

declare class IdentifoAuth {
    api: API;
    tokenService: TokenService;
    config: IdentifoConfig;
    urlBuilder: UrlBuilderInit;
    private token;
    get isAuth(): boolean;
    constructor(config?: IdentifoConfig);
    configure(config: IdentifoConfig): void;
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

declare enum Routes {
    'LOGIN' = "login",
    'REGISTER' = "register",
    'TFA_VERIFY_SMS' = "tfa/verify/sms",
    'TFA_VERIFY_EMAIL' = "tfa/verify/email",
    'TFA_VERIFY_APP' = "tfa/verify/app",
    'TFA_VERIFY_SELECT' = "tfa/verify/select",
    'TFA_SETUP_SMS' = "tfa/setup/sms",
    'TFA_SETUP_EMAIL' = "tfa/setup/email",
    'TFA_SETUP_APP' = "tfa/setup/app",
    'TFA_SETUP_SELECT' = "tfa/setup/select",
    'PASSWORD_RESET' = "password/reset",
    'PASSWORD_FORGOT' = "password/forgot",
    'PASSWORD_FORGOT_TFA_SMS' = "password/forgot/tfa/sms",
    'PASSWORD_FORGOT_TFA_EMAIL' = "password/forgot/tfa/email",
    'PASSWORD_FORGOT_TFA_APP' = "password/forgot/tfa/app",
    'PASSWORD_FORGOT_TFA_SELECT' = "password/forgot/tfa/select",
    'CALLBACK' = "callback",
    'OTP_LOGIN' = "otp/login",
    'ERROR' = "error",
    'PASSWORD_FORGOT_SUCCESS' = "password/forgot/success",
    'LOGOUT' = "logout",
    'LOADING' = "loading"
}
declare type TFASetupRoutes = Routes.TFA_SETUP_SELECT | Routes.TFA_SETUP_SMS | Routes.TFA_SETUP_EMAIL | Routes.TFA_SETUP_APP;
declare type TFALoginVerifyRoutes = Routes.TFA_VERIFY_SELECT | Routes.TFA_VERIFY_SMS | Routes.TFA_VERIFY_EMAIL | Routes.TFA_VERIFY_APP;
declare type TFAResetVerifyRoutes = Routes.PASSWORD_FORGOT_TFA_SELECT | Routes.PASSWORD_FORGOT_TFA_SMS | Routes.PASSWORD_FORGOT_TFA_EMAIL | Routes.PASSWORD_FORGOT_TFA_APP;
interface State {
    route: Routes;
}
interface StateWithError {
    error: ApiError;
}
interface StateLogin extends State, StateWithError {
    route: Routes.LOGIN;
    registrationForbidden: boolean;
    federatedProviders: FederatedLoginProvider[];
    signup: () => Promise<void>;
    signin: (email: string, password: string, remember?: boolean) => Promise<void>;
    socialLogin: (provider: FederatedLoginProvider) => Promise<void>;
    passwordForgot: () => Promise<void>;
}
interface StateRegister extends State, StateWithError {
    route: Routes.REGISTER;
    signup: (email: string, password: string) => Promise<void>;
    goback: () => Promise<void>;
}
interface StatePasswordForgot extends State, StateWithError {
    route: Routes.PASSWORD_FORGOT;
    restorePassword: (email: string) => Promise<void>;
    goback: () => Promise<void>;
}
interface StatePasswordForgotSuccess extends State {
    route: Routes.PASSWORD_FORGOT_SUCCESS;
    goback: () => Promise<void>;
}
interface StateError extends State, StateWithError {
    route: Routes.ERROR;
}
interface StateCallback extends State {
    route: Routes.CALLBACK;
    callbackUrl: string;
    result: LoginResponse;
}
interface StatePasswordReset extends State, StateWithError {
    route: Routes.PASSWORD_RESET;
    setNewPassword: (password: string) => Promise<void>;
}
interface StateLoading extends State {
    route: Routes.LOADING;
}
interface StateOTPLogin extends State {
    route: Routes.OTP_LOGIN;
    registrationForbidden: boolean;
    federatedProviders: FederatedLoginProvider[];
    signup: () => Promise<void>;
    signin: (phone: string) => Promise<void>;
    socialLogin: (provider: FederatedLoginProvider) => Promise<void>;
}
interface StateTFASetup extends State, StateWithError {
}
interface StateTFASetupApp extends StateTFASetup {
    route: Routes.TFA_SETUP_APP;
    provisioningURI: string;
    provisioningQR: string;
    setupTFA: () => Promise<void>;
}
interface StateTFASetupEmail extends StateTFASetup {
    route: Routes.TFA_SETUP_EMAIL;
    email: string;
    setupTFA: (email: string) => Promise<void>;
}
interface StateTFASetupSMS extends StateTFASetup {
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
interface StateTFASetupSelect extends StateTFASelect {
    route: Routes.TFA_SETUP_SELECT;
    tfaStatus: TFAStatus;
    setupNextTime: () => Promise<void>;
}
interface StateTFAVerifySelect extends StateTFASelect {
    route: Routes.TFA_VERIFY_SELECT;
}
interface StatePasswordForgotTFASelect extends StateTFASelect {
    route: Routes.PASSWORD_FORGOT_TFA_SELECT;
}
interface StateTFAVerifyApp extends State, StateWithError {
    route: Routes.TFA_VERIFY_APP;
    email?: string;
    phone?: string;
    verifyTFA: (code: string) => Promise<void>;
}
interface StateTFAVerifyEmailSms extends State, StateWithError {
    route: Routes.TFA_VERIFY_EMAIL | Routes.TFA_VERIFY_SMS;
    email?: string;
    phone?: string;
    resendTimeout: number;
    verifyTFA: (code: string) => Promise<void>;
    resendTFA: () => Promise<void>;
}
interface StatePasswordForgotTFAVerify extends State, StateWithError {
    route: Routes.PASSWORD_FORGOT_TFA_APP | Routes.PASSWORD_FORGOT_TFA_EMAIL | Routes.PASSWORD_FORGOT_TFA_SMS;
    email?: string;
    phone?: string;
    verifyTFA: (code: string) => Promise<void>;
}
interface StateLogout extends State {
    route: Routes.LOGOUT;
    logout: () => Promise<SuccessResponse>;
}
declare const typeToSetupRoute: {
    app: Routes;
    email: Routes;
    sms: Routes;
};
declare const typeToTFAVerifyRoute: {
    app: Routes;
    email: Routes;
    sms: Routes;
};
declare const typeToPasswordForgotTFAVerifyRoute: {
    app: Routes;
    email: Routes;
    sms: Routes;
};
declare type States = State | StateTFASetupApp | StateTFASetupEmail | StateTFASetupSMS | StatePasswordReset | StatePasswordForgot | StatePasswordForgotSuccess | StateLoading | StateCallback | StateLogin | StateRegister | StateError;

declare class CDK {
    auth: IdentifoAuth;
    settings: AppSettingsResponse;
    lastError: ApiError;
    callbackUrl?: string;
    postLogoutRedirectUri?: string;
    scopes: Set<string>;
    state: BehaviorSubject<States>;
    constructor();
    configure(authConfig: IdentifoConfig, callbackUrl: string): Promise<void>;
    login(): void;
    register(): void;
    forgotPassword(): void;
    forgotPasswordSuccess(): void;
    passwordReset(): void;
    callback(result: LoginResponse): void;
    validateEmail(email: string): boolean;
    tfaSetup(loginResponse: LoginResponse, type: TFAType): Promise<void>;
    tfaVerify(loginResponse: LoginResponse, type: TFAType): Promise<void>;
    passwordForgotTFAVerify(email: string, type: TFAType): Promise<void>;
    logout(): Promise<void>;
    private processError;
    private redirectTfaSetup;
    private tfaSetupSelect;
    private redirectTfaVerify;
    private redirectTfaForgot;
    private afterLoginRedirect;
    private loginCatchRedirect;
}

export { APIErrorCodes, ApiError, ApiRequestError, AppSettingsResponse, CDK, ClientToken, CookieStorage as CookieStorageManager, EnableTFAResponse, FederatedLoginProvider, IdentifoAuth, IdentifoConfig, JWTPayload, LocalStorage as LocalStorageManager, LoginResponse, Routes, SessionStorage as SessionStorageManager, State, StateCallback, StateError, StateLoading, StateLogin, StateLogout, StateOTPLogin, StatePasswordForgot, StatePasswordForgotSuccess, StatePasswordForgotTFASelect, StatePasswordForgotTFAVerify, StatePasswordReset, StateRegister, StateTFASetupApp, StateTFASetupEmail, StateTFASetupSMS, StateTFASetupSelect, StateTFAVerifyApp, StateTFAVerifyEmailSms, StateTFAVerifySelect, StateWithError, States, SuccessResponse, TFALoginVerifyRoutes, TFARequiredRespopnse, TFAResetVerifyRoutes, TFASetupRoutes, TFAStatus, TFAType, TokenManager, TokenResponse, TokenType, UpdateUser, UrlBuilderInit, UrlFlows, User, typeToPasswordForgotTFAVerifyRoute, typeToSetupRoute, typeToTFAVerifyRoute };

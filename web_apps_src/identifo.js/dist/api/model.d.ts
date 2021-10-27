export declare enum APIErrorCodes {
    PleaseEnableTFA = "error.api.request.2fa.please_enable",
    InvalidCallbackURL = "error.api.request.callbackurl.invalid",
    NetworkError = "error.network"
}
export declare enum TFAType {
    TFATypeApp = "app",
    TFATypeSMS = "sms",
    TFATypeEmail = "email"
}
export declare enum TFAStatus {
    DISABLED = "disabled",
    OPTIONAL = "optional",
    MANDATORY = "mandatory"
}
export interface ApiRequestError {
    error: {
        detailed_message?: string;
        id?: APIErrorCodes;
        message?: string;
        status?: number;
    };
}
export declare class ApiError extends Error {
    detailedMessage?: string;
    id?: APIErrorCodes;
    status?: number;
    constructor(error?: ApiRequestError['error']);
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
        tfa_info: {
            hotp_expired_at: string;
        };
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
export interface TokenResponse {
    access_token?: string;
    refresh_token?: string;
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
    tfaStatus: TFAStatus;
    federatedProviders: FederatedLoginProvider[];
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
export interface TFARequiredRespopnse {
    result: 'tfa-required';
}
export declare type FederatedLoginProvider = 'apple' | 'google' | 'facebook';

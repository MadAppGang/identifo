import { ApiError, FederatedLoginProvider } from '../api/model';
export declare enum Routes {
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
export declare type TFASetupRoutes = Routes.TFA_SETUP_SELECT | 'TFA_SETUP_SMS' | 'TFA_SETUP_EMAIL' | 'TFA_SETUP_APP';
export declare type TFALoginVerifyRoutes = 'TFA_VERIFY_SELECT' | 'TFA_VERIFY_SMS' | 'TFA_VERIFY_EMAIL' | 'TFA_VERIFY_APP';
export declare type TFAResetVerifyRoutes = 'PASSWORD_FORGOT_TFA_SELECT' | 'PASSWORD_FORGOT_TFA_SMS' | 'PASSWORD_FORGOT_TFA_EMAIL' | 'PASSWORD_FORGOT_TFA_APP';
export interface State {
    route: Routes;
}
export interface StateLogin extends State {
    route: Routes.LOGIN;
    registrationForbidden: boolean;
    lastError: ApiError;
    federatedProviders: FederatedLoginProvider[];
    signup: () => void;
    signin: () => void;
    socialLogin: (provider: FederatedLoginProvider) => void;
    passwordForgot: () => void;
}
export interface StateRegister extends State {
    route: Routes.REGISTER;
    registrationForbidden: boolean;
    lastError: ApiError;
    federatedProviders: FederatedLoginProvider[];
    signup: () => void;
    signin: () => void;
    socialLogin: (provider: FederatedLoginProvider) => void;
    passwordForgot: () => void;
}
export declare type States = StateLogin | StateRegister;

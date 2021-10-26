import { API } from './api/api';
import TokenService from './tokenService';
import { ClientToken, IdentifoConfig, UrlBuilderInit } from './types/types';
declare class IdentifoAuth {
    api: API;
    tokenService: TokenService;
    config: IdentifoConfig;
    urlBuilder: UrlBuilderInit;
    private token;
    isAuth: boolean;
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
export default IdentifoAuth;

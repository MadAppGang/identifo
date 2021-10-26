import { IdentifoConfig, UrlFlows } from './types/types';
export declare class UrlBuilder {
    private config;
    constructor(config: IdentifoConfig);
    getUrl(flow: UrlFlows): string;
    createSignupUrl(): string;
    createSigninUrl(): string;
    createLogoutUrl(): string;
    createRenewSessionUrl(): string;
}

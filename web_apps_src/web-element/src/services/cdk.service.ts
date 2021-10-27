import { CDK, IdentifoConfig } from '@identifo/identifo-auth-js';

class IdentifoCDKController {
  public cdk: CDK;
  public debug: boolean;
  constructor() {
    this.cdk = new CDK();
  }
  public configure(config: IdentifoConfig, callbackUrl: string, scopes: string[]): Promise<void> {
    return this.cdk.configure(config, callbackUrl, scopes);
  }
}
export const CDKService = new IdentifoCDKController();

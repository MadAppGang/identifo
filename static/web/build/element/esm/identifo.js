import { p as promiseResolve, b as bootstrapLazy } from './index-0998ee25.js';

/*
 Stencil Client Patch Browser v2.6.0 | MIT Licensed | https://stenciljs.com
 */
const patchBrowser = () => {
    const importMeta = import.meta.url;
    const opts = {};
    if (importMeta !== '') {
        opts.resourcesUrl = new URL('.', importMeta).href;
    }
    return promiseResolve(opts);
};

patchBrowser().then(options => {
  return bootstrapLazy([["identifo-form",[[1,"identifo-form",{"route":[1537],"token":[1],"appId":[513,"app-id"],"url":[513],"theme":[1],"scopes":[1],"auth":[32],"username":[32],"password":[32],"phone":[32],"email":[32],"registrationForbidden":[32],"tfaCode":[32],"tfaType":[32],"tfaMandatory":[32],"provisioningURI":[32],"provisioningQR":[32],"success":[32],"lastError":[32],"lastResponse":[32]}]]]], options);
});

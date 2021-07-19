'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

const index = require('./index-efa5969e.js');

/*
 Stencil Client Patch Esm v2.6.0 | MIT Licensed | https://stenciljs.com
 */
const patchEsm = () => {
    return index.promiseResolve();
};

const defineCustomElements = (win, options) => {
  if (typeof window === 'undefined') return Promise.resolve();
  return patchEsm().then(() => {
  return index.bootstrapLazy([["identifo-form.cjs",[[1,"identifo-form",{"route":[1537],"token":[1],"appId":[513,"app-id"],"url":[513],"theme":[1],"scopes":[1],"auth":[32],"username":[32],"password":[32],"phone":[32],"email":[32],"registrationForbidden":[32],"tfaCode":[32],"tfaType":[32],"tfaMandatory":[32],"provisioningURI":[32],"provisioningQR":[32],"success":[32],"lastError":[32],"lastResponse":[32]}]]]], options);
  });
};

exports.defineCustomElements = defineCustomElements;

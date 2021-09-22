/// <reference types="cypress" />
// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)

/**
 * @type {Cypress.PluginConfig}
 */
// eslint-disable-next-line no-unused-vars
// const _ = require('lodash');
// const del = require('del');
const { renameSync } = require('fs');
module.exports = (on, config) => {
  on('after:screenshot', ({ path }) => {
    renameSync(path, path.replace(/ \(\d*\)/i, ''));
  });
  // on('after:spec', (spec, results) => {
  //   if (results && results.video) {
  //     // Do we have failures for any retry attempts?
  //     const failures = _.some(results.tests, test => {
  //       return _.some(test.attempts, { state: 'failed' });
  //     });
  //     if (!failures) {
  //       // delete the video if the spec passed and no tests retried
  //       return del(results.video);
  //     }
  //   }
  // });
  on('before:browser:launch', (browser = {}, launchOptions) => {
    if (browser.name === 'chrome' || browser.name === 'edge') {
      launchOptions.args.push('--disable-features=SameSiteByDefaultCookies'); // bypass 401 unauthorised access on chromium-based browsers
      return launchOptions;
    }
  });
  return config;
};

const { defineConfig } = require('cypress')

module.exports = defineConfig({
  projectId: '9g4zv8',
  video: false,
  viewportHeight: 800,
  viewportWidth: 600,
  serverUrl: 'http://localhost:8081',
  chromeWebSecurity: false,
  e2e: {
    // We've imported your old cypress plugins here.
    // You may want to clean this up later by importing these.
    setupNodeEvents(on, config) {
      return require('./cypress/plugins/index.js')(on, config)
    },
    baseUrl: 'http://localhost:8081/web',
    specPattern: 'cypress/e2e/**/*.{js,jsx,ts,tsx}',
  },
})

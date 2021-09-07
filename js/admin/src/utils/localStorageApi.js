const SETTINGS_KEY = 'iap_settings';

export const localStorageApi = {
  settings: {
    set: (payload) => {
      window.localStorage.setItem(SETTINGS_KEY, JSON.stringify(payload));
    },
    get: () => {
      try {
        const settings = JSON.parse(window.localStorage.getItem(SETTINGS_KEY));
        return settings;
      } catch (error) {
        return null;
      }
    },
  },
};

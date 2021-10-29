export const verificationStatuses = {
  required: 0,
  loading: 1,
  success: 2,
  fail: 3,
};

export const notificationStatuses = {
  idle: 0,
  success: 1,
  error: 2,
  changed: 3,
};

export const localStorageKeys = {
  markdown: 'iap-markdown',
};

export const tabGroups = {
  server_group: 'server_group',
  account_group: 'account_group',
  storages_group: 'storages_group',
  external_services_group: 'external_services_group',
  apple_integration_group: 'apple_integration_group',
  edit_app_group: 'edit_app_group',
};

export const storageTypes = {
  local: 'local',
  s3: 's3',
};

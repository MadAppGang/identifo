import { CONNECTION_FAILED, CONNECTION_TEST_REQUIRED } from '~/modules/database/connectionReducer';


export const settingsActionsEnum = {
  save: 'save',
  verify: 'verify',
  close: 'close',
};

export const settingsConfig = {
  [CONNECTION_FAILED]: {
    content: 'The server was unable to connect with these settings',
    buttons: [
      { label: 'Save', data: settingsActionsEnum.save },
      { label: 'Don`t Save', data: settingsActionsEnum.close },
    ],
  },
  [CONNECTION_TEST_REQUIRED]: {
    content: 'It would be a good idea to test the connection before saving.',
    buttons: [
      { label: 'Verify', data: settingsActionsEnum.verify },
      { label: 'Save without verification', data: settingsActionsEnum.save },
    ],
  },
};

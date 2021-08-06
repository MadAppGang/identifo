import React from 'react';
import PropTypes from 'prop-types';
import { DatagridRow } from '~/components/shared/Datagrid';
import MobileIcon from '~/components/icons/MobileIcon';
import WebIcon from '~/components/icons/WebIcon';

const getIconForType = (type = '') => {
  const mobileTypes = ['ios', 'android'];

  if (mobileTypes.includes(type.toLowerCase())) {
    return MobileIcon;
  }

  return WebIcon;
};

const ApplicationRow = ({ data, config }) => {
  const Icon = getIconForType(data.type);

  return (
    <DatagridRow key={data.id} className="iap-application-row">
      <div style={{ width: config.icon.width }}>
        <div className="iap-application-row__icon-wrapper">
          <Icon className="iap-application-row__icon" />
        </div>
      </div>
      <div style={{ width: config.type.width }}>
        <p className="iap-application-row__type">{data.type}</p>
        <p className="iap-application-row__name">{data.name}</p>
      </div>
      <div style={{ width: config.clientId.width }}>
        <p className="iap-application-row__clientid">{data.id}</p>
      </div>
      <div style={{ width: config.tfaStatus.width }}>
        <p className="iap-application-row__tfa-status">{data.tfa_status || 'disabled'}</p>
      </div>
      <div style={{ width: config.settings.width }} />
    </DatagridRow>
  );
};

ApplicationRow.propTypes = {
  data: PropTypes.shape({
    id: PropTypes.string,
    clientId: PropTypes.string,
  }).isRequired,
  config: PropTypes.shape({
    icon: PropTypes.object,
    type: PropTypes.object,
    clientId: PropTypes.object,
    settings: PropTypes.object,
  }).isRequired,
};

export default ApplicationRow;

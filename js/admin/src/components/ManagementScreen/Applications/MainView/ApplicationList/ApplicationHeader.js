import React from 'react';
import PropTypes from 'prop-types';
import { DatagridHeader } from '~/components/shared/Datagrid';

const ApplicationHeader = ({ config }) => (
  <DatagridHeader>
    {Object.keys(config).map(key => (
      <div key={key} style={{ width: config[key].width }}>
        {config[key].title}
      </div>
    ))}
  </DatagridHeader>
);

ApplicationHeader.propTypes = {
  config: PropTypes.shape().isRequired,
};

export default ApplicationHeader;

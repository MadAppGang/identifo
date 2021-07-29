import React from 'react';
import PropTypes from 'prop-types';
import DatagridRow from './Row';

const DatagridHeader = ({ children }) => (
  <header className="iap-datagrid-header">
    <DatagridRow header>
      {children}
    </DatagridRow>
  </header>
);

DatagridHeader.propTypes = {
  children: PropTypes.node.isRequired,
};

export default DatagridHeader;

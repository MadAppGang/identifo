import React from 'react';
import PropTypes from 'prop-types';
import NotFoundIcon from '~/components/icons/NotFoundIcon';

const DatagridNotFound = ({ text }) => (
  <div className="iap-datagrid__not-found">
    <NotFoundIcon className="iap-datagrid__not-found-icon" />
    <p className="iap-datagrid__not-found-text">
      {text}
    </p>
  </div>
);

DatagridNotFound.propTypes = {
  text: PropTypes.string,
};

DatagridNotFound.defaultProps = {
  text: 'Not Found',
};

export default DatagridNotFound;

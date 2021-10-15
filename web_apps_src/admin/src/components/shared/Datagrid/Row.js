import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

const DatagridRow = ({ children, header, className }) => (
  <div
    className={classnames({
      'iap-datagrid-row': true,
      'iap-datagrid-row--header': header,
      [className]: !!className,
    })}
  >
    {children}
  </div>
);

DatagridRow.propTypes = {
  children: PropTypes.node.isRequired,
  header: PropTypes.bool,
  className: PropTypes.string,
};

DatagridRow.defaultProps = {
  header: false,
  className: null,
};

export default DatagridRow;

import React from 'react';
import PropTypes from 'prop-types';
import './Container.css';

const Container = ({ children }) => (
  <div className="iap-global-container">
    {children}
  </div>
);

Container.propTypes = {
  children: PropTypes.node.isRequired,
};

export default Container;

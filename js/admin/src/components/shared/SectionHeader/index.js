import React from 'react';
import PropTypes from 'prop-types';
import './index.css';

const SectionHeader = props => (
  <header>
    <span className="iap-section__title">
      {props.title}
    </span>
    <p className="iap-section__description">
      {props.description}
    </p>
  </header>
);

SectionHeader.propTypes = {
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
};

export default SectionHeader;

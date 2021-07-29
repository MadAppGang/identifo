import React, { useRef, useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import Tab from './Tab';

const getElementWidth = el => el ? el.offsetWidth : 0;

const Tabs = ({ activeTabIndex, onChange, children }) => {
  const tabContent = [...children].pop();
  const tabs = children.slice(0, -1);

  if (tabContent.type === Tab) {
    throw Error('Last child should be tab content');
  }

  const tabRefsRef = useRef([]);
  const [tabWidths, setTabWidths] = useState([]);

  useEffect(() => {
    if (!tabRefsRef.current.length) {
      return;
    }

    setTabWidths(tabRefsRef.current.map(ref => getElementWidth(ref.current)));
  }, [tabRefsRef.current]);

  return (
    <>
      <div className="iap-tabs-tablist">
        {React.Children.map(tabs, (child, index) => {
          const ref = React.createRef();

          tabRefsRef.current.push(ref);

          return React.cloneElement(child, {
            ref,
            isActive: activeTabIndex === index,
            key: child.props.title,
            onClick: () => onChange(index),
          });
        })}
        <div
          className="tab-underline"
          style={{
            width: tabWidths[activeTabIndex] || 0,
            left: tabWidths.slice(0, activeTabIndex).reduce((a, b) => a + b, 0),
          }}
        />
      </div>
      {tabContent}
    </>
  );
};

Tabs.propTypes = {
  activeTabIndex: PropTypes.number.isRequired,
  onChange: PropTypes.func.isRequired,
  children: PropTypes.node.isRequired,
};

export default Tabs;

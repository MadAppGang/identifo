import PropTypes from 'prop-types';
import React, { useEffect, useRef, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import { useQuery } from '~/hooks/useQuery';
import Tab from './Tab';

const getElementWidth = el => el ? el.offsetWidth : 0;

const getUrl = (group, value) => {
  const searchParams = new URLSearchParams(window.location.search);
  searchParams.set(group, value);
  return `${window.location.pathname}?${searchParams.toString()}`;
};

const toSnakeCase = (value) => {
  return value.toLowerCase().replaceAll(' ', '_');
};

const getActiveTabIndex = (childrens, urlTab) => {
  const idx = childrens.findIndex(child => toSnakeCase(child.props.title) === urlTab);
  return idx === -1 ? 0 : idx;
};

const Tabs = ({ onChange, group, children }) => {
  const tabUrl = useQuery().get(group);
  const tabContent = [...children].pop();
  const history = useHistory();
  const tabs = children.slice(0, -1);
  if (tabContent.type === Tab) {
    throw Error('Last child should be tab content');
  }

  const tabRefsRef = useRef([]);
  const [tabWidths, setTabWidths] = useState([]);

  const onTabChange = () => {
    if (onChange) onChange();
  };

  useEffect(() => {
    if (!tabRefsRef.current.length) {
      return;
    }

    setTabWidths(tabRefsRef.current.map(ref => getElementWidth(ref.current)));
  }, [tabRefsRef.current]);

  useEffect(() => {
    if (!tabUrl) {
      history.replace(getUrl(group, toSnakeCase(children[0].props.title)));
    }
  }, [tabUrl]);
  return (
    <>
      <div className="iap-tabs-tablist">
        {React.Children.map(tabs, (child, index) => {
          const ref = React.createRef();

          tabRefsRef.current.push(ref);
          return (
            <Link to={getUrl(group, toSnakeCase(child.props.title))}>
              {' '}
              {React.cloneElement(child, {
                ref,
                isActive: getActiveTabIndex(tabs, tabUrl) === index,
                key: child.props.title,
                onClick: onTabChange,
              })}
            </Link>
          );
        })}
        <div
          className="tab-underline"
          style={{
            width: tabWidths[getActiveTabIndex(tabs, tabUrl)],
            left: tabWidths.slice(0, getActiveTabIndex(tabs, tabUrl)).reduce((a, b) => a + b, 0),
          }}
        />
      </div>
      {tabContent}
    </>
  );
};

Tabs.propTypes = {
  onChange: PropTypes.func,
  children: PropTypes.node.isRequired,
};

Tabs.defaultProps = {
  onChange: null,
};

export default Tabs;

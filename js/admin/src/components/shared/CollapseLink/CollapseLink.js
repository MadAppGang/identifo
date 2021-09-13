import { Collapse } from 'react-collapse';
import React, { useState } from 'react';
import DropdownIcon from '~/components/icons/DropdownIcon';
import classnames from 'classnames';
import './index.css';

export const CollapseItem = (props) => {
  const iconClass = classnames('iap-collapse-link--icon', {
    'iap-collapse-link--icon-opened': props.isOpen,
  });
  return (
    <div className="iap-collapse-link">
      <div className="iap-collapse-link--title" onClick={() => props.handleOpen(props.data)} role="presentation">
        <span>{props.title}</span>
        <DropdownIcon className={iconClass} />
      </div>
      <div className="iap-collapse-link--content">
        <Collapse isOpened={props.isOpen}>
          {props.children}
        </Collapse>
      </div>
    </div>
  );
};

export const CollapseLinks = ({ collapse, children }) => {
  const [opened, setOpened] = useState([]);
  const isOpen = idx => opened.includes(idx);

  const onItemClick = (idx) => {
    const openedIdx = opened.indexOf(idx);
    if (openedIdx >= 0) {
      setOpened(s => (collapse ? s.filter(item => item !== idx) : []));
    } else {
      setOpened(s => (collapse ? [...s, idx] : [idx]));
    }
  };
  return (
    <div>
      {React.Children.map(children, (child, idx) => (
        React.cloneElement(child, { handleOpen: onItemClick, isOpen: isOpen(idx), data: idx })
      ))}
    </div>
  );
};

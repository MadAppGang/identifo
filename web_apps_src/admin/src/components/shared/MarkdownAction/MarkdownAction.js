import React from 'react';
import './MarkdownAction.css';
import { Link } from 'react-router-dom';
import { MarkdownActionIcon } from '~/components/icons/MarkdownActionIcon/MarkdownActionIcon';

export const MarkdownAction = ({ title }) => {
  return (
    <div className="iap-markdown-action">
      <div className="iap-markdown-action__in">
        <div className="iap-markdown-action__icon">
          <MarkdownActionIcon />
        </div>
        <div className="iap-markdown-action__content">
          <span>{title}</span>
          <Link to="/faq" className="iap-markdown-action__btn">
          View more
          </Link>
        </div>
      </div>
    </div>
  );
};

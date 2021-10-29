import React from 'react';
import { MarkdownAction } from '~/components/shared/MarkdownAction/MarkdownAction';
import { useMarkdown } from '~/hooks/useMarkdown';
import './WithMarkdown.css';


export const WithMarkdown = ({ children }) => {
  const { hasMarkdown, parsedMarkdown } = useMarkdown();

  return (
    <div className="iap-with-markdown">
      {children}
      {hasMarkdown
        && (
          <div className="iap-with-markdown__action">
            <MarkdownAction title={parsedMarkdown.title} />
          </div>
        )}
    </div>
  );
};

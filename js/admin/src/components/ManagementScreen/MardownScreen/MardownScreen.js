import Markdown from 'markdown-to-jsx';
import React, { useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import { Arrow } from '~/components/icons/Arrow/Arrow';
import Header from '~/components/ManagementScreen/Header';
import Container from '~/components/shared/Container/index';
import { localStorageKeys } from '~/enums';
import './MarkdownScreen.css';


export const MarkodwnScreen = () => {
  const history = useHistory();

  if (!localStorage.getItem(localStorageKeys.markdown)) {
    history.goBack();
  }

  const markdown = JSON.parse(localStorage.getItem(localStorageKeys.markdown));

  useEffect(() => {
    return () => {
      localStorage.removeItem(localStorageKeys.markdown);
    };
  }, []);

  return (
    <div className="iap-markdown-screen">
      <Header />
      <Container>
        <div className="iap-markdown-screen__in">
          <div className="iap-markdown-screen__title">Instructions</div>
          <div className="iap-markdown-screen__back-btn" onClick={history.goBack} role="presentation">
            <Arrow />
            Go back
          </div>
          <div className="iap-markdown-screen__content">
            <Markdown>
              {markdown.body}
            </Markdown>
          </div>
        </div>
      </Container>
    </div>
  );
};

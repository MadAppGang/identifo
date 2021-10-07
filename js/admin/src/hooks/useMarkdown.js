import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { markdownMapper } from '~/utils/markdowns';
import { localStorageKeys } from '~/enums';

const markdownBasePath = 'src/markdowns/';


export const useMarkdown = () => {
  const [hasMarkdown, setHasMarkdown] = useState(false);
  const [parsedMarkdown, setParsedMarkdown] = useState({ title: '', body: '' });
  const location = useLocation();

  useEffect(() => {
    const pathToFile = `${location.pathname.replaceAll('/', '-')}.md`;
    fetch(`${markdownBasePath}${pathToFile}`)
      .then((res) => {
        if (res.ok) {
          setHasMarkdown(true);
          return res.text();
        }
        return Promise.reject();
      })
      .then((text) => {
        const { title, body } = markdownMapper(text);
        setParsedMarkdown({ title, body });
      })
      .catch(() => setHasMarkdown(false));
  }, [location]);

  useEffect(() => {
    if (hasMarkdown) {
      localStorage.setItem(localStorageKeys.markdown, JSON.stringify(parsedMarkdown));
    }
  }, [hasMarkdown, parsedMarkdown]);
  return { hasMarkdown, parsedMarkdown };
};

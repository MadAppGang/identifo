import { useEffect, useState } from 'react';
import { matchPath, useLocation, useRouteMatch } from 'react-router-dom';
import { localStorageKeys } from '~/enums';
import { paramRoutes } from '~/markdowns/paramRoutes/paramRoutes';
import { markdownMapper } from '~/utils/markdowns';

const markdownBasePath = '/src/markdowns';

const searchRegex = /[=,&]/gm;
const slashRegex = /(\/:*)/g;

const notFoundMarkdowns = new Map();

export const useMarkdown = () => {
  const [hasMarkdown, setHasMarkdown] = useState(false);
  const [parsedMarkdown, setParsedMarkdown] = useState({ title: '', body: '' });
  const location = useLocation();
  const params = useRouteMatch();

  // check is url match paramRoutes to exclude url params
  const match = matchPath(location.pathname, {
    path: Object.values(paramRoutes),
    exact: true,
  });

  useEffect(() => {
    const serializedSearchStr = location.search.replace(searchRegex, '_').slice(1);
    // serialize param from url for correct filename
    const serializedMatchPath = (match && match.path.replace(params.url, '').replace(slashRegex, '_')) || '';
    // create filename or use index name
    const filename = (`${serializedMatchPath}_${serializedSearchStr}`.replace(/(^_)|(_$)/g, '')) || 'index';
    const pathToFile = `${markdownBasePath}${params.url}/${filename}.md`;

    if (notFoundMarkdowns.has(pathToFile)) {
      setHasMarkdown(false);
      return;
    }

    fetch(pathToFile)
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
      .catch(() => {
        notFoundMarkdowns.set(pathToFile, false);
        setHasMarkdown(false);
      });
  }, [location]);

  useEffect(() => {
    if (hasMarkdown) {
      localStorage.setItem(localStorageKeys.markdown, JSON.stringify(parsedMarkdown));
    }
  }, [hasMarkdown, parsedMarkdown]);

  return { hasMarkdown, parsedMarkdown };
};

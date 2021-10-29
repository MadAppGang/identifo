const markdownTitleRegex = /^#\w.+/m;

export const markdownMapper = (markdownStr) => {
  const title = markdownStr.match(markdownTitleRegex)[0].slice(1);
  const body = markdownStr.replace(markdownTitleRegex, '');
  return { title, body };
};

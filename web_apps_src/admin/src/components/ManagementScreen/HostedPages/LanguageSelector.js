import React from 'react';
import classnames from 'classnames';

const BTN_HEIGHT = 45;

const LanguageSelector = (props) => {
  const { languages, selected, onChange } = props;
  const selectedLangIndex = languages.findIndex(item => item === selected);

  return (
    <div className="template-lang-selector">
      {
        languages.map((item) => {
          const className = classnames('template-lang-selector__btn', {
            'template-lang-selector__btn--active': item === selected,
          });

          return (
            <button
              key={item}
              type="button"
              className={className}
              onClick={() => onChange(item)}
            >
              {item}
            </button>
          );
        })
      }
      <div
        className="template-lang-selector__indicator"
        style={{ top: selectedLangIndex * BTN_HEIGHT }}
      />
    </div>
  );
};

export default LanguageSelector;

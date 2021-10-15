import React, { useRef, useState } from 'react';
import Input from '~/components/shared/Input';
import ValueTag from './ValueTag';
import AddIcon from '~/components/icons/AddIcon';

const MultipleInput = (props) => {
  const { values, onChange } = props;
  const [value, setValue] = useState('');
  const inputElRef = useRef(null);

  const deleteValue = (valueToDelete) => {
    onChange(values.filter(item => item !== valueToDelete));
  };

  const addValue = () => {
    if (!value.trim()) {
      return;
    }

    if (!values.includes(value.trim())) {
      onChange(values.concat(value));
    }

    setValue('');
  };

  const handleKeyPress = (event) => {
    if (event.key !== 'Enter') {
      return;
    }

    addValue(value);
    event.preventDefault();
  };

  const handleBlur = () => {
    if (value.trim()) {
      addValue(value);
    }
  };

  return (
    <div>
      <Input
        value={value}
        onChange={e => setValue(e.target.value)}
        onKeyPress={handleKeyPress}
        onBlur={handleBlur}
        placeholder={props.placeholder}
        autoComplete="off"
        ref={inputElRef}
        errorMessage={props.errorMessage}
        renderButton={() => (
          <button
            type="button"
            className="multiple-input__add-btn"
            onClick={() => inputElRef.current.focus()}
          >
            <AddIcon className="multiple-input__add-btn-icon" />
          </button>
        )}
        style={{
          paddingRight: '50px',
        }}
      />
      <ul className="multiple-input__value-tags">
        {values.map(item => (
          <ValueTag key={item} onDeleteClick={() => deleteValue(item)}>
            {item}
          </ValueTag>
        ))}
      </ul>
    </div>
  );
};

MultipleInput.defaultProps = {
  values: [],
  onChange: () => {},
};

export default MultipleInput;

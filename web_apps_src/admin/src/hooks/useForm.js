import { useState } from 'react';
import update from '@madappgang/update-by-path';

const useForm = (initialState, validate, submit) => {
  const [values, setValues] = useState(initialState || {});
  const [errors, setErrors] = useState({});

  const handleSubmit = (event) => {
    event.preventDefault();

    if (validate) {
      const validationResult = validate(values);

      if (Object.keys(validationResult).length > 0) {
        setErrors(validationResult);
        return;
      }
    }

    submit(values);

    if (Object.keys(errors).length > 0) {
      setErrors({});
    }
  };

  const setValue = (name, value) => {
    if (name in errors) {
      setErrors(update(errors, name, ''));
    }

    setValues(update(values, name, value));
  };

  const handleChange = ({ target }) => {
    setValue(target.name, target.value);
  };

  const handleBlur = () => {};

  return {
    values,
    errors,
    setErrors,
    setValue,
    setValues,
    handleSubmit,
    handleChange,
    handleBlur,
  };
};

export default useForm;

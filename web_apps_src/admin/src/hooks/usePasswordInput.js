import { useEffect, useState } from 'react';

export const usePasswordInput = (value, symbol = '*', maskLength) => {
  const [passwordShown, setPasswordShown] = useState(false);
  const [passwordMask, setPasswordMask] = useState('');

  useEffect(() => {
    if (value) {
      const maskedValue = value.replace(/./gm, symbol);
      if (passwordShown) {
        setPasswordMask(value);
      } else if (maskLength && value.length > maskLength) {
        setPasswordMask(maskedValue.slice(0, maskLength));
      } else {
        setPasswordMask(maskedValue);
      }
    }
  }, [value, passwordShown]);


  return [passwordMask, passwordShown, setPasswordShown];
};

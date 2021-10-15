import React from 'react';
import Input from '~/components/shared/Input';
import SearchIcon from '~/components/icons/SearchIcon.svg';
import { throttle } from '~/utils/fn';

import './SearchInput.css';

const SearchInput = ({ placeholder, timeout, onChange }) => {
  const [query, setQuery] = React.useState('');

  // for the memoized function to access
  const queryRef = React.useRef(query);

  const dispatchChange = React.useMemo(() => {
    return throttle(() => onChange(queryRef.current), timeout);
  }, [timeout, onChange]);

  const handleSearchChange = (value) => {
    setQuery(value);
    queryRef.current = value;
    dispatchChange(value);
  };

  return (
    <div className="iap-search-input">
      <Input
        Icon={SearchIcon}
        value={query}
        onValue={handleSearchChange}
        placeholder={placeholder}
      />
    </div>
  );
};

export default SearchInput;

import { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { fetchServerSetings } from '~/modules/settings/actions';

export const useSettings = () => {
  const dispatch = useDispatch();

  useEffect(() => {
    const fetch = async () => {
      await dispatch(fetchServerSetings());
    };
    fetch();
  }, []);
};
